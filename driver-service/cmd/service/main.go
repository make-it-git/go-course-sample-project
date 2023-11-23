package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"driver-service/internal/config"
	"driver-service/internal/db/repository"
	driver_order "driver-service/internal/generated/proto/driver.order"
	"driver-service/internal/handlers"
	"driver-service/internal/logger"
	"driver-service/internal/otel"
	"driver-service/internal/services/driver_search"
	"driver-service/internal/services/location_updater"
	"driver-service/internal/services/order"
)

func main() {
	log := logger.New()
	cfg, err := config.FromEnv()
	if err != nil {
		log.WithError(err, "get cfg")
		os.Exit(1)
	}

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, cfg.DatabaseUrl)
	if err != nil {
		log.WithError(err, "database connect")
		os.Exit(1)
	}
	defer conn.Close()

	rdbBroker := redis.NewClient(&redis.Options{
		Addr:     cfg.BrokerURL,
		Password: "",
		Protocol: 3,
	})
	if err := rdbBroker.Ping(ctx).Err(); err != nil {
		log.WithError(err, "redis broker ping")
		os.Exit(1)
	}
	defer rdbBroker.Close()

	driverSearch := driver_search.NewDriverSearchService()
	orderRepository := repository.NewOrderRepository(conn)
	orderService := order.NewOrderService(orderRepository, driverSearch, log)

	ctxWithCancel, cancelLocationUpdater := context.WithCancel(ctx)
	locationUpdater := location_updater.NewLocationUpdater(rdbBroker, orderRepository, log)
	go locationUpdater.Run(ctxWithCancel)

	handler := handlers.NewHandler(orderService)

	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.05, 0.1, 0.5, 1, 2, 3, 5}),
		),
	)
	prometheus.MustRegister(srvMetrics)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingInterceptor(log, cfg),
			srvMetrics.UnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			srvMetrics.StreamServerInterceptor(),
			otelgrpc.StreamServerInterceptor(),
		),
	)
	reflection.Register(s)
	driver_order.RegisterOrderServer(s, handler)

	serviceName := "driver-service"
	serviceVersion := os.Getenv("SERVICE_VERSION")
	otelShutdown, err := otel.SetupOTelSDK(ctx, serviceName, serviceVersion, cfg.Env == config.ProdEnv)
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Info("Listen", "addr", cfg.ListenAddrAndPort())
		listener, err := net.Listen("tcp", cfg.ListenAddrAndPort())
		if err != nil {
			log.WithError(err, "listen")
			close(done)
			return
		}
		if err := s.Serve(listener); err != nil {
			log.WithError(err, "listen")
			close(done)
		}
	}()

	go func() {
		err := http.ListenAndServe(cfg.PromListenAddrAndPort(), nil)
		if err != nil {
			log.WithError(err, "prom listen")
		}
	}()

	<-done
	cancelLocationUpdater()
	log.Info("Listen stopped")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer func() {
		cancel()
	}()

	s.GracefulStop()
	log.Info("Shutdown completed")
}

func loggingInterceptor(log logger.Log, cfg *config.Config) grpc.UnaryServerInterceptor {
	f := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if cfg.Env == config.LocalEnv {
			log.Info("Handle", "method", info.FullMethod, "req", req)
		}
		h, err := handler(ctx, req)
		return h, err
	}
	return f
}

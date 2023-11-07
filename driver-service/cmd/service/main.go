package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"driver-service/internal/config"
	"driver-service/internal/db/repository"
	driver_order "driver-service/internal/generated/proto/driver.order"
	"driver-service/internal/handlers"
	"driver-service/internal/logger"
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

	s := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor(log, cfg)))
	reflection.Register(s)
	driver_order.RegisterOrderServer(s, handler)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

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

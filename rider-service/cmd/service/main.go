package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yarlson/chiprom"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"rider-service/internal/config"
	"rider-service/internal/db/repository"
	driver_order "rider-service/internal/generated/proto/driver.order"
	rider "rider-service/internal/generated/schema"
	"rider-service/internal/handlers"
	"rider-service/internal/logger"
	"rider-service/internal/now_time"
	"rider-service/internal/services/driver_sender"
	"rider-service/internal/services/order"
	"rider-service/internal/services/price_estimator"
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

	grpcConn, err := grpc.DialContext(
		context.Background(),
		cfg.DriverServiceLocation,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.WithError(err, "grpc connect")
		os.Exit(1)
	}
	grcpClient := driver_order.NewOrderClient(grpcConn)
	driverSenderService := driver_sender.NewDriverSenderService(grcpClient)

	priceEstimator := price_estimator.NewPriceEstimatorService()
	orderRepository := repository.NewOrderRepository(conn)
	orderService := order.NewOrderService(orderRepository, priceEstimator, now_time.Get, driverSenderService)

	handle := handlers.New(log, now_time.Get, orderService)

	r := chi.NewRouter()
	swagger, err := rider.GetSwagger()
	if err != nil {
		log.WithError(err, "get swagger")
		os.Exit(1)
	}
	r.Use(middleware.OapiRequestValidator(swagger))
	r.Use(chimiddleware.Recoverer)
	r.Use(chiprom.NewMiddleware("rider-service"))
	if cfg.Env == config.LocalEnv {
		r.Use(chimiddleware.Logger)
	}

	baseRouter := chi.NewRouter()
	baseRouter.Handle("/metrics", promhttp.Handler())

	rider.HandlerFromMux(handle, r)
	baseRouter.Mount("/", r)

	s := &http.Server{
		Handler: baseRouter,
		Addr:    cfg.ListenAddrAndPort(),
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("Listen", "addr", cfg.ListenAddrAndPort())
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err, "listen")
			close(done)
		}
	}()

	<-done
	log.Info("Listen stopped")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := s.Shutdown(ctx); err != nil {
		log.Error("Shutdown error", "error", err.Error())
		os.Exit(1)
	}
	log.Info("Shutdown completed")
}

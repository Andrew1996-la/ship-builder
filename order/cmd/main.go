package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderapi "github.com/Andrew1996-la/ship-builder/order/internal/api/order/v1"
	grpcclientInventory "github.com/Andrew1996-la/ship-builder/order/internal/client/grpc/inventory/v1"
	grpcclientPayment "github.com/Andrew1996-la/ship-builder/order/internal/client/grpc/payment/v1"
	orderrepo "github.com/Andrew1996-la/ship-builder/order/internal/repository/order"
	orderservice "github.com/Andrew1996-la/ship-builder/order/internal/service/order"
	orderv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/payment/v1"
)

const (
	httpPort = "localhost:8080"

	inventoryServiceAddress = "127.0.0.1:50051"
	paymentServiceAddress   = "127.0.0.1:50052"

	shutdownTimeout = 10 * time.Second
)

func main() {
	ctx := context.Background()
	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		slog.Error("dbURI не задан")
		os.Exit(1)
	}

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		slog.Error("создание пула соединений", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		slog.Error("PostgreSQL не отвечает", "error", err)
		os.Exit(1)
	}

	slog.Info("успешно подключились к PostgreSQL")

	// Создать gRPC соединение с InventoryService
	inventoryConn, err := grpc.NewClient(
		inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("не удалось подключиться к InventoryService", "error", err)
		os.Exit(1)
	}

	defer inventoryConn.Close()

	paymentConn, err := grpc.NewClient(
		paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("не удалось подключиться к PaymentService", "error", err)
		os.Exit(1)
	}

	defer paymentConn.Close()

	repository := orderrepo.New()

	inventoryClient := grpcclientInventory.New(
		inventoryv1.NewInventoryServiceClient(inventoryConn),
	)

	paymentClient := grpcclientPayment.New(
		paymentv1.NewPaymentServiceClient(paymentConn),
	)

	service := orderservice.New(
		repository,
		inventoryClient,
		paymentClient,
	)

	api := orderapi.New(service)

	// Создать OpenAPI сервер
	orderServer, err := orderv1.NewServer(api)
	if err != nil {
		slog.Error("ошибка создания сервера OpenAPI", "error", err)
		os.Exit(1)
	}

	// Настройка http server
	server := &http.Server{
		Addr:              httpPort,
		Handler:           orderServer,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Запуск сервера
	go func() {
		slog.Info("http server запущен на порту", "port", httpPort)
		listenServerErr := server.ListenAndServe()

		if listenServerErr != nil && !errors.Is(listenServerErr, http.ErrServerClosed) {
			slog.Error("ошибка запуска сервера", "error", listenServerErr)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("завершение работы сервера")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("произошла ошибка при остановке сервера", "error", err)
	}

	slog.Info("сервер успешно остановлен")
}

package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	inventoryapi "github.com/Andrew1996-la/ship-builder/inventory/internal/api/inventory/v1"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/interceptor"
	partrepo "github.com/Andrew1996-la/ship-builder/inventory/internal/repository/part"
	partservice "github.com/Andrew1996-la/ship-builder/inventory/internal/service/part"
	inventoryv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/inventory/v1"
)

const (
	grpcAddress = "localhost:50051"

	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
)

func main() {
	ctx := context.Background()

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		slog.Error("DB_URI не задан")
		os.Exit(1)
	}

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		slog.Error("не удалось подключиться к PostgreSQL", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		slog.Error("PostgreSQL не отвечает", "error", err)
		os.Exit(1)
	}

	slog.Info("успешно подключились к PostgreSQL")

	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: true,
		}),
	)
	repository := partrepo.New(pool)
	service := partservice.New(repository)
	api := inventoryapi.New(service)

	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)

	// Включаем reflection для postman/grpcurl
	reflection.Register(grpcServer)

	slog.Info("запуск InventoryService", "адрес", grpcAddress)

	// graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(quit)

		<-quit
		slog.Info("остановка InventoryService")

		grpcServer.GracefulStop()
	}()

	err = grpcServer.Serve(lis)
	if err != nil {
		slog.Error("ошибка запуска сервера", "error", err)
		os.Exit(1)
	}
}

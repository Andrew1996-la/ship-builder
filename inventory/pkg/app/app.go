package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	inventoryapi "github.com/Andrew1996-la/ship-builder/inventory/internal/api/inventory/v1"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/interceptor"
	partrepo "github.com/Andrew1996-la/ship-builder/inventory/internal/repository/part"
	partservice "github.com/Andrew1996-la/ship-builder/inventory/internal/service/part"
	inventoryv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/inventory/v1"
)

func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	}
}

func RegisterServices(grpcServer *grpc.Server, pool *pgxpool.Pool) {
	repository := partrepo.New(pool)
	service := partservice.New(repository)
	api := inventoryapi.New(service)

	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)
}

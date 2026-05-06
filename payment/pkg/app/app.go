package app

import (
	"google.golang.org/grpc"

	paymentapi "github.com/Andrew1996-la/ship-builder/payment/internal/api/payment/v1"
	"github.com/Andrew1996-la/ship-builder/payment/internal/interceptor"
	paymentservice "github.com/Andrew1996-la/ship-builder/payment/internal/service/payment"
	paymentv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/payment/v1"
)

func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	}
}

func RegisterServices(grpcServer *grpc.Server) {
	service := paymentservice.New()
	api := paymentapi.New(service)

	paymentv1.RegisterPaymentServiceServer(grpcServer, api)
}

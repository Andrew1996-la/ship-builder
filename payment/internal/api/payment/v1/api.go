package v1

import paymentv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/payment/v1"

type api struct {
	paymentv1.UnimplementedPaymentServiceServer
	service PaymentService
}

func New(service PaymentService) *api {
	return &api{
		service: service,
	}
}

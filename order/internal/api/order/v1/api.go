package v1

import orderv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/openapi/order/v1"

type api struct {
	orderv1.UnimplementedHandler
	service OrderService
}

func New(service OrderService) *api {
	return &api{
		service: service,
	}
}

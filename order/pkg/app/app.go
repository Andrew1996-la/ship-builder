package app

import (
	"net/http"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	orderapi "github.com/Andrew1996-la/ship-builder/order/internal/api/order/v1"
	grpcclientInventory "github.com/Andrew1996-la/ship-builder/order/internal/client/grpc/inventory/v1"
	grpcclientPayment "github.com/Andrew1996-la/ship-builder/order/internal/client/grpc/payment/v1"
	orderrepo "github.com/Andrew1996-la/ship-builder/order/internal/repository/order"
	orderservice "github.com/Andrew1996-la/ship-builder/order/internal/service/order"
	orderv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/payment/v1"
)

func NewHTTPHandler(
	pool *pgxpool.Pool,
	txManager trm.Manager,
	inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
) (http.Handler, error) {
	repository := orderrepo.New(pool, txManager)
	service := orderservice.New(
		repository,
		grpcclientInventory.New(inventoryClient),
		grpcclientPayment.New(paymentClient),
	)
	api := orderapi.New(service)

	return orderv1.NewServer(api)
}

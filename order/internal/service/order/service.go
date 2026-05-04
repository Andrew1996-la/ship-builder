package order

import "github.com/avito-tech/go-transaction-manager/trm/v2"

type service struct {
	repository      Repository
	inventoryClient InventoryClient
	paymentClient   PaymentClient
	txManager       trm.Manager
}

func New(
	repository Repository,
	inventoryClient InventoryClient,
	paymentClient PaymentClient,
) *service {
	return &service{
		repository:      repository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

func NewWithTx(
	repository Repository,
	inventoryClient InventoryClient,
	paymentClient PaymentClient,
	txManager trm.Manager,
) *service {
	return &service{
		repository:      repository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		txManager:       txManager,
	}
}

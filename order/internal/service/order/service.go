package order

type service struct {
	repository      Repository
	inventoryClient InventoryClient
	paymentClient   PaymentClient
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

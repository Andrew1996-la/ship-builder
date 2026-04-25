package v1

import inventoryv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/inventory/v1"

type api struct {
	inventoryv1.UnimplementedInventoryServiceServer
	partService PartService
}

func New(partService PartService) *api {
	return &api{
		partService: partService,
	}
}

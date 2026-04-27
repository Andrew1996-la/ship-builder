package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/Andrew1996-la/ship-builder/order/internal/model"
	inventoryv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/inventory/v1"
)

type Client struct {
	client inventoryv1.InventoryServiceClient
}

func New(client inventoryv1.InventoryServiceClient) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) ListParts(ctx context.Context, uuids []uuid.UUID) ([]model.Part, error) {
	rawUUIDs := make([]string, 0, len(uuids))

	for _, id := range uuids {
		rawUUIDs = append(rawUUIDs, id.String())
	}

	resp, err := c.client.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: rawUUIDs,
	})
	if err != nil {
		return nil, err
	}

	parts := make([]model.Part, 0, len(resp.GetParts()))

	for _, part := range resp.GetParts() {
		id, err := uuid.Parse(part.GetUuid())
		if err != nil {
			return nil, err
		}

		parts = append(parts, model.Part{
			UUID:          id,
			Name:          part.GetName(),
			Price:         part.GetPrice(),
			StockQuantity: part.GetStockQuantity(),
		})
	}

	return parts, nil
}

package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Andrew1996-la/ship-builder/order/internal/client/grpc/inventory/v1/converter"
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
	resp, err := c.client.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: converter.UUIDsToRaw(uuids),
	})
	if err != nil {
		return nil, fmt.Errorf("получить детали из сервиса склада: %w", err)
	}

	parts, err := converter.ProtoPartsToModel(resp.GetParts())
	if err != nil {
		return nil, fmt.Errorf("преобразовать ответ сервиса склада: %w", err)
	}

	return parts, nil
}

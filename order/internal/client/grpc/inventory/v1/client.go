package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Andrew1996-la/ship-builder/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/Andrew1996-la/ship-builder/order/internal/errors"
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
		Uuids: converter.ToRawUUIDs(uuids),
	})
	if err != nil {
		return nil, mapListPartsError(err)
	}

	parts, err := converter.ToModelParts(resp.GetParts())
	if err != nil {
		return nil, fmt.Errorf("преобразовать ответ сервиса склада: %w", err)
	}

	return parts, nil
}

func mapListPartsError(err error) error {
	switch status.Code(err) {
	case codes.NotFound:
		return fmt.Errorf("получить детали из сервиса склада: %w", errors.Join(errs.ErrPartNotFound, err))
	default:
		return fmt.Errorf("получить детали из сервиса склада: %w", err)
	}
}

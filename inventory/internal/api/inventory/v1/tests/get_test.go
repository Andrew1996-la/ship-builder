package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	inventoryapi "github.com/Andrew1996-la/ship-builder/inventory/internal/api/inventory/v1"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/api/inventory/v1/mocks"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/converter"
	errs "github.com/Andrew1996-la/ship-builder/inventory/internal/errors"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
	inventoryv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/inventory/v1"
)

func TestGet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	partUUID := uuid.New().String()

	expectedPart := model.Part{
		UUID:          uuid.MustParse(partUUID),
		Name:          "Ионный двигатель",
		Description:   "Базовый двигатель",
		Price:         300000,
		PartType:      model.PartTypeEngine,
		StockQuantity: 10,
		CreatedAt:     time.Now(),
	}

	tests := []struct {
		name        string
		req         *inventoryv1.GetPartRequest
		setupMock   func(service *mocks.PartService)
		expected    *inventoryv1.Part
		expectedErr error
	}{
		{
			name: "success",
			req: &inventoryv1.GetPartRequest{
				Uuid: partUUID,
			},
			setupMock: func(service *mocks.PartService) {
				service.EXPECT().
					Get(ctx, partUUID).
					Return(expectedPart, nil)
			},
			expected: converter.ModelToProtoPart(expectedPart),
		},
		{
			name: "service error",
			req: &inventoryv1.GetPartRequest{
				Uuid: partUUID,
			},
			setupMock: func(service *mocks.PartService) {
				service.EXPECT().
					Get(ctx, partUUID).
					Return(model.Part{}, errs.ErrPartNotFound)
			},
			expectedErr: errs.ErrPartNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := mocks.NewPartService(t)
			tt.setupMock(service)

			api := inventoryapi.New(service)

			resp, err := api.Get(ctx, tt.req)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Part)

			assert.Equal(t, tt.expected.Uuid, resp.Part.Uuid)
			assert.Equal(t, tt.expected.Name, resp.Part.Name)
			assert.Equal(t, tt.expected.Description, resp.Part.Description)
			assert.Equal(t, tt.expected.Price, resp.Part.Price)
			assert.Equal(t, tt.expected.PartType, resp.Part.PartType)
			assert.Equal(t, tt.expected.StockQuantity, resp.Part.StockQuantity)
		})
	}
}

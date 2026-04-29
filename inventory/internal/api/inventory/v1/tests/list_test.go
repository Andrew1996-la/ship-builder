package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inventoryapi "github.com/Andrew1996-la/ship-builder/inventory/internal/api/inventory/v1"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/api/inventory/v1/mocks"
	errs "github.com/Andrew1996-la/ship-builder/inventory/internal/errors"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
	inventoryv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/inventory/v1"
)

func TestList(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	partUUID := uuid.New().String()

	req := &inventoryv1.ListPartsRequest{
		Uuids: []string{partUUID},
	}

	expectedFilter := model.PartFilter{
		UUIDs:    []string{partUUID},
		PartType: model.PartTypeUnspecified,
	}

	expectedParts := []model.Part{
		{
			UUID:          uuid.MustParse(partUUID),
			Name:          "Ионный двигатель",
			Description:   "Базовый двигатель",
			Price:         300000,
			PartType:      model.PartTypeEngine,
			StockQuantity: 10,
			CreatedAt:     time.Now(),
		},
	}

	tests := []struct {
		name         string
		req          *inventoryv1.ListPartsRequest
		setupMock    func(service *mocks.PartService)
		expectedLen  int
		expectedErr  error
		expectedCode codes.Code
		expectedMsg  string
	}{
		{
			name: "успешный сценарий",
			req:  req,
			setupMock: func(service *mocks.PartService) {
				service.EXPECT().
					List(ctx, expectedFilter).
					Return(expectedParts, nil)
			},
			expectedLen: 1,
		},
		{
			name: "ошибка сервиса",
			req:  req,
			setupMock: func(service *mocks.PartService) {
				service.EXPECT().
					List(ctx, expectedFilter).
					Return(nil, errs.ErrPartNotFound)
			},
			expectedErr:  errs.ErrPartNotFound,
			expectedCode: codes.NotFound,
			expectedMsg:  "детали не найдены",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := mocks.NewPartService(t)
			tt.setupMock(service)

			api := inventoryapi.New(service)

			resp, err := api.ListParts(ctx, tt.req)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedCode, status.Code(err))
				assert.Equal(t, tt.expectedMsg, status.Convert(err).Message())
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Len(t, resp.Parts, tt.expectedLen)

			assert.Equal(t, expectedParts[0].UUID.String(), resp.Parts[0].Uuid)
			assert.Equal(t, expectedParts[0].Name, resp.Parts[0].Name)
			assert.Equal(t, expectedParts[0].Price, resp.Parts[0].Price)
			assert.Equal(t, expectedParts[0].StockQuantity, resp.Parts[0].StockQuantity)
		})
	}
}

package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/Andrew1996-la/ship-builder/inventory/internal/errors"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
	partservice "github.com/Andrew1996-la/ship-builder/inventory/internal/service/part"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/service/part/mocks"
)

func TestList(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	filter := model.PartFilter{
		PartType: model.PartTypeEngine,
	}

	expectedParts := []model.Part{
		{
			UUID:          uuid.New(),
			Name:          "Ионный двигатель",
			Description:   "Базовый двигатель",
			Price:         300000,
			PartType:      model.PartTypeEngine,
			StockQuantity: 10,
			CreatedAt:     time.Now(),
		},
	}

	tests := []struct {
		name        string
		filter      model.PartFilter
		setupMock   func(repository *mocks.PartRepository)
		expected    []model.Part
		expectedErr error
	}{
		{
			name:   "успешный сценарий",
			filter: filter,
			setupMock: func(repository *mocks.PartRepository) {
				repository.EXPECT().
					List(ctx, filter).
					Return(expectedParts, nil)
			},
			expected: expectedParts,
		},
		{
			name:   "ошибка репозитория",
			filter: filter,
			setupMock: func(repository *mocks.PartRepository) {
				repository.EXPECT().
					List(ctx, filter).
					Return(nil, errs.ErrPartNotFound)
			},
			expectedErr: errs.ErrPartNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := mocks.NewPartRepository(t)
			tt.setupMock(repository)

			service := partservice.New(repository)

			actual, err := service.List(ctx, tt.filter)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
				assert.Nil(t, actual)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

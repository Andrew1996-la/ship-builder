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
		setupMock   func(repository *mocks.PartRepository)
		expected    model.Part
		expectedErr error
	}{
		{
			name: "успешный сценарий",
			setupMock: func(repository *mocks.PartRepository) {
				repository.EXPECT().
					Get(ctx, partUUID).
					Return(expectedPart, nil)
			},
			expected: expectedPart,
		},
		{
			name: "ошибка репозитория",
			setupMock: func(repository *mocks.PartRepository) {
				repository.EXPECT().
					Get(ctx, partUUID).
					Return(model.Part{}, errs.ErrPartNotFound)
			},
			expected:    model.Part{},
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

			actual, err := service.Get(ctx, partUUID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
				assert.Equal(t, tt.expected, actual)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

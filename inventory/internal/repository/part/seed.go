package part

import (
	"time"

	"github.com/google/uuid"

	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/repository/record"
)

func NewWithSeed() *repository {
	now := time.Now()

	return New(map[uuid.UUID]record.Part{
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"): {
			UUID:          "550e8400-e29b-41d4-a716-446655440001",
			Name:          "Алюминиевый корпус",
			Description:   "Лёгкий корпус для небольших кораблей",
			Price:         500000,
			PartType:      string(model.PartTypeHull),
			StockQuantity: 10,
			CreatedAt:     now,
		},
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"): {
			UUID:          "550e8400-e29b-41d4-a716-446655440002",
			Name:          "Титановый корпус",
			Description:   "Прочный корпус для средних кораблей",
			Price:         1500000,
			PartType:      string(model.PartTypeHull),
			StockQuantity: 5,
			CreatedAt:     now,
		},
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"): {
			UUID:          "550e8400-e29b-41d4-a716-446655440003",
			Name:          "Ионный двигатель C",
			Description:   "Базовый ионный двигатель класса C",
			Price:         300000,
			PartType:      string(model.PartTypeEngine),
			StockQuantity: 8,
			CreatedAt:     now,
		},
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"): {
			UUID:          "550e8400-e29b-41d4-a716-446655440004",
			Name:          "Ионный двигатель B",
			Description:   "Улучшенный ионный двигатель класса B",
			Price:         800000,
			PartType:      string(model.PartTypeEngine),
			StockQuantity: 3,
			CreatedAt:     now,
		},
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"): {
			UUID:          "550e8400-e29b-41d4-a716-446655440005",
			Name:          "Энергетический щит",
			Description:   "Стандартный энергетический щит",
			Price:         400000,
			PartType:      string(model.PartTypeShield),
			StockQuantity: 6,
			CreatedAt:     now,
		},
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"): {
			UUID:          "550e8400-e29b-41d4-a716-446655440006",
			Name:          "Лазерная пушка",
			Description:   "Точная лазерная пушка",
			Price:         250000,
			PartType:      string(model.PartTypeWeapon),
			StockQuantity: 7,
			CreatedAt:     now,
		},
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"): {
			UUID:          "550e8400-e29b-41d4-a716-446655440007",
			Name:          "Плазменный корпус",
			Description:   "Прочный корпус для больших кораблей",
			Price:         2000000,
			PartType:      string(model.PartTypeHull),
			StockQuantity: 0,
			CreatedAt:     now,
		},
	})
}

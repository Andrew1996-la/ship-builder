package service

import (
	"context"
	"log/slog"
	"sort"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryv1 "github.com/Andrew1996-la/ship-builder/shared/pkg/proto/inventory/v1"
)

// Part представляет деталь космического корабля.
type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64 // в копейках
	PartType      inventoryv1.PartType
	StockQuantity int64
	CreatedAt     *timestamppb.Timestamp
}

// InventoryServer реализует gRPC сервис.
type InventoryServer struct {
	inventoryv1.UnimplementedInventoryServiceServer
	parts map[uuid.UUID]Part
	mu    sync.RWMutex
}

// NewInventoryServer создаёт сервер с предзагруженными seed-данными.
func NewInventoryServer() *InventoryServer {
	now := timestamppb.Now()

	return &InventoryServer{
		parts: map[uuid.UUID]Part{
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440001",
				Name:          "Алюминиевый корпус",
				Description:   "Лёгкий корпус для небольших кораблей",
				Price:         500000, // 5000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 10,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440002",
				Name:          "Титановый корпус",
				Description:   "Прочный корпус для средних кораблей",
				Price:         1500000, // 15000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 5,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440003",
				Name:          "Ионный двигатель C",
				Description:   "Базовый ионный двигатель класса C",
				Price:         300000, // 3000₽
				PartType:      inventoryv1.PartType_PART_TYPE_ENGINE,
				StockQuantity: 8,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440004",
				Name:          "Ионный двигатель B",
				Description:   "Улучшенный ионный двигатель класса B",
				Price:         800000, // 8000₽
				PartType:      inventoryv1.PartType_PART_TYPE_ENGINE,
				StockQuantity: 3,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440005",
				Name:          "Энергетический щит",
				Description:   "Стандартный энергетический щит",
				Price:         400000, // 4000₽
				PartType:      inventoryv1.PartType_PART_TYPE_SHIELD,
				StockQuantity: 6,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440006",
				Name:          "Лазерная пушка",
				Description:   "Точная лазерная пушка",
				Price:         250000, // 2500₽
				PartType:      inventoryv1.PartType_PART_TYPE_WEAPON,
				StockQuantity: 7,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440007",
				Name:          "Плазменный корпус",
				Description:   "Прочный корпус для больших кораблей",
				Price:         2000000, // 20000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 0,
				CreatedAt:     now,
			},
		},
	}
}

// GetPart возвращает деталь по UUID.
func (s *InventoryServer) GetPart(
	ctx context.Context,
	req *inventoryv1.GetPartRequest,
) (*inventoryv1.GetPartResponse, error) {
	// TODO: Реализовать метод
	// 1. Проверить, что uuid не пустой → INVALID_ARGUMENT
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid обязателен")
	}

	// 2. Валидировать формат UUID → INVALID_ARGUMENT
	parsedUUID, err := uuid.Parse(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "неверный формат uuid: %s", req.GetUuid())
	}

	// 3. Найти деталь в map
	s.mu.RLock()
	part, ok := s.parts[parsedUUID]
	s.mu.RUnlock()

	// 4. Если не найдена → NOT_FOUND
	if !ok {
		return nil, status.Error(codes.NotFound, "деталь не наидена")
	}

	// 5. Преобразовать в inventoryv1.Part
	p := &inventoryv1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      part.PartType,
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}

	// 6. Вернуть деталь
	return &inventoryv1.GetPartResponse{
		Part: p,
	}, nil
}

// ListParts возвращает список деталей с опциональной фильтрацией по типу.
func (s *InventoryServer) ListParts(
	ctx context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {
	slog.Info("ListParts вызван", "uuids", req.GetUuids(), "part_type", req.GetPartType())
	if len(req.GetUuids()) > 0 {
		parts := make([]*inventoryv1.Part, 0, len(req.GetUuids()))

		for _, u := range req.Uuids {
			parsedUUID, err := uuid.Parse(u)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "неверный формат uuid: %s", u)
			}

			s.mu.RLock()
			part, ok := s.parts[parsedUUID]
			s.mu.RUnlock()

			if !ok {
				return nil, status.Error(codes.NotFound, "деталь не наидена")
			}

			parts = append(parts, &inventoryv1.Part{
				Uuid:          part.UUID,
				Name:          part.Name,
				Description:   part.Description,
				Price:         part.Price,
				PartType:      part.PartType,
				StockQuantity: part.StockQuantity,
				CreatedAt:     part.CreatedAt,
			})
		}

		return &inventoryv1.ListPartsResponse{
			Parts: parts,
		}, nil
	}
	// 2. Иначе если part_type == UNSPECIFIED → вернуть все детали
	s.mu.RLock()
	parts := make([]*inventoryv1.Part, 0, len(s.parts))

	for _, part := range s.parts {
		if req.GetPartType() != inventoryv1.PartType_PART_TYPE_UNSPECIFIED &&
			part.PartType != req.GetPartType() {
			continue
		}

		parts = append(parts, &inventoryv1.Part{
			Uuid:          part.UUID,
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			PartType:      part.PartType,
			StockQuantity: part.StockQuantity,
			CreatedAt:     part.CreatedAt,
		})
	}
	s.mu.RUnlock()

	// 3. Сортировка по имени
	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Name < parts[j].Name
	})

	return &inventoryv1.ListPartsResponse{
		Parts: parts,
	}, nil
}

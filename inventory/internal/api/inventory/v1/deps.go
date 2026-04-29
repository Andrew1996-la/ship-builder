package v1

import (
	"context"

	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
)

type PartService interface {
	// Get получить деталь по UUID.
	//
	// Возвращаемые ошибки:
	// - errs.ErrInvalidUUID — если UUID некорректный
	// - errs.ErrPartNotFound — если деталь не найдена
	Get(ctx context.Context, uuid string) (model.Part, error)

	// List получить список деталей по фильтру.
	//
	// Если filter.UUIDs переданы — возвращаются только указанные детали.
	// Если filter.PartType задан — применяется фильтрация по типу.
	//
	// Возвращаемые ошибки:
	// - errs.ErrInvalidUUID — если один из UUID некорректный
	// - errs.ErrPartNotFound — если хотя бы одна деталь не найдена
	List(ctx context.Context, filter model.PartFilter) ([]model.Part, error)
}

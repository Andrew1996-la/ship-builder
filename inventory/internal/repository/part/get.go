package part

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/Andrew1996-la/ship-builder/inventory/internal/errors"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/repository/converter"
)

// GetPart возвращает деталь по UUID.
func (r *repository) Get(
	ctx context.Context,
	id string,
) (model.Part, error) {
	parsedUuid, err := uuid.Parse(id)
	if err != nil {
		return model.Part{}, errs.ErrInvalidUUID
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.parts[parsedUuid]
	if !ok {
		return model.Part{}, errs.ErrPartNotFound
	}

	return converter.ToModelPart(part)
}

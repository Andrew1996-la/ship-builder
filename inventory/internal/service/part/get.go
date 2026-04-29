package part

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/Andrew1996-la/ship-builder/inventory/internal/errors"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
)

func (s service) Get(ctx context.Context, uuidRaw string) (model.Part, error) {
	id, err := uuid.Parse(uuidRaw)
	if err != nil {
		return model.Part{}, errs.ErrInvalidUUID
	}

	if id == uuid.Nil {
		return model.Part{}, errs.ErrInvalidUUID
	}

	return s.repository.Get(ctx, uuidRaw)
}

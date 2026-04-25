package part

import (
	"context"
	"sort"

	"github.com/google/uuid"

	errs "github.com/Andrew1996-la/ship-builder/inventory/internal/errors"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/model"
	"github.com/Andrew1996-la/ship-builder/inventory/internal/repository/converter"
)

func (r *repository) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	if len(filter.UUIDs) > 0 {
		return r.listByUUIDs(filter.UUIDs)
	}

	return r.listByType(filter.PartType)
}

func (r *repository) listByUUIDs(uuids []string) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := make([]model.Part, 0, len(uuids))

	for _, strUuid := range uuids {
		parsedUuid, err := uuid.Parse(strUuid)
		if err != nil {
			return nil, errs.ErrInvalidUUID
		}

		part, ok := r.parts[parsedUuid]
		if !ok {
			return nil, errs.ErrPartNotFound
		}

		modelPart, err := converter.RepositoryToModel(part)
		if err != nil {
			return nil, err
		}

		parts = append(parts, modelPart)
	}

	return parts, nil
}

func (r *repository) listByType(partType model.PartType) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := make([]model.Part, 0, len(r.parts))

	for _, part := range r.parts {
		modelPart, err := converter.RepositoryToModel(part)
		if err != nil {
			return nil, err
		}

		if modelPart.PartType == partType || partType == model.PartTypeUnspecified {
			parts = append(parts, modelPart)
		}
	}

	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Name < parts[j].Name
	})

	return parts, nil
}

package part

import (
	"sync"

	"github.com/google/uuid"

	"github.com/Andrew1996-la/ship-builder/inventory/internal/repository/record"
)

type repository struct {
	mu    sync.RWMutex
	parts map[uuid.UUID]record.Part
}

func New(parts map[uuid.UUID]record.Part) *repository {
	return &repository{
		parts: parts,
	}
}

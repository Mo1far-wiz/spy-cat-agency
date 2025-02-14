package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrorNotFound = errors.New("store: resource not found")
var QueryTimeoutDuration = 5 * time.Second
var ErrConflict = errors.New("store: resource already exists")

type CRUD[T any] interface {
	Create(context.Context, *T) error
	Update(context.Context, *T) error
	Delete(context.Context, int64) error
	GetByID(context.Context, int64) (*T, error)
	GetAll(context.Context) ([]T, error)
}

type Storage struct {
	Cat interface {
		CRUD[Cat]
		HasIncompleteMission(context.Context, int64) (bool, error)
	}
	Mission interface {
		CRUD[Mission]
		AssignCat(context.Context, int64, int64) error
		AddTarget(context.Context, int64, *Target) error
		RemoveTarget(context.Context, int64) error
		AddNote(context.Context, *Note) error
		GetAllWithTargets(context.Context) ([]Mission, error)
		GetByIDWithTargets(context.Context, int64) (*Mission, error)
		GetAllMissionTargets(context.Context, int64) ([]Target, error)
		HasAssignedSpy(context.Context, int64) (bool, error)
		GetTargetsQuantity(context.Context, int64) (int, error)
		GetTargetByID(context.Context, int64) (*Target, error)
		UpdateTarget(context.Context, *Target) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Cat:     &CatStore{db},
		Mission: &MissionStore{db},
	}
}

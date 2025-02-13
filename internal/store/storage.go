package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrorNotFound = errors.New("store: resource not found")
var QueryTimeoutDuration = 5 * time.Second

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
	}
	Mission interface {
		CRUD[Mission]
		AssignCat(context.Context, int64) error
		RemoveCat(context.Context, int64) error
		AddTarget(context.Context, *Target) error
		RemoveTarget(context.Context, int64) error
		AddNote(context.Context, int64, string) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Cat: &CatStore{db},
	}
}

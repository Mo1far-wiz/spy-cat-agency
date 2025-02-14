package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Cat struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	YearsOfExperience int     `json:"years_of_experience"`
	Breed             string  `json:"breed"`
	Salary            float64 `json:"salary"`
}

type CatStore struct {
	db *sql.DB
}

func (cs *CatStore) Create(ctx context.Context, cat *Cat) error {
	query := `
		INSERT INTO cats (name, years_of_experience, breed, salary)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := cs.db.QueryRowContext(
		ctx,
		query,
		cat.Name,
		cat.YearsOfExperience,
		cat.Breed,
		cat.Salary,
	).Scan(&cat.ID)

	if err != nil {
		return fmt.Errorf("store: failed to create cat: %w", err)
	}

	return nil
}

func (cs *CatStore) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM cats
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := cs.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("store: failed to delete cat: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: failed to retrieve affected rows: %w", err)
	}

	if rows == 0 {
		return ErrorNotFound
	}

	return nil
}

func (cs *CatStore) Update(ctx context.Context, cat *Cat) error {
	query := `
	UPDATE cats
	SET salary = $1
	WHERE id = $2;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := cs.db.ExecContext(
		ctx,
		query,
		cat.Salary,
		cat.ID,
	)

	if err != nil {
		return fmt.Errorf("store: failed to update cat salary: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: failed to retrieve affected rows: %w", err)
	}

	if affected == 0 {
		return ErrorNotFound
	}

	return nil
}

func (cs *CatStore) GetByID(ctx context.Context, id int64) (*Cat, error) {
	query := `
	SELECT id, name, years_of_experience, breed, salary
	FROM cats
	WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var cat Cat

	err := cs.db.QueryRowContext(ctx, query, id).
		Scan(
			&cat.ID,
			&cat.Name,
			&cat.YearsOfExperience,
			&cat.Breed,
			&cat.Salary,
		)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, fmt.Errorf("store: failed to retrieve cat: %w", err)
		}
	}

	return &cat, nil
}

func (cs *CatStore) GetAll(ctx context.Context) ([]Cat, error) {
	query := `
		SELECT id, name, years_of_experience, breed, salary
		FROM cats;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := cs.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("store: failed to execute query: %w", err)
	}
	defer rows.Close()

	cats := []Cat{}
	for rows.Next() {
		var c Cat
		err = rows.Scan(&c.ID, &c.Name, &c.YearsOfExperience, &c.Breed, &c.Salary)
		if err != nil {
			return nil, fmt.Errorf("store: failed to scan row: %w", err)
		}

		cats = append(cats, c)
	}
	return cats, nil
}

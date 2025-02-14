package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Target struct {
	ID         int64  `json:"id"`
	MissionID  int64  `json:"mission_id"`
	Name       string `json:"name"`
	Country    string `json:"country"`
	IsComplete bool   `json:"is_complete"`
}

func (ms *MissionStore) GetTargetByID(ctx context.Context, id int64) (*Target, error) {
	query := `
	SELECT id, mission_id, name, country, is_complete
	FROM targets
	WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var target Target
	err := ms.db.QueryRowContext(ctx, query, id).
		Scan(
			&target.ID,
			&target.MissionID,
			&target.Name,
			&target.Country,
			&target.IsComplete,
		)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, fmt.Errorf("store: failed to retrieve mission: %w", err)
		}
	}

	return &target, nil
}

func (ms *MissionStore) GetAllMissionTargets(ctx context.Context, missionID int64) ([]Target, error) {
	query := `
		SELECT id, mission_id, name, country, is_complete
		FROM targets
		WHERE mission_id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := ms.db.QueryContext(ctx, query, missionID)
	if err != nil {
		return nil, fmt.Errorf("store: failed to execute query: %w", err)
	}
	defer rows.Close()

	targets := []Target{}
	for rows.Next() {
		var t Target
		err = rows.Scan(&t.ID, &t.MissionID, &t.Name, &t.Country, &t.IsComplete)
		if err != nil {
			return nil, fmt.Errorf("store: failed to scan row: %w", err)
		}

		targets = append(targets, t)
	}
	return targets, nil
}

func (ms *MissionStore) GetTargetsQuantity(ctx context.Context, missionID int64) (int, error) {
	query := `
		SELECT COUNT(id)
		FROM targets
		WHERE mission_id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var quantity int
	err := ms.db.QueryRowContext(ctx, query, missionID).Scan(&quantity)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, ErrorNotFound
		default:
			return 0, fmt.Errorf("store: failed to retrieve targets: %w", err)
		}
	}

	return quantity, nil
}

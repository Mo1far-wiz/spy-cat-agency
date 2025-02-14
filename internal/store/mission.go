package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Mission struct {
	ID         int64    `json:"id"`
	CatID      *int64   `json:"cat_id"`
	IsComplete bool     `json:"is_complete"`
	Targets    []Target `json:"targets" validation:"required"`
}

type MissionStore struct {
	db *sql.DB
}

func (ms *MissionStore) Create(ctx context.Context, mission *Mission) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	tx, err := ms.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	queryCreateMission := `
		INSERT INTO missions (cat_id, is_complete)
		Values (NULL, false)
		RETURNING id;
	`

	err = tx.QueryRowContext(
		ctx,
		queryCreateMission,
	).Scan(&mission.ID)

	if err != nil {
		return fmt.Errorf("store: failed to create mission: %w", err)
	}

	queryCreateTargets := `
		INSERT INTO targets (mission_id, name, country, is_complete)
		Values ($1, $2, $3, false)
		RETURNING id;
	`
	for idx, t := range mission.Targets {
		var id int64
		err = tx.QueryRowContext(
			ctx,
			queryCreateTargets,
			mission.ID,
			t.Name,
			t.Country,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("store: failed to create target: %w", err)
		}

		mission.Targets[idx].ID = id
		mission.Targets[idx].MissionID = mission.ID
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("store: failed to commit transaction: %w", err)
	}

	return nil
}

// update that doesn't updates but rather changes mission marking
func (ms *MissionStore) Update(ctx context.Context, mission *Mission) error {
	query := `
	UPDATE missions
	SET is_complete = $1
	WHERE id = $2;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := ms.db.ExecContext(
		ctx,
		query,
		mission.IsComplete,
		mission.ID,
	)

	if err != nil {
		return fmt.Errorf("store: failed to update mission: %w", err)
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

func (ms *MissionStore) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM missions
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := ms.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("store: failed to delete mission: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: failed to retrieve affected rows: %w", err)
	}

	if rows == 0 {
		return ErrorNotFound
	}

	// on delete cascade will do the thing with targets and notes

	return nil
}

func (ms *MissionStore) GetByID(ctx context.Context, id int64) (*Mission, error) {
	query := `
	SELECT id, cat_id, is_complete
	FROM missions
	WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var mission Mission
	err := ms.db.QueryRowContext(ctx, query, id).
		Scan(
			&mission.ID,
			&mission.CatID,
			&mission.IsComplete,
		)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, fmt.Errorf("store: failed to retrieve mission: %w", err)
		}
	}

	return &mission, nil
}

func (ms *MissionStore) GetByIDWithTargets(ctx context.Context, id int64) (*Mission, error) {
	mission, err := ms.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("store: failed to get mission by ID: %w", err)
	}

	targets, err := ms.GetAllMissionTargets(ctx, mission.ID)
	if err != nil {
		return nil, fmt.Errorf("store: failed to get targets for mission %d: %w", mission.ID, err)
	}

	mission.Targets = targets

	return mission, nil
}

func (ms *MissionStore) GetAll(ctx context.Context) ([]Mission, error) {
	query := `
		SELECT id, cat_id, is_complete
		FROM missions;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := ms.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("store: failed to get missions: %w", err)
	}
	defer rows.Close()

	missions := []Mission{}
	for rows.Next() {
		var m Mission
		err = rows.Scan(&m.ID, &m.CatID, &m.IsComplete)
		if err != nil {
			return nil, fmt.Errorf("store: failed to scan row: %w", err)
		}

		missions = append(missions, m)
	}
	return missions, nil
}

func (ms *MissionStore) GetAllWithTargets(ctx context.Context) ([]Mission, error) {
	missions, err := ms.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("store: failed to get missions: %w", err)
	}

	for i, mission := range missions {
		targets, err := ms.GetAllMissionTargets(ctx, mission.ID)
		if err != nil {
			return nil, fmt.Errorf("store: failed to get targets for mission %d: %w", mission.ID, err)
		}
		missions[i].Targets = targets
	}

	return missions, nil
}

func (ms *MissionStore) AssignCat(ctx context.Context, catID int64, missionID int64) error {
	query := `
		UPDATE missions
		SET cat_id = $1
		WHERE id = $2;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := ms.db.ExecContext(
		ctx,
		query,
		catID,
		missionID,
	)

	if err != nil {
		return fmt.Errorf("store: failed to assign cat: %w", err)
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

func (ms *MissionStore) HasAssignedSpy(ctx context.Context, missionID int64) (bool, error) {
	query := `
		SELECT cat_id
		FROM missions
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var catID *int64
	err := ms.db.QueryRowContext(ctx, query, missionID).Scan(&catID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return false, ErrorNotFound
		default:
			return false, fmt.Errorf("store: failed to check if mission has assigned spy: %w", err)
		}
	}

	return catID != nil, nil
}

func (ms *MissionStore) AddTarget(ctx context.Context, id int64, target *Target) error {
	query := `
		INSERT INTO targets (mission_id, name, country, is_complete)
		VALUES ($1, $2, $3, false)
		RETURNING id;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := ms.db.QueryRowContext(
		ctx,
		query,
		id,
		target.Name,
		target.Country,
	).Scan(&target.ID)

	if err != nil {
		return fmt.Errorf("store: failed to add target: %w", err)
	}

	return nil
}

func (ms *MissionStore) RemoveTarget(ctx context.Context, targetId int64) error {
	query := `
		DELETE FROM targets
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := ms.db.ExecContext(ctx, query, targetId)
	if err != nil {
		return fmt.Errorf("store: failed to remove target: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: failed to retrieve affected rows: %w", err)
	}

	if rows == 0 {
		return ErrorNotFound
	}

	// on delete cascade will do the thing with notes

	return nil
}

func (ms *MissionStore) UpdateTarget(ctx context.Context, target *Target) error {
	query := `
	UPDATE targets
	SET is_complete = $1
	WHERE id = $2;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := ms.db.ExecContext(
		ctx,
		query,
		target.IsComplete,
		target.ID,
	)

	if err != nil {
		return fmt.Errorf("store: failed to update target: %w", err)
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

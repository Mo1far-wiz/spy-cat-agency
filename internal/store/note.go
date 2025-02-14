package store

import (
	"context"
	"fmt"
	"time"
)

type Note struct {
	ID int64 `json:"id"`
	// a little denormalization
	MissionID int64  `json:"mission_id"`
	TargetID  int64  `json:"target_id"`
	Note      string `json:"note"`
	// i know that it wasn't in task, but it's just makes sense
	CreatedAt time.Time `json:"created_at"`
}

func (ms *MissionStore) AddNote(ctx context.Context, note *Note) error {
	query := `
		INSERT INTO notes (mission_id, target_id, note)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := ms.db.QueryRowContext(
		ctx,
		query,
		note.MissionID,
		note.TargetID,
		note.Note,
	).Scan(&note.ID)

	if err != nil {
		return fmt.Errorf("store: failed to create note: %w", err)
	}

	return nil
}

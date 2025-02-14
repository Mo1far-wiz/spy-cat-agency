CREATE TABLE notes (
    id bigserial PRIMARY KEY,
    mission_id BIGINT NOT NULL REFERENCES missions(id),
    target_id BIGINT NOT NULL REFERENCES targets(id) ON DELETE CASCADE,
    note TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

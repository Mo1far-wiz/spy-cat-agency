CREATE TABLE notes (
    id bigserial PRIMARY KEY,
    target_id BIGINT NOT NULL REFERENCES targets(id) ON DELETE CASCADE,
    note TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE targets (
    id bigserial PRIMARY KEY,
    mission_id BIGINT NOT NULL REFERENCES missions(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    is_complete BOOLEAN NOT NULL DEFAULT FALSE
);
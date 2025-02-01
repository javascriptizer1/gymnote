-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS exercises (
    id UUID,
    name String,
    muscle_group String,
    equipment String,
    created_at DateTime,
) ENGINE = MergeTree()
  PRIMARY KEY id
  ORDER BY (id, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exercises;
-- +goose StatementEnd

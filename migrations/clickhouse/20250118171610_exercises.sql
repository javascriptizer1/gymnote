-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS exercises (
    id UUID,
    name String,
    muscle_group String,
    equipment String,
    createdAt DateTime,
) ENGINE = MergeTree()
  PRIMARY KEY id
  ORDER BY (id, createdAt);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exercises;
-- +goose StatementEnd

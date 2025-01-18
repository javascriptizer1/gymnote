-- +goose Up
-- +goose StatementBegin
CREATE TABLE training_sessions (
    id UUID,
    user_id String,
    date Date,
    total_volume Float32,
    exercise_count UInt8,
    notes String,
) ENGINE = MergeTree()
  PRIMARY KEY id
  ORDER BY (id, date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS training_sessions;
-- +goose StatementEnd

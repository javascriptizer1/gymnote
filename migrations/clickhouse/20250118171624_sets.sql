-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS training_logs (
    id UUID,
    user_id String,
    session_date Date,
    exercise_id UUID,
    set_number UInt8,
    weight Float32,
    reps UInt8,
    difficulty String,
    notes String,
) ENGINE = MergeTree()
  PRIMARY KEY id
  ORDER BY (id, session_date, exercise_id, set_number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS training_logs;
-- +goose StatementEnd

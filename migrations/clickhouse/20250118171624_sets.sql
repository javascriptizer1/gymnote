-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS training_logs (
    id UUID,
    user_id String,
    session_id UUID,
    exercise_id UUID,
    session_date Date,
    exercise_name String,
    exercise_number UInt8,
    set_number UInt8,
    weight Float32,
    reps UInt8,
    difficulty String,
    notes String,
    muscle_group String,
    created_at DateTime
) ENGINE = MergeTree()
  PRIMARY KEY (id, session_id, exercise_id, exercise_number, set_number, created_at)
  ORDER BY (id, session_id, exercise_id, exercise_number, set_number, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS training_logs;
-- +goose StatementEnd

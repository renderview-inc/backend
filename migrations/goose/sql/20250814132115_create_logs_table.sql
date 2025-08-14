-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS logs (
    timestamp DateTime DEFAULT now(),
    level Enum('debug'=1, 'info'=2, 'warn'=3, 'error'=4, 'dpanic'=5, 'panic'=6, 'fatal'=7),
    msg String,
    fields JSON,
    correlation_id String
) ENGINE = MergeTree()
ORDER BY (timestamp, correlation_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS logs;
-- +goose StatementEnd

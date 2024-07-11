-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    chat_id BIGINT PRIMARY KEY,
    binance_name varchar,
    bybit_name varchar
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd

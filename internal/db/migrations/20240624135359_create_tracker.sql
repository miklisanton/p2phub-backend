-- +goose Up
-- +goose StatementBegin
CREATE TABLE trackers (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    exchange varchar NOT NULL,
    currency varchar(3) NOT NULL,
    side varchar(4) NOT NULL,
    waiting_adv boolean DEFAULT false,
    UNIQUE (user_id, exchange, currency, side),
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(chat_id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE trackers;
-- +goose StatementEnd

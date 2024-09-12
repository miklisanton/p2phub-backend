-- +goose Up
-- +goose StatementBegin
CREATE TABLE trackers (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    exchange varchar NOT NULL,
    currency varchar(3) NOT NULL,
    username varchar NOT NULL,
    side varchar(4) NOT NULL,
    waiting_adv boolean DEFAULT false,
    outbided boolean DEFAULT false,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE trackers;
-- +goose StatementEnd

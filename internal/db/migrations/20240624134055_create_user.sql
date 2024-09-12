-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT UNIQUE, 
    email varchar UNIQUE,
    password_enc varchar,
    CHECK (chat_id IS NOT NULL OR email IS NOT NULL),
    CHECK (email IS NULL OR password_enc IS NOT NULL)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd

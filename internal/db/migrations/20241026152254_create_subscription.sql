-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE subscription (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    valid_until TIMESTAMP DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 month',
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
    
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE subscription;
-- +goose StatementEnd

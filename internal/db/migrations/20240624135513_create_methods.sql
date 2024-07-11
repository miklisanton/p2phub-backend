-- +goose Up
-- +goose StatementBegin
CREATE TABLE methods (
    tracker_id INT NOT NULL,
    payment_method  varchar(64),
    UNIQUE (tracker_id, payment_method),
    CONSTRAINT fk_tracker
        FOREIGN KEY (tracker_id)
        REFERENCES trackers(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE methods;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE addresses (
    address_id SERIAL PRIMARY KEY,
    street VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    country VARCHAR(100) NOT NULL DEFAULT 'ETHIOPIA'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE addresses;
-- +goose StatementEnd
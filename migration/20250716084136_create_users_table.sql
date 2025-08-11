-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists users (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    Hash_Password VARCHAR(250) NOT NULL,
    address_id INT,
    CONSTRAINT fk_address FOREIGN KEY (address_id) REFERENCES addresses (address_id) ON DELETE RESTRICT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
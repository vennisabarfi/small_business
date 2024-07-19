-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    email VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- 'Down' is executed when migration is rolled back
-- +goose StatementBegin
DROP TABLE users
-- +goose StatementEnd

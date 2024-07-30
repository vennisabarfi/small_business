-- +goose Up
-- +goose StatementBegin
CREATE TABLE supplier (
	id   INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	name	VARCHAR(255) UNIQUE NOT NULL,
	contact_email VARCHAR(255) UNIQUE,
    phone VARCHAR(20),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE suppliers;

-- +goose StatementEnd

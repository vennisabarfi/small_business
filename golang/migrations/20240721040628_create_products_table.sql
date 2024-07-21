-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
       id SERIAL PRIMARY KEY,
       name VARCHAR UNIQUE NOT NULL,
       description VARCHAR,
       price DECIMAL(6,2) NOT NULL,
       stock INTEGER NOT NULL,
       supplier_id INTEGER,
       FOREIGN KEY(supplier_id) 
        REFERENCES supplier(id)
        ON DELETE CASCADE,
       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
   );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP 
-- +goose StatementEnd

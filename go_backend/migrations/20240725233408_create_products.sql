-- +goose Up
-- +goose StatementBegin
CREATE TABLE products(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	name VARCHAR UNIQUE NOT NULL,
	description VARCHAR, 
    supplier_id INT NOT NULL,
    price DECIMAL(6,2) NOT NULL,
    stock INTEGER NOT NULL,
	minimum_stock SMALLINT  NOT NULL,
	CONSTRAINT fk_supplier FOREIGN KEY(supplier_id) REFERENCES "supplier"(id),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd

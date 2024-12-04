-- Classes Table
-- ONE TO MANY in relaton to Ascendancies

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL, -- Description of the class
    image_url TEXT, -- URL for class image
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS classes
-- +goose StatementEnd

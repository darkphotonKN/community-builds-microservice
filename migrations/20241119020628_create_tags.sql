-- Build Tags Table
-- Join Table between Tags and Builds

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create function to auto-update the updated_at column
CREATE OR REPLACE FUNCTION update_tags_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for the tags table
CREATE TRIGGER set_tags_updated_at
BEFORE UPDATE ON tags
FOR EACH ROW
EXECUTE FUNCTION update_tags_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop trigger and function
DROP TRIGGER IF EXISTS set_tags_updated_at ON tags;
DROP FUNCTION IF EXISTS update_tags_updated_at;

-- Drop table
DROP TABLE IF EXISTS tags;
-- +goose StatementEnd

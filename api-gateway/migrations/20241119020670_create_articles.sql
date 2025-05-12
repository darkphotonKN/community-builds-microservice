-- Build articles Table
-- Join Table between articles and Builds

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS articles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create function to auto-update the updated_at column
CREATE OR REPLACE FUNCTION update_articles_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for the articles table
CREATE TRIGGER set_articles_updated_at
BEFORE UPDATE ON articles
FOR EACH ROW
EXECUTE FUNCTION update_articles_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop trigger and function
DROP TRIGGER IF EXISTS set_articles_updated_at ON articles;
DROP FUNCTION IF EXISTS update_articles_updated_at;

-- Drop table
DROP TABLE IF EXISTS articles;
-- +goose StatementEnd

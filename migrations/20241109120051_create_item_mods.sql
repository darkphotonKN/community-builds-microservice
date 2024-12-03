-- Items Table
-- Many to Many with Builds

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS item_mods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    affix TEXT NOT NULL,
    name TEXT NOT NULL,
	level TEXT NOT NULL,
	stat TEXT NOT NULL,
	tags TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS item_mods;
-- +goose StatementEnd

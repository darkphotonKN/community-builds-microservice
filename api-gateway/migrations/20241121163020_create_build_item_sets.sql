-- Build Item Links Table
-- - Many-to-One: Relates to `builds` (each group belongs to one build).
-- - One-to-Many: Related to `build_items` (each group contains multiple items).
-- - Many-to-Many: Indirectly links to `items` via the `build_items` join table.

-- +goose Up
-- +goose StatementBegin
CREATE TABLE build_item_sets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    build_id UUID NOT NULL REFERENCES builds(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS build_item_sets
-- +goose StatementEnd

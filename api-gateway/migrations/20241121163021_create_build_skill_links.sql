-- Build Skill Links Table
-- - Many-to-One: Relates to `builds` (each group belongs to one build).
-- - One-to-Many: Related to `build_skills` (each group contains multiple skills).
-- - Many-to-Many: Indirectly links to `skills` via the `build_skills` join table.

-- +goose Up
-- +goose StatementBegin
CREATE TABLE build_skill_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    build_id UUID NOT NULL REFERENCES builds(id) ON DELETE CASCADE,  
    name TEXT NOT NULL, -- Name of  (e.g., "Main DPS", "Mobility", "defensive")
    is_main BOOLEAN DEFAULT FALSE, -- Indicates the primary skill group for the build
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS build_skill_links
-- +goose StatementEnd

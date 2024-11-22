-- Build Skill Groups Table
-- Many to One in relation with builds
-- One to Many build_skills
-- Many to Many to Skills via join table build_skills

-- +goose Up
-- +goose StatementBegin
CREATE TABLE build_skill_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    build_id UUID NOT NULL REFERENCES builds(id) ON DELETE CASCADE,  
    name TEXT NOT NULL, -- Group name (e.g., "Main DPS", "Mobility", "defensive")
    is_main BOOLEAN DEFAULT FALSE, -- Indicates the primary skill group for the build
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP IF EXISTS build_skill_groups
-- +goose StatementEnd

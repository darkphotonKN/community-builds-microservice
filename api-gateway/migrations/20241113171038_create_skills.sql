-- Skills Table
-- Many to Many in relation with Builds (join table build_skills)

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS skills (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,  -- Skill name (e.g., "Lightning Strike")
    type TEXT NOT NULL CHECK  (type IN ('active', 'support')), -- 1: Active or 2: Support skill 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS skills;
-- +goose StatementEnd

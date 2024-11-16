-- Builds Table
-- Many to One in relation with Members
-- Many to One relation with Skills
-- Many to Many in relation with Items (join table build_items)
-- Many to Many in relation with Skills (join table build_skills)

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS builds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    member_id UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE, -- FK to members table
    title TEXT NOT NULL, -- main build title
    description TEXT NOT NULL,
    main_skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE RESTRICT, -- FK to skills table
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS builds;
-- +goose StatementEnd


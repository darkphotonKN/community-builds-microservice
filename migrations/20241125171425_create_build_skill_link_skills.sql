-- Build Skill Link Skills Join Table
-- Join Table between build_skill_links with skills

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS build_skill_link_skills (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    build_skill_link_id UUID NOT NULL REFERENCES build_skill_links(id) ON DELETE CASCADE,
    skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS build_skill_link_skills;
-- +goose StatementEnd

-- Build Tags Table
-- Join Table between Tags and Builds

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS build_tag (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    build_id UUID NOT NULL REFERENCES builds(id) ON DELETE CASCADE,
    -- not repeatable
    UNIQUE(build_id, tag_id)
);
-- optimize index 
CREATE INDEX idx_build_tag_build_id ON build_tag(build_id);
CREATE INDEX idx_build_tag_tag_id ON build_tag(tag_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS build_tag;
-- +goose StatementEnd

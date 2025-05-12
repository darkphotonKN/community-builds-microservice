-- Build Item Link Items Join Table
-- Join Table between build_item_links with items

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS build_item_set_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    build_item_set_id UUID NOT NULL REFERENCES build_item_sets(id) ON DELETE CASCADE,
    item_id UUID REFERENCES items(id) ON DELETE CASCADE,
    slot TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS build_item_set_items;
-- +goose StatementEnd

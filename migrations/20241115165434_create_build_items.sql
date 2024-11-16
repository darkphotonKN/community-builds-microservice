-- Build Items Table
-- Join Table between Items and Builds

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS build_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    build_id UUID NOT NULL REFERENCES builds(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    slot TEXT NOT NULL -- Slot type (Helm, Chest, Boots, Belt, Gloves, WeaponOne, WeaponTwo, RingOne, RingTwo, Necklace)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS build_items;
-- +goose StatementEnd

-- Items Table
-- Many to Many with Builds

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category TEXT NOT NULL CHECK (category <> ''), -- E.g., "Two Handed Weapon", "Gems"
    class TEXT NOT NULL CHECK (class <> ''), -- Class of Item, below category. E.g., "Body Armours", "Two Hand Swords"
    type TEXT NOT NULL CHECK (type <> ''), -- Type of Item, below class. E.g., "Glourious Plate"
    name TEXT NOT NULL UNIQUE CHECK (name <> ''), -- Name of the item (e.g., "Kaom's Heart")
    image_url TEXT, -- Path or URL for the item's image
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
-- +goose StatementEnd



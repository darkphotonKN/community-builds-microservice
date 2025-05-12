-- Items Table
-- Many to Many with Builds

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    member_id UUID REFERENCES members(id) ON DELETE RESTRICT,
    base_item_id UUID REFERENCES base_items(id) ON DELETE RESTRICT,
    -- category TEXT NOT NULL CHECK (category <> ''), -- E.g., "Two Handed Weapon", "Gems"
    -- class TEXT NOT NULL CHECK (class <> ''), -- Class of Item, below category. E.g., "Body Armours", "Two Hand Swords"
    -- type TEXT NOT NULL CHECK (type <> ''), -- Type of Item, below class. E.g., "Glourious Plate"
    -- name TEXT NOT NULL CHECK (name <> ''), -- Name of the item (e.g., "Kaom's Heart")
    category TEXT  NOT NULL,
    class TEXT  NOT NULL,
    type TEXT NOT NULL,
    name TEXT NOT NULL,
    unique_item BOOLEAN NOT NULL,
    slot TEXT NOT NULL,
    description TEXT, -- equip's story or description
    image_url TEXT, -- Path or URL for the item's image
    
    required_level TEXT,
    required_strength TEXT,
    required_dexterity TEXT,
    required_intelligence TEXT,
    armour TEXT,
    block TEXT,
    energy_shield TEXT,
    evasion TEXT,
    ward TEXT,

    damage TEXT,
    aps TEXT,
    crit TEXT,
    pdps TEXT,
    edps TEXT,
    dps TEXT,

    life TEXT,
    mana TEXT,
    duration TEXT,
    usage TEXT,
    capacity TEXT,

    additional TEXT,
    stats TEXT[],
    implicit TEXT[],

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
-- +goose StatementEnd

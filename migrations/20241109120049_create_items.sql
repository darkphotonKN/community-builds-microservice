-- Items Table
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES members(id) ON DELETE CASCADE, -- Links item to the member listing it
    item_id UUID NOT NULL,  -- Unique identifier for each item type
    name VARCHAR(255) NOT NULL,  -- Name of the item
    description TEXT,  -- Description of the item
    price_per_unit DECIMAL(10, 2) NOT NULL,  -- Price for each unit
    stock_quantity INT NOT NULL CHECK (stock_quantity >= 0),  -- Number of items available for sale
    category VARCHAR(255) NOT NULL,  -- Broad Category of Item
    type VARCHAR(255) NOT NULL,  -- Specific Type of Item
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
-- +goose StatementEnd



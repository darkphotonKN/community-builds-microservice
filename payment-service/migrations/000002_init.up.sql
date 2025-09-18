CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- CREATE INDEX idx_members_email ON members(email);
-- CREATE INDEX idx_members_created_at ON members(created_at);
CREATE INDEX idx_products_created_at ON products(created_at);
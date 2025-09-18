ALTER TABLE members 
ADD COLUMN stripe_customer_id VARCHAR(255) UNIQUE;

CREATE INDEX idx_members_stripe_customer_id ON members(stripe_customer_id);
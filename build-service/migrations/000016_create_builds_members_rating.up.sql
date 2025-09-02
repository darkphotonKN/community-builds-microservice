CREATE TABLE IF NOT EXISTS builds_members_rating (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    to_member_id UUID NOT NULL,
    from_member_id UUID NOT NULL,
    build_id UUID NOT NULL REFERENCES builds(id) ON DELETE CASCADE,
    rating INT NOT NULL CHECK (rating >= 0 AND rating <= 5)
); 
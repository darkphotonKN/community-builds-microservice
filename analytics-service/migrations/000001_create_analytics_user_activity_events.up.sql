
-- Create member_activity_events table
CREATE TABLE member_activity_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for common queries
CREATE INDEX idx_member_activity_events_member_id ON member_activity_events(member_id);
CREATE INDEX idx_member_activity_events_date ON member_activity_events(date);
CREATE INDEX idx_member_activity_events_event_type ON member_activity_events(event_type);
CREATE INDEX idx_member_activity_events_timestamp ON member_activity_events(timestamp);

-- Create daily_member_stats table
CREATE TABLE daily_member_stats (
    date DATE PRIMARY KEY,
    new_signups INTEGER DEFAULT 0,
    active_members INTEGER DEFAULT 0,
    total_logins INTEGER DEFAULT 0,
    total_members_to_date INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for date range queries
CREATE INDEX idx_daily_member_stats_date ON daily_member_stats(date);

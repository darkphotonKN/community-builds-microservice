-- notification-service/migrations/000001_create_notifications.up.sql

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Simple notifications table for web consumption
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Who gets this notification
    member_id UUID NOT NULL,  -- Member ID from auth service
    
    -- Notification content
    type VARCHAR(50) NOT NULL,     -- 'welcome', 'build_commented', 'build_rated', 'build_liked'
    title VARCHAR(200) NOT NULL,   -- "Your build received a comment!"
    message VARCHAR(500) NOT NULL, -- "Sarah commented on your Lightning Sorceress build"
    
    -- Web-specific fields
    read BOOLEAN DEFAULT false,    -- Simple read/unread status
    
    -- Simple status
    email_sent BOOLEAN DEFAULT false,
    
    -- Metadata
    source_id UUID,                -- ID of the build/rating/etc that triggered this
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for web queries
CREATE INDEX idx_notifications_member_unread ON notifications(member_id) WHERE read = false;
CREATE INDEX idx_notifications_member_recent ON notifications(member_id, created_at DESC);
CREATE INDEX idx_notifications_type ON notifications(type);

-- Web notification types for community builds
COMMENT ON COLUMN notifications.type IS 'Types: welcome, build_commented, build_rated, build_liked';

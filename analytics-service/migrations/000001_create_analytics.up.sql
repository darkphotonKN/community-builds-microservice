-- analytics-service/migrations/000001_create_analytics.up.sql

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Analytics table for tracking user behavior and events
CREATE TABLE IF NOT EXISTS analytics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Who performed this action
    member_id UUID NOT NULL,  -- Member ID from auth service
    
    -- Event details
    event_type VARCHAR(50) NOT NULL,     -- 'user_activity', 'system_event', 'performance'
    event_name VARCHAR(100) NOT NULL,    -- 'member_signup', 'build_created', 'page_view'
    data JSONB,                          -- Additional event data as JSON
    
    -- Session and context
    session_id UUID,                     -- Session identifier
    ip_address INET,                     -- User's IP address
    user_agent TEXT,                     -- Browser/client information
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for analytics queries
CREATE INDEX idx_analytics_member_id ON analytics(member_id);
CREATE INDEX idx_analytics_event_type ON analytics(event_type);
CREATE INDEX idx_analytics_event_name ON analytics(event_name);
CREATE INDEX idx_analytics_created_at ON analytics(created_at DESC);
CREATE INDEX idx_analytics_session_id ON analytics(session_id) WHERE session_id IS NOT NULL;

-- Analytics event types for community builds
COMMENT ON COLUMN analytics.event_type IS 'Types: user_activity, system_event, performance';
COMMENT ON COLUMN analytics.event_name IS 'Examples: member_signup, build_created, page_view, search_performed';

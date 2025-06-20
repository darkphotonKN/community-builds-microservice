-- analytics-service/migrations/000001_create_analytics.down.sql

-- Drop indexes first
DROP INDEX IF EXISTS idx_daily_user_stats_date;
DROP INDEX IF EXISTS idx_user_activity_events_timestamp;
DROP INDEX IF EXISTS idx_user_activity_events_event_type;
DROP INDEX IF EXISTS idx_user_activity_events_date;
DROP INDEX IF EXISTS idx_user_activity_events_user_id;

-- Drop tables
DROP TABLE IF EXISTS daily_user_stats;
DROP TABLE IF EXISTS user_activity_events;

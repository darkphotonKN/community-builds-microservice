-- Drop indexes first
DROP INDEX IF EXISTS idx_daily_member_stats_date;
DROP INDEX IF EXISTS idx_member_activity_events_timestamp;
DROP INDEX IF EXISTS idx_member_activity_events_activity_type;
DROP INDEX IF EXISTS idx_member_activity_events_date;
DROP INDEX IF EXISTS idx_member_activity_events_member_id;

-- Drop tables
DROP TABLE IF EXISTS daily_member_stats;
DROP TABLE IF EXISTS member_activity_events;

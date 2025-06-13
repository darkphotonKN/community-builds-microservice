-- analytics-service/migrations/000001_create_analytics.down.sql

-- Drop indexes
DROP INDEX IF EXISTS idx_analytics_member_id;
DROP INDEX IF EXISTS idx_analytics_event_type;
DROP INDEX IF EXISTS idx_analytics_event_name;
DROP INDEX IF EXISTS idx_analytics_created_at;
DROP INDEX IF EXISTS idx_analytics_session_id;

-- Drop table
DROP TABLE IF EXISTS analytics;

-- Drop extension (be careful - other services might use it)
DROP EXTENSION IF EXISTS "uuid-ossp";

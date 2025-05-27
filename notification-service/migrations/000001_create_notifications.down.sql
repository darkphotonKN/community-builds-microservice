-- notification-service/migrations/000001_create_notifications.down.sql

-- Drop indexes
DROP INDEX IF EXISTS idx_notifications_member_unread;
DROP INDEX IF EXISTS idx_notifications_member_recent;
DROP INDEX IF EXISTS idx_notifications_type;

-- Drop table
DROP TABLE IF EXISTS notifications;

-- Drop extension (be careful - other services might use it)
-- DROP EXTENSION IF EXISTS "uuid-ossp";

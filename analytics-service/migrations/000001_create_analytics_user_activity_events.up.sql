/*
User Stories:
As a Product Owner:

I want to see how many users signed up this week so I can track growth
I want to know how many users are active daily so I can measure engagement
I want to see which days users are most active so I can plan feature releases
I want to track user retention - do users come back after 1 week, 1 month?

As a Community Manager:

I want to see signup trends over time so I can correlate with marketing efforts
I want to know when users are most active so I can schedule announcements
I want to identify power users (login frequently) vs casual users

As a Developer:

I want to track login success vs failure rates to monitor system health
I want to see geographic patterns of where users are signing up from
I want to know which features trigger the most user activity

Event Flow Stories:
"User Logs In" Story:

User enters credentials in frontend
Auth service validates login
Auth service publishes "user logged in" event
Analytics service captures login time, user ID, location
Analytics service updates "daily active users" counter
Dashboard shows real-time active user count

"Track User Retention" Story:

User signs up (existing event you have)
Analytics service marks user as "new user"
Analytics service tracks when same user logs in again
Calculate: "Did user return within 7 days?"
Display retention rates on admin dashboard

"Popular Times" Story:

Track all login events throughout the day
Aggregate by hour: "Most logins happen at 8pm"
Help you decide when to deploy features or send notifications
*/

-- Create user_activity_events table
CREATE TABLE user_activity_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for common queries
CREATE INDEX idx_user_activity_events_user_id ON user_activity_events(user_id);
CREATE INDEX idx_user_activity_events_date ON user_activity_events(date);
CREATE INDEX idx_user_activity_events_event_type ON user_activity_events(event_type);
CREATE INDEX idx_user_activity_events_timestamp ON user_activity_events(timestamp);

-- Create daily_user_stats table
CREATE TABLE daily_user_stats (
    date DATE PRIMARY KEY,
    new_signups INTEGER DEFAULT 0,
    active_users INTEGER DEFAULT 0,
    total_logins INTEGER DEFAULT 0,
    total_users_to_date INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for date range queries
CREATE INDEX idx_daily_user_stats_date ON daily_user_stats(date);

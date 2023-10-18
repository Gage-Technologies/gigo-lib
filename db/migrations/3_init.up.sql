-- delete all the rows in user_stats and user_daily_usage
delete from user_stats;
delete from user_daily_usage;

-- initialize the user_stats table with the user_id and the date
insert into user_stats (_id, user_id, challenges_completed, streak_active, current_streak, longest_streak, total_time_spent, avg_time, days_on_platform, days_on_fire, streak_freezes, streak_freeze_used, xp_gained, date, expiration)
    select
        -- we need a unique id that conforms to the snowflake standard so this is hacky way to preserve an approximate
        -- correlation to when the first user stats should have been created
        _id + FLOOR(RAND() * 1000) as _id,
        _id as user_id,
        0 as challenges_completed,
        false as streak_active,
        0 as current_streak,
        0 as longest_streak,
        0 as total_time_spent,
        0 as avg_time,
        0 as days_on_platform,
        0 as days_on_fire,
        0 as streak_freezes,
        0 as streak_freeze_used,
        0 as xp_gained,
        DATE(CONVERT_TZ(NOW(), @@session.time_zone, 'America/Los_Angeles')) as date,
        DATE_ADD(DATE(CONVERT_TZ(NOW(), @@session.time_zone, 'America/Los_Angeles')), INTERVAL 24 HOUR) as expiration
    from users;
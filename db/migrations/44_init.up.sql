ALTER TABLE user_inactivity MODIFY COLUMN notify_on datetime NULL;

INSERT INTO user_inactivity (user_id, last_login, last_notified, send_week, send_month, notify_on, email)
SELECT
    u._id AS user_id,
    IFNULL(wt.max_timestamp, u.created_at) AS last_login,
    NOW() - INTERVAL 31 DAY AS last_notified,  -- Assuming current time as the last notified time
    FALSE AS send_week,       -- Default values for send_week and send_month
    FALSE AS send_month,
    NULL AS notify_on,        -- Assuming NULL for notify_on, adjust as needed
    u.email AS email
FROM
    users u
        LEFT JOIN
    (SELECT user_id, MAX(timestamp) AS max_timestamp FROM web_tracking GROUP BY user_id) wt
    ON
            u._id = wt.user_id
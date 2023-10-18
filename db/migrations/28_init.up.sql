-- Add a json field tutorials to the users table
ALTER TABLE friend_requests ADD COLUMN notification_id bigint NOT NULL DEFAULT 0;
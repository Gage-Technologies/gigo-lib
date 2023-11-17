ALTER TABLE user_inactivity DROP COLUMN should_notify;
ALTER TABLE user_inactivity DROP COLUMN week_sent;
ALTER TABLE user_inactivity ADD COLUMN send_week boolean NOT NULL;
ALTER TABLE user_inactivity ADD COLUMN send_month boolean NOT NULL;
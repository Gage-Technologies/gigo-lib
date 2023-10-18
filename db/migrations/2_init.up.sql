-- Add expiration and closed to the user_stats table to enable exact expiration tracking
ALTER TABLE user_stats ADD COLUMN expiration timestamp not null;
ALTER TABLE user_stats ADD COLUMN closed boolean not null default false;
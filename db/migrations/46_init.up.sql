-- MySQL v5.7

-- Add referred_by bigint nullable to users table
ALTER TABLE users ADD COLUMN referred_by bigint;
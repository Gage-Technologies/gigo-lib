-- MySQL v5.7

-- Add used_free_trial boolean column to the users table with a default false
ALTER TABLE users ADD COLUMN used_free_trial boolean NOT NULL DEFAULT FALSE;

-- Update all existing users used_free_trial column to true
UPDATE users SET used_free_trial = true;
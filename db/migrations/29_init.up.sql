-- Add muted column to chat user table as a boolean that defaults to false
ALTER TABLE chat_users ADD COLUMN muted boolean DEFAULT false;
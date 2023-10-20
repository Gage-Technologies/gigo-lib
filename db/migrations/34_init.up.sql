-- Add is vnc column
ALTER TABLE workspace ADD COLUMN is_vnc bool NOT NULL DEFAULT false;
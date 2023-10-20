-- Add is vnc column
ALTER TABLE workspaces ADD COLUMN is_vnc bool NOT NULL DEFAULT false;
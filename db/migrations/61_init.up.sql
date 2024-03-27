-- MySQL v5.7

-- add last_state_update to workspace_pool
ALTER TABLE workspace_pool ADD COLUMN last_state_update TIMESTAMP;

-- add expiration to workspace_pool
ALTER TABLE workspace_pool ADD COLUMN expiration TIMESTAMP;
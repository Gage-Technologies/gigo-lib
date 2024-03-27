-- MySQL v5.7

-- add creation_start_timestamp to workspace_pool
ALTER TABLE workspace_pool ADD COLUMN create_start_timestamp TIMESTAMP;

-- add expiration to workspace_pool
ALTER TABLE workspace_pool ADD COLUMN expiration TIMESTAMP;
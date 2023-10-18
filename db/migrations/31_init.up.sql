-- Add a temp user field to users table
ALTER TABLE users ADD COLUMN is_ephemeral boolean DEFAULT false;
-- Add a temp workspace field to workspaces table
ALTER TABLE workspaces ADD COLUMN is_ephemeral boolean DEFAULT false;
-- Add table to track usage of temp workspaces
create table if not exists ephemeral_shared_workspaces (
   workspace_id BIGINT NOT NULL,
   ip BIGINT NOT NULL,
   date DATETIME NOT NULL,
   user_id BIGINT NOT NULL,
   challenge_id BIGINT NOT NULL,
   primary key (workspace_id, ip, challenge_id)
)
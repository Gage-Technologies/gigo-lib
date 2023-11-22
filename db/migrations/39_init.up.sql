-- MySQL v5.7

-- Add ziti_id and ziti_token columns to the workspace_agent table
ALTER TABLE workspace_agent ADD COLUMN ziti_id varchar(50);
ALTER TABLE workspace_agent ADD COLUMN ziti_token text;
create table if not exists user_inactivity(
      user_id bigint not null primary key,
      last_login datetime not null,
      last_notified datetime not null,
      should_notify boolean not null
)
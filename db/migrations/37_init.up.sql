-- Mysql v5.7

-- Add start_time bigint column to workspace table

alter table `workspaces` add column `start_time` bigint;

-- Add start_time bigint column to attempt table

alter table `attempt` add column `start_time` bigint;

-- Add start_time bigint column to post table

alter table `post` add column `start_time` bigint;
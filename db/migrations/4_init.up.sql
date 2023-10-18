-- Add deleted to post table for deleting post
ALTER TABLE post ADD COLUMN deleted boolean not null default false;
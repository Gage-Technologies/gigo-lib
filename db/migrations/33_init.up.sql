-- Add share_hash column to post table
ALTER TABLE post ADD COLUMN share_hash binary(16);

-- Add storage class column to volpool_volume table
ALTER TABLE volpool_volume ADD COLUMN storage_class varchar(255);
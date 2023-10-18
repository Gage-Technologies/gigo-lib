-- Add views column to recommended_post table
ALTER TABLE recommended_post ADD COLUMN views BIGINT NOT NULL DEFAULT 0;
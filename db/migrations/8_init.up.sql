-- Fix bad unique index on recommendation table
ALTER TABLE recommended_post DROP INDEX IF EXISTS idx_recommended_post_user_id;
ALTER TABLE recommended_post DROP INDEX IF EXISTS idx_recommended_post_post_id;
ALTER TABLE recommended_post ADD CONSTRAINT uk_recommended_p_post_i_user_i UNIQUE (post_id, user_id);
-- Add post_type to attempts
ALTER TABLE attempt ADD COLUMN post_type int NOT NULL DEFAULT 0;

-- update attempts with the post type of their parent posts
UPDATE attempt a
    INNER JOIN post p ON a.post_id = p._id
SET a.post_type = p.post_type;

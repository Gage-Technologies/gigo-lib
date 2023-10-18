-- Add accepted boolean to recommended_post table
ALTER TABLE recommended_post ADD COLUMN accepted boolean NOT NULL DEFAULT false;

-- Add xp_reason table
CREATE TABLE IF NOT EXISTS xp_reasons (
    _id BIGINT NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    date TIMESTAMP NOT NULL,
    reason VARCHAR(280) NOT NULL,
    xp BIGINT NOT NULL
);
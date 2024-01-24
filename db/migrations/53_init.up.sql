ALTER TABLE byte_attempts ADD COLUMN completed_easy bool NOT NULL default false;
ALTER TABLE byte_attempts ADD COLUMN completed_medium bool NOT NULL default false;
ALTER TABLE byte_attempts ADD COLUMN completed_hard bool NOT NULL default false;
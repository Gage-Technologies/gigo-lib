ALTER TABLE bytes RENAME COLUMN description TO description_medium;
ALTER TABLE bytes ADD COLUMN description_easy varchar(500) NOT NULL;
ALTER TABLE bytes ADD COLUMN  description_hard varchar(500) NOT NULL;

ALTER TABLE bytes RENAME COLUMN outline_content TO outline_content_medium;
ALTER TABLE bytes ADD COLUMN outline_content_easy longtext NOT NULL;
ALTER TABLE bytes ADD COLUMN outline_content_hard longtext NOT NULL;

ALTER TABLE bytes RENAME COLUMN dev_steps TO dev_steps_medium;
ALTER TABLE bytes ADD COLUMN dev_steps_easy longtext;
ALTER TABLE bytes ADD COLUMN dev_steps_hard longtext;


ALTER TABLE byte_attempts RENAME COLUMN content TO content_medium;
ALTER TABLE byte_attempts ADD COLUMN content_easy longtext NOT NULL;
ALTER TABLE byte_attempts ADD COLUMN content_hard longtext NOT NULL;
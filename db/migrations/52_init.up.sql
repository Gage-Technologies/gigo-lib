ALTER TABLE bytes CHANGE  description_medium varchar(500) NOT NULL;
ALTER TABLE bytes ADD COLUMN description_easy varchar(500) NOT NULL;
ALTER TABLE bytes ADD COLUMN  description_hard varchar(500) NOT NULL;

ALTER TABLE bytes CHANGE outline_content outline_content_medium longtext NOT NULL;
ALTER TABLE bytes ADD COLUMN outline_content_easy longtext NOT NULL;
ALTER TABLE bytes ADD COLUMN outline_content_medium longtext NOT NULL;

ALTER TABLE bytes CHANGE dev_steps dev_steps_medium longtext;
ALTER TABLE bytes ADD COLUMN dev_steps_easy longtext;
ALTER TABLE bytes ADD COLUMN dev_steps_hard longtext;


ALTER TABLE byte_attempts CHANGE content content_easy longtext NOT NULL;
ALTER TABLE byte_attempts ADD COLUMN content_easy longtext NOT NULL;
ALTER TABLE byte_attempts ADD COLUMN content_hard longtext NOT NULL;
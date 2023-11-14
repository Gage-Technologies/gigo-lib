ALTER TABLE workspace_config ADD COLUMN uses int NOT NULL default 0;
ALTER TABLE workspace_config ADD COLUMN completions int NOT NULL default 0;
ALTER TABLE user_inactivity ADD COLUMN week_sent boolean NOT NULL;
ALTER TABLE `journey_units` ADD column `author_id` bigint NOT NULL;
ALTER TABLE `journey_units` ADD column `visibility` int NOT NULL;

ALTER TABLE workspace_config ADD COLUMN uses int NOT NULL default 0;
ALTER TABLE workspace_config ADD COLUMN completions int NOT NULL default 0;
ALTER TABLE user_inactivity ADD COLUMN week_sent boolean NOT NULL;
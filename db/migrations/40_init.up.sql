ALTER TABLE workspace_config ADD COLUMN uses int NOT NULL default 0;
ALTER TABLE workspace_config ADD COLUMN completions int NOT NULL default 0;
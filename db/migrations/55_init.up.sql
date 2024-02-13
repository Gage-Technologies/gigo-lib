-- Create journey_units table first because other tables reference it
CREATE TABLE IF NOT EXISTS journey_units(
            _id BIGINT NOT NULL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            description VARCHAR(500) NOT NULL,
            unit_above BIGINT,
            unit_below BIGINT,
            langs JSON NOT NULL,
            published BOOLEAN NOT NULL DEFAULT false,
            color VARCHAR(7) NOT NULL DEFAULT '#29C18C'
);

-- Create journey_tasks table with a foreign key reference to journey_units
CREATE TABLE IF NOT EXISTS journey_tasks(
            _id BIGINT NOT NULL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            description VARCHAR(500) NOT NULL,
            journey_unit_id BIGINT NOT NULL,
            node_above BIGINT,
            node_below BIGINT,
            code_source_id BIGINT,
            code_source_type INT NOT NULL,
            lang INT NOT NULL,
            published BOOLEAN NOT NULL DEFAULT false
);

-- Create journey_detour table with primary key and foreign key references
CREATE TABLE IF NOT EXISTS journey_detour(
             detour_unit_id BIGINT NOT NULL,
             user_id BIGINT NOT NULL,
             task_id BIGINT NOT NULL,
             started_at DATETIME NOT NULL,
             PRIMARY KEY (detour_unit_id, user_id)
);

-- Create journey_detour_recommendation table with a primary key
CREATE TABLE IF NOT EXISTS journey_detour_recommendation(
            _id BIGINT NOT NULL PRIMARY KEY,
            user_id BIGINT NOT NULL,
            recommended_unit BIGINT NOT NULL,
            created_at DATETIME NOT NULL,
            from_task_id BIGINT NOT NULL,
            accepted BOOLEAN NOT NULL DEFAULT false
);

-- Create journey_user_map table with a composite primary key
CREATE TABLE IF NOT EXISTS journey_user_map(
           user_id BIGINT NOT NULL,
           unit_id BIGINT NOT NULL,
           started_at DATETIME NOT NULL
);

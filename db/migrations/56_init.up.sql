-- Create journey_unit_tags table with a composite primary key
CREATE TABLE IF NOT EXISTS journey_unit_tags(
       unit_id BIGINT NOT NULL,
       value varchar(256) NOT NULL,
       primary key (unit_id, value)
);

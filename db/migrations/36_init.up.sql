create table if not exists journey_info(
    _id bigint not null primary key,
    user_id bigint not null,
    learning_goal varchar(280),
    selected_language varchar(280),
    end_goal varchar(280),
    experience_level varchar(280),
    familiarity_ide varchar(280),
    familiarity_linux varchar(280),
    tried varchar(280),
    tried_online varchar(280),
    aptitude_level varchar(280)
);
-- MySQL 5.7

CREATE TABLE IF NOT EXISTS web_tracking (
    `_id` bigint not null primary key,
    `user_id` bigint default NULL,
    `ip` bigint not null,
    `host` varchar(255) not null,
    `event` varchar(255) not null,
    `timestamp` datetime not null,
    `timespent` bigint default NULL,
    `path` varchar(255) not null,
    `lattitude` double not null,
    `longitude` double not null,
    `metadata` json
)
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
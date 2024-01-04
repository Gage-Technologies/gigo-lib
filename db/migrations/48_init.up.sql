-- MySQL v5.7

CREATE TABLE IF NOT EXISTS workspace_pool (
    `_id` bigint not null primary key,
    `container` varchar(500) not null,
    `state` int not null,
    `memory` bigint not null,
    `cpu` bigint not null,
    `volume_size` bigint not null,
    `secret`  binary(16) not null,
    `agent_id` bigint not null,
    `workspace_table_id` bigint not null
    );

create table if not exists bytes(
    _id bigint not null primary key,
    name varchar(255) not null,
    description varchar(500) not null,
    outline_content longtext not null,
    dev_steps longtext
);

create table if not exists byte_attempts(
    _id bigint not null primary key,
    byte_id bigint not null,
    author_id bigint not null,
    content longtext not null
);
-- MySQL v5.7

CREATE TABLE IF NOT EXISTS workspace_pool (
    `_id` bigint not null primary key,
    `container` varchar(500) not null,
    `state` int not null,
    `memory` bigint not null,
    `cpu` bigint not null,
    `storage` bigint not null,
    `secret`  binary(16) not null,
    `workspace_table_id` bigint not null
    )
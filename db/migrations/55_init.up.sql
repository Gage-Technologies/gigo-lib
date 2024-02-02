create table if not exists journey_tasks(
        `_id` bigint not null primary key,
        `name` varchar(255) not null,
        `description` varchar(500) not null,
        `journey_unit_id` bigint not null,
        `node_above` bigint,
        `node_below` bigint,
        `code_source_id`  bigint,
        `code_source_type` int not null,
        `published` boolean not null default false
);

create table if not exists journey_units(
        `_id` bigint not null primary key,
        `name` varchar(255) not null,
        `description` varchar(500) not null,
        `unit_above` bigint,
        `unit_below` bigint,
        `langs` json not null,
        `published` boolean not null default false,
        `color` VARCHAR(7) NOT NULL DEFAULT '#29C18C'
);

create table if not exists journey_detour(
         `detour_unit_id` bigint not null key,
         `user_id` bigint not null,
         `task_id` bigint not null,
         `started_at` datetime not null,
         primary key (detour_unit_id, user_id)
);

create table if not exists journey_detour_recommendation(
        `_id` bigint not null primary key,
        `user_id` bigint not null,
        `recommended_unit` bigint not null,
        `created_at` datetime not null,
        `from_task_id` bigint not null,
        `accepted` boolean not null default false
);
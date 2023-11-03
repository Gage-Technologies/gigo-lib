-- Mysql v5.7

-- Add start_time bigint column to workspace table

alter table `workspaces` add column `start_time` bigint;

-- Add start_time bigint column to attempt table

alter table `attempt` add column `start_time` bigint;

-- Add start_time bigint column to post table

alter table `post` add column `start_time` bigint;
create table if not exists journey_unit (
    _id bigint not null primary key,
    title varchar(255) not null,
    unit_focus varchar(255) not null,
    description varchar(500) not null,
    workspace_config bigint not null,
    workspace_config_revision bigint not null,
    workspace_settings json
);

create table if not exists journey_unit_langs (
    unit_id bigint not null,
    lang_id bigint not null,
    primary key (unit_id, lang_id)
);

create table if not exists journey_unit_projects (
    _id bigint not null primary key,
    unit_id bigint not null,
    repo_id bigint not null,
    title varchar(255) not null,
    description varchar(500) not null,
    project_language int not null,
    estimated_time_completion bigint
);

create table if not exists journey_unit_project_tags (
    project_id bigint not null,
    tag_id bigint not null,
    primary key (project_id, tag_id)
);

create table if not exists journey_unit_project_dependencies (
    project_id bigint not null,
    dependency_id bigint not null,
    primary key (project_id, dependency_id)
);
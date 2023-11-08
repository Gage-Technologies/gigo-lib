ALTER TABLE attempt ADD COLUMN title varchar(50);
create table if not exists journey_units (
    _id bigint not null primary key,
    title varchar(255) not null,
    unit_focus varchar(255) not null,
    description varchar(500) not null,
    repo_id bigint(20) not null,
    created_at datetime not null,
    updated_at datetime not null,
    challenge_cost varchar(16),
    completions bigint not null,
    attempts bigint not null,
    tier int not null,
    embedded boolean not null default false,
    workspace_config bigint not null,
    workspace_config_revision bigint not null,
    workspace_settings json,
    estimated_tutorial_time bigint
);

create table if not exists journey_unit_languages (
    unit_id bigint not null,
    value varchar(255) not null,
    is_attempt boolean not null default false,
    primary key (unit_id, lang_id)
);

create table if not exists journey_unit_projects (
    _id bigint not null primary key,
    unit_id bigint not null,
    completions bigint not null,
    working_directory varchar(255) not null,
    title varchar(50) not null,
    description varchar(500) not null,
    project_language int not null,
    estimated_time_completion bigint
);

create table if not exists journey_unit_project_tags (
    journey_id bigint not null,
    value varchar(255) not null,
    type int not null,
    primary key (journey_id, tag_id)
);

create table if not exists journey_unit_project_dependencies (
    project_id bigint not null,
    dependency_id bigint not null,
    primary key (project_id, dependency_id)
);

create table if not exists journey_unit_project_attempts (
    _id bigint not null primary key,
    unit_id bigint not null,
    parent_project bigint not null,
    is_completed boolean not null default false,
    working_directory varchar(255) not null,
    title varchar(50) not null,
    description varchar(500) not null,
    project_language int not null,
    estimated_tutorial_time bigint not null
)

create table if not exists journey_unit_attempts (
    _id bigint not null primary key,
    title varchar(50) not null,
    user_id bigint not null,
    parent_unit bigint not null,
    unit_focus varchar(255) not null,
    description varchar(500) not null,
    repo_id bigint(20) not null,
    created_at datetime not null,
    updated_at datetime not null,
    challenge_cost varchar(16) not null,
    completions bigint not null,
    attempts bigint not null,
    tier int not null,
    embedded boolean not null default false,
    workspace_config bigint not null,
    workspace_config_revision bigint not null,
    workspace_settings json,
    estimated_tutorial_time bigint
)
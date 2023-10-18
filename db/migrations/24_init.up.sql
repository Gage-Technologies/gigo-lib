create table if not exists curated_post (
    _id bigint not null primary key,
    post_id bigint not null,
    proficiency_type int not null,
    post_language int not null
);
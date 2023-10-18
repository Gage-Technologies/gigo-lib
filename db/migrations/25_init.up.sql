create table if not exists curated_post_type (
     curated_id bigint not null,
     proficiency_type int not null,
     primary key (curated_id, proficiency_type)
);
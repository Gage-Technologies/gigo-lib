create table if not exists report_issue (
    _id bigint not null primary key,
    user_id bigint not null,
    date datetime not null,
    issue longtext not null,
    page varchar(36) not null
)
create table if not exists exclusive_content_purchases (
    user_id bigint not null,
    post bigint not null,
    date datetime not null,
    primary key (user_id, post)
);
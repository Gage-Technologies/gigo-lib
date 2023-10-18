create table if not exists chat (
    _id bigint not null primary key,
    name varchar(280) not null,
    type int not null,
    last_message_time datetime
);

create table if not exists chat_users (
    chat_id bigint not null,
    user_id bigint not null,
    created_at datetime not null,
    primary key (chat_id, user_id)
);

create table if not exists chat_messages (
    _id bigint not null primary key,
    chat_id bigint not null,
    author_id bigint not null,
    author varchar(280) not null,
    message varchar(500) not null,
    -- We use datetime(6) to store the milliseconds - really important for chat messages
    created_at datetime(6) not null,
    revision bigint not null,
    type int not null
);

-- Add the global chat for all users
insert into chat (_id, name, type, last_message_time) values (0, 'Global', 0, null);

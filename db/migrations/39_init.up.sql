create table if not exists user_inactivity(
      user_id bigint not null primary key,
      last_login datetime not null,
      last_notified datetime not null,
      should_notify boolean not null
)
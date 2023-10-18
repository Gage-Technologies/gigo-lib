-- Add table for tracking volumes in the volume pool
create table if not exists volpool_volume(
    _id bigint not null primary key,
    size int not null,
    state int not null,
    pvc_name varchar(500) not null,
    workspace_id bigint
)
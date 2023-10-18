package gen

-- name: InsertDatabaseVersion :exec
insert
ignore into database_versions (
    version,
    date
) values (?, utc_timestamp());

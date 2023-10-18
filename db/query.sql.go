package ti

import "context"

const insertDatabaseVersion = `-- name: InsertDatabaseVersion :exec
insert
ignore into database_versions (
    version,
    date
) values (?, utc_timestamp())
`

func (q *Queries) InsertDatabaseVersion(ctx context.Context, version int64) error {
	_, err := q.db.ExecContext(ctx, insertDatabaseVersion, version)
	return err
}

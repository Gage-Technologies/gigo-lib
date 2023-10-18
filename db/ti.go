package ti

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"os"
	"path/filepath"
	"time"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Database struct {
	*sql.DB
	DBName string
	Host   string
	Port   string
	User   string
	Pass   string
	*Queries
}

func CreateDatabase(host string, port string, driverName string, username string, password string, databaseName string) (*Database, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true&tidb_skip_isolation_level_check=1", username, password, host, port)

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("create database if not exists " + databaseName)
	if err != nil {
		return nil, err
	}

	err = db.Close()
	if err != nil {
		return nil, err
	}

	dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&tidb_skip_isolation_level_check=1", username, password, host, port, databaseName)

	db, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(16384)
	db.SetMaxOpenConns(16384)
	db.SetConnMaxLifetime(time.Minute * 5)

	dataB := Database{
		DB:     db,
		DBName: databaseName,
		Host:   host,
		Port:   port,
		User:   username,
		Pass:   password,
	}

	// create a tmp directory and write all migrations to it
	tmp, err := os.MkdirTemp("", "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tmp)
	migrationFiles, err := migrations.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %s", err)
	}
	for _, file := range migrationFiles {
		buf, err := migrations.ReadFile(filepath.Join("migrations", file.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file: %s", err)
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", tmp, file.Name()), buf, 0600)
		if err != nil {
			return nil, fmt.Errorf("failed to write migration file: %s", err)
		}
	}

	migrationDriver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migration driver: %s", err)
	}
	migrator, err := migrate.NewWithDatabaseInstance("file://"+tmp, "mysql", migrationDriver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %s", err)
	}

	err = migrator.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			return &dataB, nil
		}
		return nil, fmt.Errorf("failed to run migrations: %s", err)
	}

	return &dataB, nil
}

func Close(db *Database) error {
	err := db.DB.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to close database: %d", err))
	}
	return nil
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func (db *Database) QueryContext(ctx context.Context, span *trace.Span, callerName *string, query string, args ...interface{}) (*sql.Rows, error) {
	if span != nil && callerName != nil {
		nSpan := *span
		dctx, dbspan := nSpan.TracerProvider().Tracer("gigo-core").Start(ctx, fmt.Sprintf("%v-db-query-call", *callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(query))
		defer dbspan.End()
		ctx = dctx
	}

	return db.DB.QueryContext(ctx, query, args...)
}

func (db *Database) QueryRowContext(ctx context.Context, span *trace.Span, callerName *string, query string, args ...interface{}) *sql.Row {
	if span != nil && callerName != nil {
		nSpan := *span
		dctx, dbspan := nSpan.TracerProvider().Tracer("gigo-core").Start(ctx, fmt.Sprintf("%v-db-queryrow-call", *callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(query))
		defer dbspan.End()
		ctx = dctx
	}

	return db.DB.QueryRowContext(ctx, query, args...)
}

func (db *Database) ExecContext(ctx context.Context, span *trace.Span, callerName *string, query string, args ...interface{}) (sql.Result, error) {
	if span != nil && callerName != nil {
		nSpan := *span

		dctx, dbspan := nSpan.TracerProvider().Tracer("gigo-core").Start(ctx, fmt.Sprintf("%v-db-exec-call", callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(query))
		defer dbspan.End()
		ctx = dctx
	}

	return db.DB.ExecContext(ctx, query, args...)
}

func (db *Database) PrepareContext(ctx context.Context, span *trace.Span, callerName *string, query string) (*sql.Stmt, error) {
	if span != nil && callerName != nil {
		nSpan := *span

		dctx, dbspan := nSpan.TracerProvider().Tracer("gigo-core").Start(ctx, fmt.Sprintf("%v-db-prepare-call", callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(query))
		defer dbspan.End()
		ctx = dctx
	}
	return db.DB.PrepareContext(ctx, query)
}

func (db *Database) Query(ctx context.Context, span *trace.Span, callerName *string, query string, args ...interface{}) (*sql.Rows, error) {
	if span != nil && callerName != nil {
		nSpan := *span

		dctx, dbspan := nSpan.TracerProvider().Tracer("gigo-core").Start(ctx, fmt.Sprintf("%v-db-query-call", *callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(query))
		defer dbspan.End()
		ctx = dctx
	}

	return db.DB.Query(query, args...)

}

func (db *Database) QueryRow(ctx context.Context, span *trace.Span, callerName *string, query string, args ...interface{}) *sql.Row {
	if span != nil && callerName != nil {
		nSpan := *span

		dctx, dbspan := nSpan.TracerProvider().Tracer("gigo-core").Start(ctx, fmt.Sprintf("%v-db-queryrow-call", *callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(query))
		defer dbspan.End()
		ctx = dctx
	}

	return db.DB.QueryRow(query, args...)

}

func (db *Database) Exec(ctx context.Context, span *trace.Span, callerName *string, query string, args ...interface{}) (sql.Result, error) {
	if span != nil && callerName != nil {
		nSpan := *span

		dctx, dbspan := nSpan.TracerProvider().Tracer("gigo-core").Start(ctx, fmt.Sprintf("%v-db-exec-call", *callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(query))
		defer dbspan.End()
		ctx = dctx
	}

	return db.DB.Exec(query, args...)

}

func (db *Database) Prepare(ctx context.Context, span *trace.Span, callerName *string, query string) (*sql.Stmt, error) {
	if span != nil && callerName != nil {
		nSpan := *span

		dctx, dbspan := nSpan.TracerProvider().Tracer("gigo-core").Start(ctx, fmt.Sprintf("%v-db-prepare-call", *callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(query))
		defer dbspan.End()
		ctx = dctx
	}

	return db.DB.Prepare(query)

}

type Tx struct {
	*sql.Tx
	span       *trace.Span
	ctx        context.Context
	callerName *string
}

func (db *Database) BeginTx(ctx context.Context, span *trace.Span, callerName *string, opts *sql.TxOptions) (*Tx, error) {
	if span != nil && callerName != nil {
		nSpan := *span

		dctx, dbspan := nSpan.TracerProvider().Tracer("gigo-core").Start(ctx, fmt.Sprintf("%v-db-tx-call", *callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(""))
		ctx = dctx
	}
	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Tx{tx, span, ctx, callerName}, nil
}

func (tx *Tx) Query(callerName *string, query string, args ...any) (*sql.Rows, error) {
	if tx.span != nil {
		span := *tx.span
		_, dbspan := span.TracerProvider().Tracer("gigo-core").Start(tx.ctx, fmt.Sprintf("%v-db-query-tx-call", *callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(query))
		defer dbspan.End()
	}

	return tx.Tx.Query(query, args...)
}

func (tx *Tx) QueryRow(callerName *string, query string, args ...any) *sql.Row {
	if tx.span != nil {
		span := *tx.span
		_, dbspan := span.TracerProvider().Tracer("gigo-core").Start(tx.ctx, fmt.Sprintf("%v-db-queryrow-tx-call", *callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(query))
		defer dbspan.End()
	}

	return tx.Tx.QueryRow(query, args...)
}

func (tx *Tx) QueryContext(ctx context.Context, callerName *string, query string, args ...any) (*sql.Rows, error) {
	if tx.span != nil {
		span := *tx.span
		_, dbspan := span.TracerProvider().Tracer("gigo-core").Start(tx.ctx, fmt.Sprintf("%v-db-query-tx-call", *callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(query))

		defer dbspan.End()
		ctx = tx.ctx
	}

	return tx.Tx.QueryContext(ctx, query, args...)
}

func (tx *Tx) Exec(callerName *string, query string, args ...any) (sql.Result, error) {
	if tx.span != nil {
		span := *tx.span
		_, dbspan := span.TracerProvider().Tracer("gigo-core").Start(tx.ctx, fmt.Sprintf("%v-db-exec-tx-call", *callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(query))
		defer dbspan.End()
	}

	return tx.Tx.Exec(query, args...)

}

func (tx *Tx) ExecContext(ctx context.Context, callerName *string, query string, args ...any) (sql.Result, error) {
	if tx.span != nil {
		span := *tx.span
		_, dbspan := span.TracerProvider().Tracer("gigo-core").Start(tx.ctx, fmt.Sprintf("%v-db-exec-tx-call", *callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(query))
		defer dbspan.End()
		ctx = tx.ctx
	}

	return tx.Tx.ExecContext(ctx, query, args...)

}

func (tx *Tx) Prepare(callerName *string, query string) (*sql.Stmt, error) {
	if tx.span != nil {
		span := *tx.span
		_, dbspan := span.TracerProvider().Tracer("gigo-core").Start(tx.ctx, fmt.Sprintf("%v-db-prepare-tx-call", *callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(query))
		defer dbspan.End()
	}

	return tx.Tx.Prepare(query)

}

func (tx *Tx) PrepareContext(ctx context.Context, callerName *string, query string) (*sql.Stmt, error) {
	if tx.span != nil {
		span := *tx.span
		_, dbspan := span.TracerProvider().Tracer("gigo-core").Start(tx.ctx, fmt.Sprintf("%v-db-prepare-tx-call", *callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(query))
		defer dbspan.End()
		ctx = tx.ctx
	}

	return tx.Tx.PrepareContext(ctx, query)

}

func (tx *Tx) Commit(callerName *string) error {
	if tx.span != nil {
		span := *tx.span
		_, dbspan := span.TracerProvider().Tracer("gigo-core").Start(tx.ctx, fmt.Sprintf("%v-db-tx-call", *callerName), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(""))
		defer func() {
			dbspan.End()
			(*tx.span).End()
		}()
	}

	return tx.Tx.Commit()
}

func (tx *Tx) Rollback() error {
	if tx.span != nil {
		span := *tx.span
		name := "rollback"
		if tx.callerName != nil {
			name = fmt.Sprintf("%v-%v", *tx.callerName, name)
		}
		_, dbspan := span.TracerProvider().Tracer("gigo-core").Start(tx.ctx, fmt.Sprintf("%v-db-tx-call", name), trace.WithSpanKind(trace.SpanKindInternal))
		dbspan.SetAttributes(semconv.DBSystemKey.String("tidb"),
			semconv.DBNameKey.String("gigo"),
			semconv.DBStatementKey.String(""))
		defer func() {
			dbspan.End()
			(*tx.span).End()
		}()
	}

	return tx.Tx.Rollback()
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}

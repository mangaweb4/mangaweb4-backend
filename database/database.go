package database

import (
	"context"
	"database/sql"

	"entgo.io/ent/dialect"
	dialect_sql "entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/wutipong/mangaweb3-backend/configuration"
	"github.com/wutipong/mangaweb3-backend/ent"
	"github.com/wutipong/mangaweb3-backend/ent/migrate"
	"github.com/wutipong/mangaweb3-backend/errors"
	_ "modernc.org/sqlite"
)

var pool *pgxpool.Pool

var databaseType string
var connectionStr string

func Open(ctx context.Context, dbType string, connStr string) error {
	databaseType = dbType
	connectionStr = connStr

	switch dbType {
	case dialect.Postgres:
		if p, e := pgxpool.New(ctx, connStr); e == nil {
			pool = p
			return nil
		} else {
			return e
		}

	case dialect.SQLite:
		return nil
	}

	return errors.ErrNotImplemented

}

func openDB(dbType string) (db *dialect_sql.Driver, err error) {
	switch dbType {
	case dialect.Postgres:
		return openPostgres()

	case dialect.SQLite:
		return openSqlite()
	}

	return nil, errors.ErrNotImplemented
}

func openPostgres() (db *dialect_sql.Driver, err error) {
	db = dialect_sql.OpenDB(dialect.Postgres, stdlib.OpenDBFromPool(pool))

	return
}

func openSqlite() (db *dialect_sql.Driver, err error) {
	if d, e := sql.Open("sqlite", connectionStr); e == nil {
		db = dialect_sql.OpenDB(dialect.SQLite, d)
		return
	} else {
		e = err
		return
	}
}

func Close() {
	if pool != nil {
		pool.Close()
	}
}

func CreateEntClient() *ent.Client {
	db, err := openDB(databaseType)
	if err != nil {
		return nil
	}

	options := []ent.Option{
		ent.Driver(db),
	}

	config := configuration.Get()
	if config.DebugMode {
		options = append(options,
			ent.Debug(),
			ent.Log(func(params ...any) {
				if databaseType == "postgres" {
					stat := pool.Stat()

					log.Debug().
						Any("params", params).
						Int32("Acquired Conns", stat.AcquiredConns()).
						Int32("Idle Conns", stat.IdleConns()).
						Int32("Constructed Conns", stat.ConstructingConns()).
						Msg("Ent Debug")
				}
			}),
		)
	}

	client := ent.NewClient(options...)

	return client
}

func CreateSchema(ctx context.Context) error {
	client := CreateEntClient()
	defer client.Close()

	return client.Schema.Create(ctx, migrate.WithDropColumn(true), migrate.WithDropIndex(true))
}

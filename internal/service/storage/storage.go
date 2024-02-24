package storage

import (
	"context"
	"fmt"
	"l0/internal/config"
	"log"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgxIface interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
	Ping(context.Context) error
}

type Storage struct {
	Db PgxIface
}

func New(cfg *config.DB) *Storage {
	connectionStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DatabaseName,
	)
	db, err := pgx.Connect(context.Background(), connectionStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	if err := db.Ping(context.TODO()); err != nil {
		panic(err)
	}
	if err = migrate(db); err != nil {
		return nil
	}
	return &Storage{
		Db: db,
	}
}

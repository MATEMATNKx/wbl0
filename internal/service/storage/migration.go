package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

var migrations = []func(tx *pgx.Tx) error{
	initial,
}
var maxVersion = len(migrations)

func migrate(db *pgx.Conn) error {
	for v := 1; v <= maxVersion; v++ {
		err := migrateVesion(v, db)
		if err != nil {
			return err
		}
	}
	return nil
}
func migrateVesion(v int, db *pgx.Conn) error {
	var err error
	var tx pgx.Tx
	migrationFunc := migrations[v-1]

	if tx, err = db.Begin(context.TODO()); err != nil {
		log.Printf("migration[%d] failed to start transaction: %s\n", v, err.Error())
		return err
	}

	if err = migrationFunc(&tx); err != nil {
		log.Printf("migration[%d] failed to migrate: %s\n", v, err.Error())
		tx.Rollback(context.TODO())
		return err
	}
	if err = tx.Commit(context.TODO()); err != nil {
		log.Printf("migration[%d] failed to commit changes: %s\n", v, err.Error())
		return err
	}
	return nil
}

func initial(tx *pgx.Tx) error {
	query := `
	create table if not exists orders (
		order_uid text primary key,
		data text not null
	);
	`
	_, err := (*tx).Exec(context.TODO(), query)
	return err
}

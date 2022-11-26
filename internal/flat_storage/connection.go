package flat_storage

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/kontsevoye/rentaflat/internal/config"
	_ "github.com/lib/pq"
)

func CreateDbConnection(cfg *config.AppConfig) (*sql.DB, error) {
	db, err := sql.Open(
		"postgres",
		cfg.DatabaseDSN,
	)

	return db, err
}

func CreateDbxConnection(db *sql.DB) *sqlx.DB {
	return sqlx.NewDb(db, "postgres")
}

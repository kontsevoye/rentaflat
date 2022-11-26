package parser

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func CreateDbConnection(cfg *AppConfig) (*sql.DB, error) {
	db, err := sql.Open(
		"postgres",
		cfg.DatabaseDSN,
	)

	return db, err
}

func CreateDbxConnection(db *sql.DB) *sqlx.DB {
	return sqlx.NewDb(db, "postgres")
}

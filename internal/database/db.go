package database

import (
	"github.com/felipecaputo/go-voting-api/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func migrateUp(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user (
			id varchar(40) NOT NULL,
			email varchar(250) NOT NULL,
			password varchar(250) NOT NULL,
			name varchar(200) NOT NULL,
			is_admin BOOL DEFAULT false NOT NULL,
			CONSTRAINT user_PK PRIMARY KEY (id),
			CONSTRAINT email_UNQ UNIQUE KEY (email)
		)
		ENGINE=InnoDB
		DEFAULT CHARSET=utf8mb4
		COLLATE=utf8mb4_0900_ai_ci;
		`)

	return err
}

func NewDB(config *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", config.ConnectionString)

	if err != nil {
		return nil, err
	}

	if config.Environment == "local" {
		if err = migrateUp(db); err != nil {
			return nil, err
		}
	}

	return db, err
}

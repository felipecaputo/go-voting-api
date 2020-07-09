package database

import (
	"github.com/felipecaputo/go-voting-api/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func migrateUp(db *sqlx.DB) {
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

	db.Exec("CREATE INDEX user_email_IDX USING BTREE ON user (email);")

	if err != nil {
		panic(err)
	}
}

func NewDB(config *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", config.ConnectionString)

	if err == nil && config.Environment == "local" {
		migrateUp(db)
	}

	return db, err
}

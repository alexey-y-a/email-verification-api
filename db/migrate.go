package db

import "database/sql"

func Migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXIST verifications (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL,
		hash TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP,
		expires_at TIMESTAMP NOT NULL,
		verified BOOLEAN NOT NULL DEFAULT false
	);`
	_, err := db.Exec(query)
	return err
}

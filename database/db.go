package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Create db tables and connection
func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "memory.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages(id INTEGER PRIMARY KEY, role TEXT, text TEXT)`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

var DB = InitDB()

package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

var DB *sql.DB

func InitDB() error {
	// Create data directory if not exists
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		if err := os.Mkdir("data", 0755); err != nil {
			return fmt.Errorf("failed to create data directory: %v", err)
		}
	}

	db, err := sql.Open("sqlite3", "data/task.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Create tasks table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		file TEXT,
		description TEXT,
		status TEXT,
		dependencies TEXT
	)
	`)
	if err != nil {
		return fmt.Errorf("failed to create tasks table: %v", err)
	}

	// Create file_documents table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS file_documents (
		file_path TEXT PRIMARY KEY,
		content TEXT
	)
	`)
	if err != nil {
		return fmt.Errorf("failed to create file_documents table: %v", err)
	}

	DB = db
	return nil
}
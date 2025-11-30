package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

// DB is the global database connection
var DB *sql.DB

// Init initializes the database connection and creates the schema
func Init() error {
	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	// Open SQLite database
	var err error
	DB, err = sql.Open("sqlite", "data/todolist.db")
	if err != nil {
		return err
	}

	// Test connection
	if err = DB.Ping(); err != nil {
		return err
	}

	// Create schema
	if err = createSchema(); err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// createSchema creates the database tables if they don't exist
func createSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		completed INTEGER DEFAULT 0,
		notes TEXT,
		carried_from_date TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_date ON tasks(date);
	`

	_, err := DB.Exec(schema)
	return err
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

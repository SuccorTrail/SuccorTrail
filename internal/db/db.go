package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "succortrail.db")
	if err != nil {
		return err
	}

	// Create tables
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS donations (
        id TEXT PRIMARY KEY,
        name TEXT,
        email TEXT,
        phone TEXT,
        meal_type TEXT,
        quantity INTEGER,
        expiry_date DATETIME,
        location TEXT,
        notes TEXT,
        created_at DATETIME,
        updated_at DATETIME
    );
    CREATE TABLE IF NOT EXISTS receivers (
        id TEXT PRIMARY KEY,
        name TEXT,
        phone TEXT,
        location TEXT,
        family_size INTEGER,
        dietary_restrictions TEXT,
        created_at DATETIME
    );
    CREATE TABLE IF NOT EXISTS feedback (
        id TEXT PRIMARY KEY,
        donation_id TEXT,
        quality INTEGER,
        comments TEXT,
        created_at DATETIME,
        FOREIGN KEY(donation_id) REFERENCES donations(id)
    );
    CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        name TEXT,
        email TEXT UNIQUE,
        password TEXT,
        user_type TEXT,
        created_at DATETIME
    );`
	_, err = db.Exec(sqlStmt)
	return err
}

func CloseDB() {
	db.Close()
}

func GetDB() *sql.DB {
	return db
}

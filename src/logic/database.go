package logic

import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitData() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println(err)
	}

	reset := flag.Bool("reset", false, "reset the database")
	flag.Parse()

	if *reset {
		deleteData()
		fmt.Println("Database reset")
	}

	createData()
	fmt.Println("Database initialized")
}

// NORMAL, GITHUB OR GOOGLE ACCOUNT
func createData() {
	query := `
	CREATE TABLE Users (
		user_id INTEGER PRIMARY KEY,
		username VARCHAR UNIQUE,
		email VARCHAR UNIQUE,
		password VARCHAR,
		creation_date DATETIME,
		rank_id INTEGER,
		profile_picture VARCHAR,
		account_type VARCHAR,
		FOREIGN KEY (rank_id) REFERENCES Ranks(rank_id)
	);`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
}

func deleteData() {
	query := `DROP TABLE users;`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
}

func checkUser(username, email string) bool {
	query := `SELECT * FROM users WHERE username = ? OR email = ?;`
	rows, err := db.Query(query, username, email)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	return rows.Next()
}

func checkPassword(username, password string) bool {
	query := `SELECT * FROM users WHERE username = ? AND password = ?;`
	rows, err := db.Query(query, username, password)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	return rows.Next()
}

func insertUser(username, email, password, profile_picture, account string) error {
	query := `INSERT INTO users (username, email, password, creation_date, rank_id, profile_picture, account_type) VALUES (?, ?, ?, datetime('now'), 1, ?, ?);`
	_, err := db.Exec(query, username, email, password, profile_picture, account)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

package logic

import (
	"database/sql"
	"flag"
	"fmt"

	"github.com/gofrs/uuid"
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
		uuid TEXT UNIQUE,
		username VARCHAR,
		email VARCHAR,
		password VARCHAR,
		creation_date DATETIME,
		rank_id INTEGER,
		profile_picture VARCHAR,
		account_type VARCHAR,
		FOREIGN KEY (rank_id) REFERENCES Ranks(rank_id)
	);
	
	CREATE TABLE Ranks (
		rank_id INTEGER PRIMARY KEY,
		rank_name VARCHAR
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

	query = `DROP TABLE ranks;`
	_, err = db.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
}

func checkUser(username, email string) bool {
	query := `SELECT username, email FROM users WHERE username = ? OR email = ?;`
	row := db.QueryRow(query, username, email)
	var dbUsername, dbEmail string
	err := row.Scan(&dbUsername, &dbEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		fmt.Println(err)
		return true
	}
	return true
}

func getPassword(username string) string {
	query := `SELECT password FROM users WHERE username = ?;`
	row := db.QueryRow(query, username)
	var password string
	err := row.Scan(&password)
	if err != nil {
		fmt.Println(err)
	}
	return password
}

func newUser(username, email, password, profile_picture, account string) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		return err
	}
	query := `INSERT INTO users (uuid, username, email, password, creation_date, rank_id, profile_picture, account_type) VALUES (?, ?, ?, ?, datetime('now'), 1, ?, ?);`
	_, err = db.Exec(query, uuid.String(), username, email, password, profile_picture, account)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

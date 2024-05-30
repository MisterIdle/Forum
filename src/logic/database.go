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
		uuid TEXT UNIQUE,
		username VARCHAR,
		email VARCHAR,
		password VARCHAR,
		identity TEXT,
		code TEXT,
		creation DATETIME,
		rank_id INTEGER,
		picture VARCHAR,
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

func checkUserMailOrIdentidy(identifier string) bool {
	query := `SELECT COALESCE(email, identity) FROM users WHERE email = ? OR identity = ?;`
	row := db.QueryRow(query, identifier, identifier)
	var result string
	err := row.Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

func setCodeByEmail(email, code string) error {
	query := `UPDATE users SET code = ? WHERE email = ?;`
	_, err := db.Exec(query, code, email)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func getEmailFromCode(code string) string {
	query := `SELECT email FROM users WHERE code = ?;`
	row := db.QueryRow(query, code)
	var email string
	err := row.Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		}
		fmt.Println(err)
		return ""
	}
	return email
}

func resetPassword(email, password string) error {
	query := `UPDATE users SET password = ? WHERE email = ?;`
	_, err := db.Exec(query, password, email)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func resetCode(email string) error {
	query := `UPDATE users SET code = null WHERE email = ?;`
	_, err := db.Exec(query, email)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func getCredentialsByUsernameOrEmail(identifier string) (string, string) {
	query := `SELECT password, COALESCE(username, email) FROM users WHERE username = ? OR email = ?;`
	row := db.QueryRow(query, identifier, identifier)
	var password, username string
	err := row.Scan(&password, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ""
		}
		fmt.Println(err)
		return "", ""
	}
	return password, username
}

func newUser(username, email, password, identity, picture string) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Set code to null if empty

	query := `INSERT INTO users (uuid, username, email, password, identity, code, creation, rank_id, picture) VALUES (?, ?, ?, ?, ?, null, datetime('now'), 1, ?);`
	_, err = db.Exec(query, uuid.String(), username, email, password, identity, picture)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

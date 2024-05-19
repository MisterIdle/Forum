package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable()

	fmt.Println("Starting the application...")
	Handle()
	launchApp()
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		password TEXT NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func launchApp() {
	fmt.Println("Server is running on port 3030")
	http.ListenAndServe(":3030", nil)
}

func Handle() {
	http.Handle(("/styles/"), http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))
	http.Handle(("/javascript/"), http.StripPrefix("/javascript/", http.FileServer(http.Dir("javascript"))))

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplateGlobal(w, "templates/index.html", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		data := Data{
			Username: username,
			Password: password,
		}

		if checkIfUserExists(username) {
			fmt.Println("User already exists")
			http.Redirect(w, r, "templates/index.html", http.StatusSeeOther)
		} else {
			insertData(data)
			RenderTemplateGlobal(w, "templates/register.html", "User registered successfully")
			return
		}
	}

	RenderTemplateGlobal(w, "templates/index.html", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		query := `SELECT username, password FROM users WHERE username = ? AND password = ?`
		row := db.QueryRow(query, username, password)
		var u, p string
		err := row.Scan(&u, &p)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("User not found")
				http.Redirect(w, r, "templates/index.html", http.StatusSeeOther)
			} else {
				log.Fatal(err)
			}
		}
		RenderTemplateGlobal(w, "templates/register.html", "User logged in successfully")
		return
	}

	RenderTemplateGlobal(w, "templates/index.html", nil)
}

func RenderTemplateGlobal(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func checkIfUserExists(username string) bool {
	query := `SELECT username FROM users WHERE username = ?`
	row := db.QueryRow(query, username)
	var u string
	err := row.Scan(&u)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	return true
}

func insertData(data Data) {
	query := `INSERT INTO users (username, password) VALUES (?, ?)`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(data.Username, data.Password)
	if err != nil {
		log.Fatal(err)
	}
}

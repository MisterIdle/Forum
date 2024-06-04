package handlers

import (
    "database/sql"
    "net/http"
    "html/template"
    _ "github.com/mattn/go-sqlite3"
)

func CreateCategory(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        tmpl, _ := template.ParseFiles("templates/create_category.html")
        tmpl.Execute(w, nil)
    } else if r.Method == http.MethodPost {
        name := r.FormValue("name")
        description := r.FormValue("description")

        db, _ := sql.Open("sqlite3", "./Forum3.db")
        defer db.Close()

        _, err := db.Exec("INSERT INTO Categories (name, description) VALUES (?, ?)", name, description)
        if err != nil {
            http.Error(w, "Error creating category", http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/categories", http.StatusSeeOther)
    }
}

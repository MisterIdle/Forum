package logic

import (
	"database/sql"
	"html/template"
	"net/http"
)

func ViewCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("id")

	db, _ := sql.Open("sqlite3", "./Forum3.db")
	defer db.Close()

	var category Category
	err := db.QueryRow(`
        SELECT category_id, name, description
        FROM Categories
        WHERE category_id = ?`, categoryID).Scan(&category.CategoryID, &category.Name)
	if err != nil {
		http.Error(w, "Error retrieving category", http.StatusInternalServerError)
		return
	}

	tmpl, _ := template.ParseFiles("templates/view_category.html")
	tmpl.Execute(w, category)
}

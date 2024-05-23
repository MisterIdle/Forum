package handlers

import (
    "database/sql"
    "net/http"
    "html/template"
    _ "github.com/mattn/go-sqlite3"
    "strconv"
)

type PostFormData struct {
    CategoryID int
    UserID     int
    Title      string
    Content    string
    Error      string
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        tmpl, _ := template.ParseFiles("templates/create_post.html")
        tmpl.Execute(w, nil)
    } else if r.Method == http.MethodPost {
        categoryID := r.FormValue("category_id")
        userID := r.FormValue("user_id") // Assume you get user_id from session in real implementation
        title := r.FormValue("title")
        content := r.FormValue("content")

        db, _ := sql.Open("sqlite3", "./Forum3.db")
        defer db.Close()

        var existingPostID int
        err := db.QueryRow("SELECT post_id FROM Posts WHERE title = ?", title).Scan(&existingPostID)
        if err == nil {
            tmpl, _ := template.ParseFiles("templates/create_post.html")
            tmpl.Execute(w, PostFormData{
                CategoryID: atoi(categoryID),
                UserID:     atoi(userID),
                Title:      title,
                Content:    content,
                Error:      "A post with this title already exists.",
            })
            return
        } else if err != sql.ErrNoRows {
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }

        _, err = db.Exec("INSERT INTO Posts (category_id, user_id, title, content, timestamp) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)", categoryID, userID, title, content)
        if err != nil {
            http.Error(w, "Error creating post", http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/category?id="+categoryID, http.StatusSeeOther)
    }
}

func atoi(str string) int {
    value, _ := strconv.Atoi(str)
    return value
}

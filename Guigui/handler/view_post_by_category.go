package handlers

import (
    "database/sql"
    "net/http"
    "html/template"
    _ "github.com/mattn/go-sqlite3"
)

func ViewCategoryPosts(w http.ResponseWriter, r *http.Request) {
    categoryID := r.URL.Query().Get("id")

    db, _ := sql.Open("sqlite3", "./Forum3.db")
    defer db.Close()

    rows, err := db.Query(`
        SELECT p.post_id, p.content, p.timestamp, u.username 
        FROM Posts p
        JOIN Users u ON p.user_id = u.user_id
        WHERE p.category_id = ?
        ORDER BY p.timestamp DESC`, categoryID)
    if err != nil {
        http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []struct {
        PostID    int
        Content   string
        Timestamp string
        Username  string
    }
    for rows.Next() {
        var post struct {
            PostID    int
            Content   string
            Timestamp string
            Username  string
        }
        if err := rows.Scan(&post.PostID, &post.Content, &post.Timestamp, &post.Username); err != nil {
            http.Error(w, "Error scanning posts", http.StatusInternalServerError)
            return
        }
        posts = append(posts, post)
    }

    tmpl, _ := template.ParseFiles("templates/view_category_posts.html")
    tmpl.Execute(w, posts)
}

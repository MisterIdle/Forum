package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    _ "github.com/mattn/go-sqlite3"
)

func Search(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("query")

    db, _ := sql.Open("sqlite3", "./Forum3.db")
    defer db.Close()

    var results SearchResult

    // Search by posts
    postRows, err := db.Query(`
        SELECT p.post_id, p.title, c.name AS category_name, COALESCE(COUNT(DISTINCT cm.user_id), 0) AS unique_users
        FROM Posts p
        JOIN Categories c ON p.category_id = c.category_id
        LEFT JOIN Comments cm ON p.post_id = cm.post_id
        WHERE p.title LIKE ? OR p.content LIKE ?
        GROUP BY p.post_id, c.name
        LIMIT 10`, "%"+query+"%", "%"+query+"%")
    if err != nil {
        http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
        return
    }
    defer postRows.Close()

    for postRows.Next() {
        var post PostResult
        if err := postRows.Scan(&post.PostID, &post.Title, &post.CategoryName, &post.UniqueUsers); err != nil {
            http.Error(w, "Error scanning posts", http.StatusInternalServerError)
            return
        }
        results.Posts = append(results.Posts, post)
    }

    // Search by users
    userRows, err := db.Query(`SELECT user_id, username FROM Users WHERE username LIKE ? LIMIT 10`, "%"+query+"%")
    if err != nil {
        http.Error(w, "Error retrieving users", http.StatusInternalServerError)
        return
    }
    defer userRows.Close()

    for userRows.Next() {
        var user UserResult
        if err := userRows.Scan(&user.UserID, &user.Username); err != nil {
            http.Error(w, "Error scanning users", http.StatusInternalServerError)
            return
        }
        results.Users = append(results.Users, user)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}

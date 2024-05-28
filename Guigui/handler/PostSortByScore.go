package handlers

import (
    "database/sql"
    "net/http"
    "html/template"
    "github.com/mattn/go-sqlite3"
)

func ViewPost(w http.ResponseWriter, r *http.Request) {
    postID := r.URL.Query().Get("id")

    db, _ := sql.Open("sqlite3", "./Forum3.db")
    defer db.Close()

    var post struct {
        PostID    int
        Content   string
        Timestamp string
        Username  string
    }
    err := db.QueryRow(`
        SELECT p.post_id, p.content, p.timestamp, u.username 
        FROM Posts p
        JOIN Users u ON p.user_id = u.user_id
        WHERE p.post_id = ?`, postID).Scan(&post.PostID, &post.Content, &post.Timestamp, &post.Username)
    if err != nil {
        http.Error(w, "Error retrieving post", http.StatusInternalServerError)
        return
    }

    rows, err := db.Query(`
        SELECT c.comment_id, c.content, c.timestamp, u.username, c.score
        FROM Comments c
        JOIN Users u ON c.user_id = u.user_id
        WHERE c.post_id = ?
        ORDER BY c.score DESC, c.timestamp DESC`, postID)
    if err != nil {
        http.Error(w, "Error retrieving comments", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var comments []struct {
        CommentID int
        Content   string
        Timestamp string
        Username  string
        Score     int
    }
    for rows.Next() {
        var comment struct {
            CommentID int
            Content   string
            Timestamp string
            Username  string
            Score     int
        }
        if err := rows.Scan(&comment.CommentID, &comment.Content, &comment.Timestamp, &comment.Username, &comment.Score); err != nil {
            http.Error(w, "Error scanning comments", http.StatusInternalServerError)
            return
        }
        comments = append(comments, comment)
    }

    tmpl, _ := template.ParseFiles("templates/view_post.html")
    tmpl.Execute(w, struct {
        Post     interface{}
        Comments interface{}
    }{
        Post:     post,
        Comments: comments,
    })
}

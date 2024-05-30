package handlers

import (
    "database/sql"
    "net/http"
    "html/template"
    "strings"
    _ "github.com/mattn/go-sqlite3"
)

func ViewPost(w http.ResponseWriter, r *http.Request) {
    postID := r.URL.Query().Get("id")

    db, _ := sql.Open("sqlite3", "./Forum3.db")
    defer db.Close()

    var post Post
    var imagePaths string
    err := db.QueryRow(`
        SELECT p.post_id, p.title, p.content, p.timestamp, u.username, p.image_paths 
        FROM Posts p
        JOIN Users u ON p.user_id = u.user_id
        WHERE p.post_id = ?`, postID).Scan(&post.PostID, &post.Title, &post.Content, &post.Timestamp, &post.Username, &imagePaths)
    if err != nil {
        http.Error(w, "Error retrieving post", http.StatusInternalServerError)
        return
    }
    post.ImagePaths = strings.Split(imagePaths, ",")

    rows, err := db.Query(`
        SELECT c.comment_id, c.content, c.timestamp, u.username, c.score, c.image_paths
        FROM Comments c
        JOIN Users u ON c.user_id = u.user_id
        WHERE c.post_id = ?
        ORDER BY c.score DESC, c.timestamp DESC`, postID)
    if err != nil {
        http.Error(w, "Error retrieving comments", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var comments []Comment
    for rows.Next() {
        var comment Comment
        var imagePaths string
        if err := rows.Scan(&comment.CommentID, &comment.Content, &comment.Timestamp, &comment.Username, &comment.Score, &imagePaths); err != nil {
            http.Error(w, "Error scanning comments", http.StatusInternalServerError)
            return
        }
        comment.ImagePaths = strings.Split(imagePaths, ",")
        comments = append(comments, comment)
    }

    tmpl, _ := template.ParseFiles("templates/view_post.html")
    tmpl.Execute(w, struct {
        Post     Post
        Comments []Comment
    }{
        Post:     post,
        Comments: comments,
    })
}

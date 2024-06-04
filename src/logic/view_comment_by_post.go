package logic

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"
)

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("id")

	db, _ := sql.Open("sqlite3", "./Forum3.db")
	defer db.Close()

	var post Post
	row := db.QueryRow(`
        SELECT p.post_id, p.title, p.content, p.image_paths, p.timestamp, p.user_id, u.username
        FROM Posts p
        JOIN Users u ON p.user_id = u.user_id
        WHERE p.post_id = ?`, postID)
	var imagePaths string
	err := row.Scan(&post.PostID, &post.Title, &post.Content, &imagePaths, &post.Timestamp, &post.UserID, &post.Username)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	post.ImagePaths = strings.Split(imagePaths, ",")

	rows, err := db.Query(`
        SELECT c.comment_id, c.content, c.image_paths, c.timestamp, c.user_id, u.username, c.score
        FROM Comments c
        JOIN Users u ON c.user_id = u.user_id
        WHERE c.post_id = ?
        ORDER BY c.timestamp ASC`, postID)
	if err != nil {
		http.Error(w, "Error retrieving comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		var commentImagePaths string
		if err := rows.Scan(&comment.CommentID, &comment.Content, &commentImagePaths, &comment.Timestamp, &comment.UserID, &comment.Username, &comment.Score); err != nil {
			http.Error(w, "Error scanning comments", http.StatusInternalServerError)
			return
		}
		comment.ImagePaths = strings.Split(commentImagePaths, ",")
		comments = append(comments, comment)
	}

	tmpl, _ := template.ParseFiles("templates/view_post.html")
	tmpl.Execute(w, PostData{Post: post, Comments: comments})
}

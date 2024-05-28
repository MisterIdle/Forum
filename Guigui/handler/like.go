package handlers

import (
    "database/sql"
    "net/http"
    _ "github.com/mattn/go-sqlite3"
)

func LikeComment(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        commentID := r.FormValue("comment_id")
        userID := r.FormValue("user_id") // Assume you get user_id from session in real implementation
        likeType := r.FormValue("type") // "like" or "dislike"

        db, _ := sql.Open("sqlite3", "./Forum3.db")
        defer db.Close()

        var scoreChange int
        if likeType == "like" {
            scoreChange = 1
        } else if likeType == "dislike" {
            scoreChange = -1
        }

        tx, err := db.Begin()
        if err != nil {
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }

        _, err = tx.Exec("INSERT INTO Likes (comment_id, user_id, type) VALUES (?, ?, ?)", commentID, userID, likeType)
        if err != nil {
            tx.Rollback()
            http.Error(w, "Error liking comment", http.StatusInternalServerError)
            return
        }

        _, err = tx.Exec("UPDATE Comments SET score = score + ? WHERE comment_id = ?", scoreChange, commentID)
        if err != nil {
            tx.Rollback()
            http.Error(w, "Error updating comment score", http.StatusInternalServerError)
            return
        }

        err = tx.Commit()
        if err != nil {
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }

        postID := r.FormValue("post_id")
        http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
    }
}

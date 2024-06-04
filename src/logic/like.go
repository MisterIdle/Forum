package logic

import (
	"database/sql"
	"net/http"
)

func LikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		commentID := r.FormValue("comment_id")
		userID := r.FormValue("user_id")
		likeType := r.FormValue("type")
		postID := r.FormValue("post_id")

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

		http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
	}
}

package logic

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postID := r.FormValue("post_id")
		userID := r.FormValue("user_id")
		content := r.FormValue("content")

		files := r.MultipartForm.File["images"]
		var imagePaths []string
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, "Error opening image", http.StatusInternalServerError)
				return
			}
			defer file.Close()

			imagePath := filepath.Join("uploads", fileHeader.Filename)
			f, err := os.Create(imagePath)
			if err != nil {
				http.Error(w, "Error saving image", http.StatusInternalServerError)
				return
			}
			defer f.Close()
			_, err = io.Copy(f, file)
			if err != nil {
				http.Error(w, "Error saving image", http.StatusInternalServerError)
				return
			}
			imagePaths = append(imagePaths, imagePath)
		}

		db, _ := sql.Open("sqlite3", "./Forum3.db")
		defer db.Close()

		_, err := db.Exec("INSERT INTO Comments (post_id, user_id, content, image_paths, score, timestamp) VALUES (?, ?, ?, ?, 0, CURRENT_TIMESTAMP)", postID, userID, content, strings.Join(imagePaths, ","))
		if err != nil {
			http.Error(w, "Error creating comment", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
	}
}

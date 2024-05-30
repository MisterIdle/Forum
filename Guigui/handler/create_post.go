package handlers

import (
    "database/sql"
    "net/http"
    "html/template"
    "strconv"
    "os"
    "path/filepath"
    "strings"
    "io"
    _ "github.com/mattn/go-sqlite3"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        categoryID, _ := strconv.Atoi(r.URL.Query().Get("category_id"))
        tmpl, _ := template.ParseFiles("templates/create_post.html")
        tmpl.Execute(w, PostFormData{
            CategoryID: categoryID,
        })
    } else if r.Method == http.MethodPost {
        categoryID := r.FormValue("category_id")
        userID := r.FormValue("user_id")
        title := r.FormValue("title")
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

        _, err = db.Exec("INSERT INTO Posts (category_id, user_id, title, content, image_paths, timestamp) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)", categoryID, userID, title, content, strings.Join(imagePaths, ","))
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

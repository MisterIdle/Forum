package logic

import (
	"database/sql"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)


func fetchCategories(db *sql.DB) ([]Category, error) {
    rows, err := db.Query("SELECT id, name FROM categories")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []Category
    for rows.Next() {
        var category Category
        if err := rows.Scan(&category.CategoryID, &category.Name); err != nil {
            return nil, err
        }
        categories = append(categories, category)
    }
    return categories, nil
}

func isValidCategory(db *sql.DB, categoryID int) bool {
    var id int
    err := db.QueryRow("SELECT id FROM categories WHERE id = ?", categoryID).Scan(&id)
    return err == nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    t, err := template.ParseFiles(tmpl)
    if err != nil {
        http.Error(w, "Template parsing error", http.StatusInternalServerError)
        return
    }
    if err := t.Execute(w, data); err != nil {
        http.Error(w, "Template executing error", http.StatusInternalServerError)
    }
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
    db, err := sql.Open("sqlite3", "./database.db")
    if err != nil {
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    if r.Method == http.MethodGet {
        categories, err := fetchCategories(db)
        if err != nil {
            http.Error(w, "Error fetching categories", http.StatusInternalServerError)
            return
        }

        categoryID, _ := strconv.Atoi(r.URL.Query().Get("category_id"))
        renderTemplate(w, "templates/create_post.html", PostFormData{CategoryID: categoryID, Categories: categories})
        return
    } else if r.Method == http.MethodPost {
        categoryIDStr := r.FormValue("category_id")
        userID := r.FormValue("user_id")
        title := r.FormValue("title")
        content := r.FormValue("content")

        categoryID, err := strconv.Atoi(categoryIDStr)
        if err != nil || (!isValidCategory(db, categoryID) && categoryID != 0) {
            categories, _ := fetchCategories(db)
            renderTemplate(w, "templates/create_post.html", PostFormData{
                CategoryID: categoryID,
                UserID:     atoi(userID),
                Title:      title,
                Content:    content,
                Error:      "Please select a valid category",
                Categories: categories,
            })
            return
        }

        //image uploads
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

        _, err = db.Exec("INSERT INTO Posts (category_id, user_id, title, content, image_paths, timestamp) VALUES (?, ?, ?, ?, ?, datetime('now'))",
            categoryID, userID, title, content, strings.Join(imagePaths, ","))
        if err != nil {
            http.Error(w, "Error saving post", http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func atoi(s string) int {
    i, _ := strconv.Atoi(s)
    return i
}

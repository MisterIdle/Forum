package handlers

import (
    "database/sql"
    "net/http"
    "html/template"
    "strings"
    _ "github.com/mattn/go-sqlite3"
)

func ViewCategoryPosts(w http.ResponseWriter, r *http.Request) {
    categoryID := r.URL.Query().Get("id")

    db, _ := sql.Open("sqlite3", "./Forum3.db")
    defer db.Close()

    rows, err := db.Query(`
        SELECT p.post_id, p.title, p.content, p.image_paths, 
               (SELECT COUNT(*) FROM Comments c WHERE c.post_id = p.post_id) as comment_count
        FROM Posts p
        WHERE p.category_id = ?
        ORDER BY p.timestamp DESC`, categoryID)
    if err != nil {
        http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []PostPreview
    for rows.Next() {
        var post PostPreview
        var content string
        var imagePaths string
        if err := rows.Scan(&post.PostID, &post.Title, &content, &imagePaths, &post.CommentCount); err != nil {
            http.Error(w, "Error scanning posts", http.StatusInternalServerError)
            return
        }
        post.ContentPreview = getPreview(content, 2)
        if imagePaths != "" {
            post.FirstImage = strings.Split(imagePaths, ",")[0]
        }
        posts = append(posts, post)
    }

    tmpl, _ := template.ParseFiles("templates/view_category_posts.html")
    tmpl.Execute(w, CategoryPostsData{
        CategoryID: atoi(categoryID),
        Posts:      posts,
    })
}

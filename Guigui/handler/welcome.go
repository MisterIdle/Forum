package handlers

import (
    "database/sql"
    "net/http"
    "html/template"
    _ "github.com/mattn/go-sqlite3"
    "strings"
    "sort"
)



func Welcome(w http.ResponseWriter, r *http.Request) {
    db, _ := sql.Open("sqlite3", "./Forum3.db")
    defer db.Close()

    rows, err := db.Query(`SELECT category_id, name FROM Categories`)
    if err != nil {
        http.Error(w, "Error retrieving categories", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var categories []Category
    for rows.Next() {
        var category Category
        if err := rows.Scan(&category.CategoryID, &category.Name); err != nil {
            http.Error(w, "Error scanning categories", http.StatusInternalServerError)
            return
        }
        categories = append(categories, category)
    }

    rows, err = db.Query(`
        SELECT p.post_id, p.title, p.content, p.image_paths, 
               COALESCE(SUM(c.score), 0) as total_score,
               COALESCE(COUNT(DISTINCT c.user_id), 1) as unique_users,
               COUNT(c.comment_id) as comment_count
        FROM Posts p
        LEFT JOIN Comments c ON p.post_id = c.post_id
        GROUP BY p.post_id
        ORDER BY total_score / unique_users DESC, p.timestamp DESC
        LIMIT 10`)
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
        var totalScore float64
        var uniqueUsers int
        if err := rows.Scan(&post.PostID, &post.Title, &content, &imagePaths, &totalScore, &uniqueUsers, &post.CommentCount); err != nil {
            http.Error(w, "Error scanning posts", http.StatusInternalServerError)
            return
        }
        post.ContentPreview = getPreview(content, 2)
        if imagePaths != "" {
            post.FirstImage = strings.Split(imagePaths, ",")[0]
        }
        post.PopularityScore = totalScore / float64(uniqueUsers)
        posts = append(posts, post)
    }

    sort.Slice(posts, func(i, j int) bool {
        return posts[i].PopularityScore > posts[j].PopularityScore
    })

    tmpl, _ := template.ParseFiles("templates/welcome.html")
    tmpl.Execute(w, WelcomeData{
        Categories: categories,
        Posts:      posts,
    })
}

func getPreview(content string, lines int) string {
    linesArray := strings.Split(content, "\n")
    if len(linesArray) < lines {
        return content
    }
    return strings.Join(linesArray[:lines], "\n")
}

package logic

import (
	"database/sql"
	"html/template"
	"net/http"
	"sort"
	"strings"
)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "./database.db")
	defer db.Close()

	//category
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

	// Retrieve post statistics
	var postCount, commentCount, uniqueUserCount int
	var latestMember string

	db.QueryRow(`SELECT COUNT(*) FROM Posts`).Scan(&postCount)
	db.QueryRow(`SELECT COUNT(*) FROM Comments`).Scan(&commentCount)
	db.QueryRow(`SELECT COUNT(DISTINCT user_id) FROM Users`).Scan(&uniqueUserCount)
	db.QueryRow(`SELECT username FROM Users ORDER BY user_id DESC LIMIT 1`).Scan(&latestMember)

	// Retrieve posts with score calculation
	rows, err = db.Query(`
        SELECT p.post_id, p.title, p.content, p.image_paths, 
               COALESCE(SUM(c.score), 0) as total_score,
               COALESCE(COUNT(DISTINCT c.user_id), 0) as unique_users,
               COUNT(c.comment_id) as comment_count
        FROM Posts p
        LEFT JOIN Comments c ON p.post_id = c.post_id
        GROUP BY p.post_id
        ORDER BY total_score DESC, p.timestamp DESC
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
		if uniqueUsers > 0 {
			post.PopularityScore = totalScore / float64(uniqueUsers)
		} else {
			post.PopularityScore = 0
		}
		posts = append(posts, post)
	}

	// Sort posts by score
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PopularityScore > posts[j].PopularityScore
	})

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, WelcomeData{
		Categories:      categories,
		Posts:           posts,
		PostCount:       postCount,
		CommentCount:    commentCount,
		UniqueUserCount: uniqueUserCount,
		LatestMember:    latestMember,
	})
}

func getPreview(content string, lines int) string {
	linesArray := strings.Split(content, "\n")
	if len(linesArray) < lines {
		return content
	}
	return strings.Join(linesArray[:lines], "\n")
}

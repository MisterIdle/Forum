package logic

import (
	"html/template"
	"net/http"
	"strings"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	var searchResults SearchResult

	if query != "" {
		// Search by posts
		postRows, err := db.Query(`
			SELECT p.post_id, p.title, c.name AS category_name, COALESCE(COUNT(DISTINCT cm.user_id), 0) AS unique_users
			FROM Posts p
			JOIN Categories c ON p.category_id = c.category_id
			LEFT JOIN Comments cm ON p.post_id = cm.post_id
			WHERE p.title LIKE ? OR p.content LIKE ?
			GROUP BY p.post_id, c.name
			LIMIT 10`, "%"+strings.ToLower(query)+"%", "%"+strings.ToLower(query)+"%")
		if err != nil {
			http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
			return
		}
		defer postRows.Close()

		for postRows.Next() {
			var post PostResult
			if err := postRows.Scan(&post.PostID, &post.Title, &post.CategoryName, &post.UniqueUsers); err != nil {
				http.Error(w, "Error scanning posts", http.StatusInternalServerError)
				return
			}
			searchResults.Posts = append(searchResults.Posts, post)
		}

		// Search by users
		userRows, err := db.Query(`SELECT user_id, username FROM Users WHERE username LIKE ? LIMIT 10`, "%"+strings.ToLower(query)+"%")
		if err != nil {
			http.Error(w, "Error retrieving users", http.StatusInternalServerError)
			return
		}
		defer userRows.Close()

		for userRows.Next() {
			var user UserResult
			if err := userRows.Scan(&user.UserID, &user.Username); err != nil {
				http.Error(w, "Error scanning users", http.StatusInternalServerError)
				return
			}
			searchResults.Users = append(searchResults.Users, user)
		}
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, PageData{
		SearchResults: searchResults,
		Query:         query,
	})
}

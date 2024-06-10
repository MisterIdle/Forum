package logic

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitData() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	reset := flag.Bool("reset", false, "Reset the database")
	force := flag.Bool("force", false, "Force the database reset")
	flag.Parse()

	if *force {
		os.Remove("./database.db")
		createData()
	}

	if *reset {
		resetUsers()
		resetCategories()
		resetPosts()
		resetComments()
		createBasicCategories()

		fmt.Println("Database reset")
	}

	fmt.Println("Database initialized")
}

func createData() {
	query := `
	CREATE TABLE IF NOT EXISTS Ranks (
		rank_id INTEGER PRIMARY KEY,
		rank_name VARCHAR
	);

	CREATE TABLE IF NOT EXISTS Users (
		uuid TEXT PRIMARY KEY,
		username VARCHAR,
		email VARCHAR,
		password VARCHAR,
		identity TEXT,
		code TEXT,
		creation DATETIME,
		rank_id INTEGER,
		picture VARCHAR,
		FOREIGN KEY (rank_id) REFERENCES Ranks(rank_id)
	);

	CREATE TABLE IF NOT EXISTS Categories (
		category_id INTEGER PRIMARY KEY,
		name VARCHAR,
		description TEXT,
		global TEXT
	);

	CREATE TABLE IF NOT EXISTS Posts (
		post_id INTEGER PRIMARY KEY,
		title TEXT,
		content TEXT,
		timestamp DATETIME,
		category_id INTEGER
	);
	
	CREATE TABLE IF NOT EXISTS Comments (
		comment_id INTEGER PRIMARY KEY,
		content TEXT,
		timestamp DATETIME,
		post_id INTEGER
	);`

	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
}

func resetUsers() {
	query := `DELETE FROM Users;`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
}

func resetCategories() {
	query := `DELETE FROM Categories;`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
}

func resetPosts() {
	query := `DELETE FROM Posts;`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
}

func resetComments() {
	query := `DELETE FROM Comments;`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
}

func checkUserMailOrIdentidy(identifier string) bool {
	query := `SELECT COALESCE(email, identity) FROM Users WHERE email = ? OR identity = ?;`
	row := db.QueryRow(query, identifier, identifier)
	var result string
	err := row.Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

func getUserInfoByMailOrIdentidy(identifier string) (string, string, int, string) {
	query := `SELECT uuid, creation, rank_id, picture FROM Users WHERE email = ? OR identity = ?;`
	row := db.QueryRow(query, identifier, identifier)
	var uuid, creation, picture string
	var rank_id int
	err := row.Scan(&uuid, &creation, &rank_id, &picture)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", 0, ""
		}
		fmt.Println(err)
		return "", "", 0, ""
	}
	return uuid, creation, rank_id, picture
}

func setCodeByEmail(email, code string) error {
	query := `UPDATE Users SET code = ? WHERE email = ?;`
	_, err := db.Exec(query, code, email)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func getEmailFromCode(code string) string {
	query := `SELECT email FROM Users WHERE code = ?;`
	row := db.QueryRow(query, code)
	var email string
	err := row.Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		}
		fmt.Println(err)
		return ""
	}
	return email
}

func resetPassword(email, password string) error {
	query := `UPDATE Users SET password = ? WHERE email = ?;`
	_, err := db.Exec(query, password, email)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func resetCode(email string) error {
	query := `UPDATE Users SET code = null WHERE email = ?;`
	_, err := db.Exec(query, email)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func getCredentialsByUsernameOrEmail(identifier string) (string, string) {
	query := `SELECT password, COALESCE(username, email) FROM Users WHERE (username = ? OR email = ?) AND identity = 'LOCAL';`
	row := db.QueryRow(query, identifier, identifier)
	var password, username string
	err := row.Scan(&password, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ""
		}
		fmt.Println(err)
		return "", ""
	}
	return password, username
}

func newUser(username, email, password, identity, picture string) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		return err
	}

	query := `INSERT INTO Users (uuid, username, email, password, identity, code, creation, rank_id, picture) VALUES (?, ?, ?, ?, ?, null, datetime('now'), 1, ?);`
	_, err = db.Exec(query, uuid.String(), username, email, password, identity, picture)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// Category

func createCategory(name, description, global string) error {
	query := `INSERT INTO Categories (name, description, global) VALUES (?, ?, ?);`
	_, err := db.Exec(query, name, description, global)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func createBasicCategories() {
	createCategory("General", "General discussion", "Forum général")
	createCategory("Windows", "Windows discussion", "Informatique")
	createCategory("Linux", "Linux discussion", "Informatique")
	createCategory("Mac", "Mac discussion", "Informatique")
	createCategory("Golang", "Golang discussion", "Programmation")
	createCategory("Python", "Python discussion", "Programmation")
	createCategory("Java", "Java discussion", "Programmation")
}

func fetchCategories() (map[string][]Category, error) {
	query := `SELECT category_id, name, global FROM Categories;`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	categories := make(map[string][]Category)
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.CategoryID, &category.Name, &category.Global); err != nil {
			fmt.Println(err)
			return nil, err
		}
		categories[category.Global] = append(categories[category.Global], category)
	}

	return categories, nil
}

func getCategoryName(categoryID int) string {
	query := `SELECT name FROM Categories WHERE category_id = ?;`
	row := db.QueryRow(query, categoryID)
	var name string
	err := row.Scan(&name)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return name
}

func getPostsByCategoryID(categoryID int) []Post {
	query := `SELECT post_id, title, content, timestamp FROM Posts WHERE category_id = ? ORDER BY timestamp DESC;`
	rows, err := db.Query(query, categoryID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostID, &post.Title, &post.Content, &post.Timestamp); err != nil {
			fmt.Println(err)
			return nil
		}
		posts = append(posts, post)
	}
	return posts
}

func newPost(categoryID int, title, content string) (int, error) {
	query := `INSERT INTO Posts (title, content, timestamp, category_id) VALUES (?, ?, datetime('now'), ?);`
	result, err := db.Exec(query, title, content, categoryID)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return int(postID), nil
}

// Post

func fetchPostByID(postID int) (Post, error) {
	query := `SELECT title, content, timestamp FROM Posts WHERE post_id = ?;`
	row := db.QueryRow(query, postID)
	var post Post
	err := row.Scan(&post.Title, &post.Content, &post.Timestamp)
	if err != nil {
		fmt.Println(err)
		return Post{}, err
	}
	return post, nil
}

// Comment

func getCommentsByPostID(postID int) []Comment {
	query := `SELECT comment_id, content, timestamp FROM Comments WHERE post_id = ? ORDER BY timestamp DESC;`
	rows, err := db.Query(query, postID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.CommentID, &comment.Content, &comment.Timestamp); err != nil {
			fmt.Println(err)
			return nil
		}
		comments = append(comments, comment)
	}
	return comments
}

func newComment(postID int, content string) error {
	query := `INSERT INTO Comments (content, timestamp, post_id) VALUES (?, datetime('now'), ?);`
	_, err := db.Exec(query, content, postID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

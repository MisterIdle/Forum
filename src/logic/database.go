package logic

import (
	"database/sql"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Initialization functions
func InitData() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		return
	}

	reset := flag.Bool("reset", false, "Reset the database")
	force := flag.Bool("force", false, "Force the database reset")
	flag.Parse()

	if *force {
		os.Remove("./database.db")
		createData()
		resetAll()
		createBasicCategories()
		createBasicRanks()
		createAdminUser()

		fmt.Println("Database has been removed and reset")
	}

	if *reset {
		resetAll()
		createBasicCategories()
		createBasicRanks()
		createAdminUser()

		fmt.Println("Database has been reset")
	}

	fmt.Println("Database has been initialized")
}

func createData() {
	query := `
    CREATE TABLE IF NOT EXISTS Users (
        user_id INTEGER PRIMARY KEY,
        uuid TEXT UNIQUE,
        username VARCHAR,
        email VARCHAR UNIQUE,
        password VARCHAR,
        creation DATETIME,
        rank_id INTEGER,
        picture VARCHAR
    );

	CREATE TABLE IF NOT EXISTS Ranks (
        rank_id INTEGER PRIMARY KEY,
        rank_name VARCHAR
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
		username TEXT,
        timestamp DATETIME,
        category_id INTEGER
    );

    CREATE TABLE IF NOT EXISTS Likes (
        like_id INTEGER PRIMARY KEY,
        post_id INTEGER,
		comment_id INTEGER,
        user_id INTEGER
    );

    CREATE TABLE IF NOT EXISTS Dislikes (
        dislike_id INTEGER PRIMARY KEY,
        post_id INTEGER,
		comment_id INTEGER,
        user_id INTEGER
    );

	CREATE TABLE IF NOT EXISTS Images (
	    image_id INTEGER PRIMARY KEY,
	    post_id INTEGER,
	    image_name TEXT,
	    FOREIGN KEY (post_id) REFERENCES Posts(post_id)
	);

    CREATE TABLE IF NOT EXISTS Comments (
        comment_id INTEGER PRIMARY KEY,
        content TEXT,
        timestamp DATETIME,
		username TEXT,
        post_id INTEGER
    );`

	db.Exec(query)
}

func resetAll() {
	resetUsers()
	resetCategories()
	resetPosts()
	resetComments()
	resetLikes()
	resetDislikes()
	resetImages()
	resetRanks()
}

// Reset functions
func resetUsers() {
	query := `DELETE FROM Users;`
	db.Exec(query)
}

func resetCategories() {
	query := `DELETE FROM Categories;`
	db.Exec(query)
}

func resetPosts() {
	query := `DELETE FROM Posts;`
	db.Exec(query)
}

func resetComments() {
	query := `DELETE FROM Comments;`
	db.Exec(query)
}

func resetLikes() {
	query := `DELETE FROM Likes;`
	db.Exec(query)
}

func resetDislikes() {
	query := `DELETE FROM Dislikes;`
	db.Exec(query)
}

func resetRanks() {
	query := `DELETE FROM Ranks;`
	db.Exec(query)
}

func resetImages() {
	query := `DELETE FROM Images;`
	db.Exec(query)

	resetImageFolder("./img/upload/")
	resetImageFolder("./img/profile/")
}

func resetImageFolder(folder string) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.Name() != "Default.png" {
			os.Remove(folder + file.Name())
		}
	}
}

// User functions
func createAdminUser() {
	password := newRandomPassword()
	newUser("Admin", "Admin", hashedPassword(password), "Default.png", 3)
	fmt.Println("Admin password: ", password)
}

func newRandomPassword() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖabcdefghijklmnopqrstuvwxyzåäö0123456789")
	length := 10
	var password strings.Builder
	for i := 0; i < length; i++ {
		password.WriteRune(chars[rand.Intn(len(chars))])
	}

	os.Create("password.txt")
	file, err := os.OpenFile("password.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return ""
	}
	defer file.Close()

	file.WriteString(password.String() + "\n")

	return password.String()
}

func checkUserEmail(email string) bool {
	query := `SELECT email FROM Users WHERE email = ?;`
	row := db.QueryRow(query, email)
	var result string
	err := row.Scan(&result)
	if err != nil {
		return false
	}
	return true
}

func checkUserUsername(username string) bool {
	query := `SELECT username FROM Users WHERE username = ?;`
	row := db.QueryRow(query, username)
	var result string
	err := row.Scan(&result)
	if err != nil {
		return false
	}
	return true
}

func getCredentialsByEmail(email string) (string, string) {
	query := `SELECT password, COALESCE(username, email) FROM Users WHERE email = ?;`
	row := db.QueryRow(query, email)
	var password, username string
	err := row.Scan(&password, &username)
	if err != nil {
		return "", ""
	}
	return password, username
}

func getCredentialsByUsername(username string) (string, string) {
	query := `SELECT password, COALESCE(username, email) FROM Users WHERE username = ?;`
	row := db.QueryRow(query, username)
	var password, email string
	err := row.Scan(&password, &email)
	if err != nil {
		return "", ""
	}
	return password, email
}

func getIDByUUID(uuid string) int {
	query := `SELECT user_id FROM Users WHERE uuid = ?;`
	row := db.QueryRow(query, uuid)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0
	}
	return id
}

func getUsernameByUUID(uuid string) string {
	query := `SELECT username FROM Users WHERE uuid = ?;`
	row := db.QueryRow(query, uuid)
	var username string
	err := row.Scan(&username)
	if err != nil {
		return ""
	}
	return username
}

func getUUIDByUsername(username string) string {
	query := `SELECT uuid FROM Users WHERE username = ?;`
	row := db.QueryRow(query, username)
	var uuid string
	err := row.Scan(&uuid)
	if err != nil {
		return ""
	}
	return uuid
}

func newUser(username, email, password, picture string, rankID int) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	query := `INSERT INTO Users (uuid, username, email, password, creation, rank_id, picture) VALUES (?, ?, ?, ?, datetime('now'), ?, ?);`
	_, err = db.Exec(query, uuid.String(), username, email, password, rankID, picture)
	if err != nil {
		return err
	}
	return nil
}

func changeProfileUsername(username, uuid string) error {
	query := `UPDATE Posts SET username = ? WHERE username = (SELECT username FROM Users WHERE uuid = ?);`
	_, err := db.Exec(query, username, uuid)
	if err != nil {
		return err
	}
	query = `UPDATE Comments SET username = ? WHERE username = (SELECT username FROM Users WHERE uuid = ?);`
	_, err = db.Exec(query, username, uuid)
	if err != nil {
		return err
	}
	query = `UPDATE Users SET username = ? WHERE uuid = ?;`
	_, err = db.Exec(query, username, uuid)
	if err != nil {
		return err
	}
	return nil
}

func changeProfilePassword(password, uuid string) error {
	query := `UPDATE Users SET password = ? WHERE uuid = ?;`
	_, err := db.Exec(query, password, uuid)
	if err != nil {
		return err
	}
	return nil
}

func changeProfileEmail(email, uuid string) error {
	query := `UPDATE Users SET email = ? WHERE uuid = ?;`
	_, err := db.Exec(query, email, uuid)
	if err != nil {
		return err
	}
	return nil
}

func changeProfilePicture(picture, uuid string) error {
	query := `UPDATE Users SET picture = ? WHERE uuid = ?;`
	_, err := db.Exec(query, picture, uuid)
	if err != nil {
		return err
	}
	return nil
}

// Category functions
func createCategory(name, description, global string) error {
	query := `INSERT INTO Categories (name, description, global) VALUES (?, ?, ?);`
	_, err := db.Exec(query, name, description, global)
	if err != nil {
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

func fetchGlobalCategories() (map[string][]Category, error) {
	query := `SELECT category_id, name, description, global, (SELECT COUNT(*) FROM Posts WHERE category_id = c.category_id) AS total_posts, (SELECT COUNT(*) FROM Comments WHERE post_id IN (SELECT post_id FROM Posts WHERE category_id = c.category_id)) AS total_comments FROM Categories c ORDER BY global;`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	categories := make(map[string][]Category)
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.CategoryID, &category.Name, &category.Description, &category.Global, &category.TotalPosts, &category.TotalComments); err != nil {
			return nil, err
		}
		categories[category.Global] = append(categories[category.Global], category)
	}

	return categories, nil
}

func fetchGlobalCategoriesByName(global string) (map[string][]Category, error) {
	query := `SELECT category_id, name, description, global, (SELECT COUNT(*) FROM Posts WHERE category_id = c.category_id) AS total_posts, (SELECT COUNT(*) FROM Comments WHERE post_id IN (SELECT post_id FROM Posts WHERE category_id = c.category_id)) AS total_comments FROM Categories c WHERE global = ? ORDER BY global;`
	rows, err := db.Query(query, global)
	if err != nil {
		return nil, err
	}

	categories := make(map[string][]Category)
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.CategoryID, &category.Name, &category.Description, &category.Global, &category.TotalPosts, &category.TotalComments); err != nil {
			return nil, err
		}
		categories[category.Global] = append(categories[category.Global], category)
	}

	return categories, nil
}

func fetchCategoriesName() []string {
	query := `SELECT name FROM Categories;`
	rows, err := db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil
		}
		names = append(names, name)
	}

	return names
}

func fetchGlobalCategoriesName() []string {
	query := `SELECT global FROM Categories;`
	rows, err := db.Query(query)

	if err != nil {
		return nil
	}

	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil
		}
		names = append(names, name)
	}

	names = removeDuplicates(names)

	return names
}

func removeDuplicates(names []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range names {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

func checkCategoryName(name string) bool {
	query := `SELECT name FROM Categories WHERE name = ?;`
	row := db.QueryRow(query, name)
	var result string
	err := row.Scan(&result)
	if err != nil {
		return false
	}
	return true
}

func getCategoryName(categoryID int) string {
	query := `SELECT name FROM Categories WHERE category_id = ?;`
	row := db.QueryRow(query, categoryID)
	var name string
	err := row.Scan(&name)
	if err != nil {
		return ""
	}
	return name
}

func getCategoryDescription(categoryID int) string {
	query := `SELECT description FROM Categories WHERE category_id = ?;`
	row := db.QueryRow(query, categoryID)
	var description string
	err := row.Scan(&description)
	if err != nil {
		return ""
	}
	return description
}

func deleteCategory(categoryName string) error {
	categoryID := getCategoryIDByName(categoryName)

	posts := getPostsByCategoryID(categoryID)
	for _, post := range posts {
		deleteImageByPostID(post.PostID)
	}

	query := `DELETE FROM Likes WHERE post_id IN (SELECT post_id FROM Posts WHERE category_id = ?);`
	if _, err := db.Exec(query, categoryID); err != nil {
		return err
	}

	query = `DELETE FROM Dislikes WHERE post_id IN (SELECT post_id FROM Posts WHERE category_id = ?);`
	if _, err := db.Exec(query, categoryID); err != nil {
		return err
	}

	query = `DELETE FROM Comments WHERE post_id IN (SELECT post_id FROM Posts WHERE category_id = ?);`
	if _, err := db.Exec(query, categoryID); err != nil {
		return err
	}

	query = `DELETE FROM Posts WHERE category_id = ?;`
	if _, err := db.Exec(query, categoryID); err != nil {
		return err
	}

	query = `DELETE FROM Categories WHERE category_id = ?;`
	if _, err := db.Exec(query, categoryID); err != nil {
		return err
	}

	return nil
}

func getCategoryIDByName(categoryName string) int {
	query := `SELECT category_id FROM Categories WHERE name = ?;`
	row := db.QueryRow(query, categoryName)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0
	}
	return id
}

// Post functions
func fetchPostByID(postID int) (Posts, error) {
	query := `SELECT title, content, timestamp, username, (SELECT COUNT(*) FROM Likes WHERE post_id = ?) AS likes, (SELECT COUNT(*) FROM Dislikes WHERE post_id = ?) AS dislikes FROM Posts WHERE post_id = ?;`
	row := db.QueryRow(query, postID, postID, postID)
	var post Posts
	err := row.Scan(&post.Title, &post.Content, &post.Timestamp, &post.Username, &post.LikesPost, &post.DislikesPost)
	if err != nil {
		return Posts{}, err
	}

	return post, nil
}

func checkPostTitle(title string) bool {
	query := `SELECT title FROM Posts WHERE title = ?;`
	row := db.QueryRow(query, title)
	var result string
	err := row.Scan(&result)
	if err != nil {
		return false
	}
	return true
}

func getPostsByCategoryID(categoryID int) []Posts {
	query := `SELECT post_id, title, content, timestamp, username FROM Posts WHERE category_id = ? ORDER BY timestamp DESC;`
	rows, err := db.Query(query, categoryID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var posts []Posts
	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.PostID, &post.Title, &post.Content, &post.Timestamp, &post.Username); err != nil {
			return nil
		}
		posts = append(posts, post)
	}

	return posts
}

func getPostTotalsByCategoryID(categoryID int) int {
	query := `SELECT COUNT(*) FROM Posts WHERE category_id = ?;`
	row := db.QueryRow(query, categoryID)
	var total int
	err := row.Scan(&total)
	if err != nil {
		return 0
	}
	return total
}

func newPost(categoryID int, title, content, username string) (int, error) {
	query := `INSERT INTO Posts (title, content, timestamp, category_id, username) VALUES (?, ?, datetime('now'), ?, ?);`
	result, err := db.Exec(query, title, content, categoryID, username)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func deletePost(postID int) error {
	deleteImageByPostID(postID)

	query := `DELETE FROM Likes WHERE post_id = ?;`
	if _, err := db.Exec(query, postID); err != nil {
		return err
	}

	query = `DELETE FROM Dislikes WHERE post_id = ?;`
	if _, err := db.Exec(query, postID); err != nil {
		return err
	}

	query = `DELETE FROM Comments WHERE post_id = ?;`
	if _, err := db.Exec(query, postID); err != nil {
		return err
	}

	query = `DELETE FROM Posts WHERE post_id = ?;`
	if _, err := db.Exec(query, postID); err != nil {
		return err
	}

	return nil
}

func getLikesByPostID(postID int) int {
	query := `SELECT COUNT(*) FROM Likes WHERE post_id = ?;`
	row := db.QueryRow(query, postID)
	var likes int
	err := row.Scan(&likes)
	if err != nil {
		return 0
	}
	return likes
}

func getDislikesByPostID(postID int) int {
	query := `SELECT COUNT(*) FROM Dislikes WHERE post_id = ?;`
	row := db.QueryRow(query, postID)
	var dislikes int
	err := row.Scan(&dislikes)
	if err != nil {
		return 0
	}
	return dislikes
}

func getCategoryNameByPostID(postID int) string {
	query := `SELECT name FROM Categories WHERE category_id = (SELECT category_id FROM Posts WHERE post_id = ?);`
	row := db.QueryRow(query, postID)
	var name string
	err := row.Scan(&name)
	if err != nil {
		return ""
	}
	return name
}

func getCategoryIDByPostID(postID int) int {
	query := `SELECT category_id FROM Posts WHERE post_id = ?;`
	row := db.QueryRow(query, postID)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0
	}
	return id
}

func getUsernameByPostID(postID int) string {
	query := `SELECT username FROM Posts WHERE post_id = ?;`
	row := db.QueryRow(query, postID)
	var username string
	err := row.Scan(&username)
	if err != nil {
		return ""
	}
	return username
}

func hasUserLikedPost(postID, userID int) bool {
	query := `SELECT COUNT(*) FROM Likes WHERE post_id = ? AND user_id = ?;`
	row := db.QueryRow(query, postID, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func newLikePost(postID, userID int) error {
	query := `INSERT INTO Likes (post_id, user_id) VALUES (?, ?);`
	_, err := db.Exec(query, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func hasUserDislikedPost(postID, userID int) bool {
	query := `SELECT COUNT(*) FROM Dislikes WHERE post_id = ? AND user_id = ?;`
	row := db.QueryRow(query, postID, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func newDislikePost(postID, userID int) error {
	query := `INSERT INTO Dislikes (post_id, user_id) VALUES (?, ?);`
	_, err := db.Exec(query, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func removeDislikePost(postID, userID int) error {
	query := `DELETE FROM Dislikes WHERE post_id = ? AND user_id = ?;`
	_, err := db.Exec(query, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func removeLikePost(postID, userID int) error {
	query := `DELETE FROM Likes WHERE post_id = ? AND user_id = ?;`
	_, err := db.Exec(query, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

// Comment functions
func fetchCommentsByPostID(postID int) ([]Comments, error) {
	query := `SELECT comment_id, content, timestamp, username FROM Comments WHERE post_id = ? ORDER BY timestamp DESC;`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comments
	for rows.Next() {
		var comment Comments
		if err := rows.Scan(&comment.CommentID, &comment.Content, &comment.Timestamp, &comment.Username); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func newComment(postID int, content, username string) error {
	query := `INSERT INTO Comments (content, timestamp, username, post_id) VALUES (?, datetime('now'), ?, ?);`
	_, err := db.Exec(query, content, username, postID)
	if err != nil {
		return err
	}
	return nil
}

func deleteComment(commentID int) error {
	query := `DELETE FROM Likes WHERE comment_id = ?;`
	if _, err := db.Exec(query, commentID); err != nil {
		return err
	}

	query = `DELETE FROM Dislikes WHERE comment_id = ?;`
	if _, err := db.Exec(query, commentID); err != nil {
		return err
	}

	query = `DELETE FROM Comments WHERE comment_id = ?;`
	if _, err := db.Exec(query, commentID); err != nil {
		return err
	}
	return nil
}

func getUsernameByCommentID(commentID int) string {
	query := `SELECT username FROM Comments WHERE comment_id = ?;`
	row := db.QueryRow(query, commentID)
	var username string
	err := row.Scan(&username)
	if err != nil {
		return ""
	}
	return username
}

func getLikesByCommentID(commentID int) int {
	query := `SELECT COUNT(*) FROM Likes WHERE comment_id = ?;`
	row := db.QueryRow(query, commentID)
	var likes int
	err := row.Scan(&likes)
	if err != nil {
		return 0
	}
	return likes
}

func getDislikesByCommentID(commentID int) int {
	query := `SELECT COUNT(*) FROM Dislikes WHERE comment_id = ?;`
	row := db.QueryRow(query, commentID)
	var dislikes int
	err := row.Scan(&dislikes)
	if err != nil {
		return 0
	}
	return dislikes
}

func hasUserLikedComment(commentID, userID int) bool {
	query := `SELECT COUNT(*) FROM Likes WHERE comment_id = ? AND user_id = ?;`
	row := db.QueryRow(query, commentID, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func newLikeComment(commentID, userID int) error {
	query := `INSERT INTO Likes (comment_id, user_id) VALUES (?, ?);`
	_, err := db.Exec(query, commentID, userID)
	if err != nil {
		return err
	}
	return nil
}

func hasUserDislikedComment(commentID, userID int) bool {
	query := `SELECT COUNT(*) FROM Dislikes WHERE comment_id = ? AND user_id = ?;`
	row := db.QueryRow(query, commentID, userID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func newDislikeComment(commentID, userID int) error {
	query := `INSERT INTO Dislikes (comment_id, user_id) VALUES (?, ?);`
	_, err := db.Exec(query, commentID, userID)
	if err != nil {
		return err
	}
	return nil
}

func removeDislikeComment(commentID, userID int) error {
	query := `DELETE FROM Dislikes WHERE comment_id = ? AND user_id = ?;`
	_, err := db.Exec(query, commentID, userID)
	if err != nil {
		return err
	}
	return nil
}

func removeLikeComment(commentID, userID int) error {
	query := `DELETE FROM Likes WHERE comment_id = ? AND user_id = ?;`
	_, err := db.Exec(query, commentID, userID)
	if err != nil {
		return err
	}
	return nil
}

// Image functions
func uploadImage(postID int, imageName string) error {
	query := `INSERT INTO Images (post_id, image_name) VALUES (?, ?);`
	_, err := db.Exec(query, postID, imageName)
	if err != nil {
		return err
	}
	return nil
}

func getImagesByPostID(postID int) []string {
	query := `SELECT image_name FROM Images WHERE post_id = ?;`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil
	}

	var images []string
	for rows.Next() {
		var image string
		if err := rows.Scan(&image); err != nil {
			return nil
		}
		images = append(images, image)
	}
	return images
}

func deleteImageByPostID(postID int) error {
	query := `DELETE FROM Images WHERE post_id = ?;`
	_, err := db.Exec(query, postID)
	if err != nil {
		return err
	}

	resetPostImages(postID)

	return nil
}

func resetPostImages(postID int) {
	files, err := os.ReadDir("./img/upload/")
	if err != nil {
		return
	}

	for _, file := range files {
		if file.Name()[:len(fmt.Sprint(postID))] == fmt.Sprint(postID) {
			os.Remove("./img/upload/" + file.Name())
		}
	}
}

// Profile functions
func fetchAllUsernames() []string {
	query := `SELECT username FROM Users;`
	rows, err := db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil
		}
		usernames = append(usernames, username)
	}
	return usernames
}

func fetchProfile(username string) (Profile, error) {
	query := `SELECT username, uuid, (SELECT rank_name FROM Ranks WHERE rank_id = (SELECT rank_id FROM Users WHERE username = ?)), creation, (SELECT COUNT(*) FROM Posts WHERE username = ?), (SELECT COUNT(*) FROM Comments WHERE username = ?), (SELECT COUNT(*) FROM Likes WHERE user_id = (SELECT user_id FROM Users WHERE username = ?)), (SELECT COUNT(*) FROM Dislikes WHERE user_id = (SELECT user_id FROM Users WHERE username = ?)) FROM Users WHERE username = ?;`
	row := db.QueryRow(query, username, username, username, username, username, username)
	var profile Profile
	err := row.Scan(&profile.Username, &profile.UUID, &profile.Rank, &profile.Timestamp, &profile.TotalPosts, &profile.TotalComments, &profile.TotalLikes, &profile.TotalDislikes)
	if err != nil {
		return Profile{}, err
	}
	return profile, nil
}

func fetchProfilePosts(username string) []Posts {
	query := `SELECT post_id, title, content, timestamp, (SELECT name FROM Categories WHERE category_id = (SELECT category_id FROM Posts WHERE post_id = p.post_id)) FROM Posts p WHERE username = ? ORDER BY timestamp DESC;`
	rows, err := db.Query(query, username)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var posts []Posts
	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.PostID, &post.Title, &post.Content, &post.Timestamp, &post.CategoryName); err != nil {
			return nil
		}
		posts = append(posts, post)
	}
	return posts
}

func fetchProfileComments(username string) []Comments {
	query := `SELECT comment_id, content, timestamp, (SELECT title FROM Posts WHERE post_id = c.post_id), post_id FROM Comments c WHERE username = ? ORDER BY timestamp DESC;`
	rows, err := db.Query(query, username)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var comments []Comments
	for rows.Next() {
		var comment Comments
		if err := rows.Scan(&comment.CommentID, &comment.Content, &comment.Timestamp, &comment.Title, &comment.PostID); err != nil {
			return nil
		}
		comments = append(comments, comment)
	}
	return comments
}

func getProfilePictureByUUID(uuid string) string {
	query := `SELECT picture FROM Users WHERE uuid = ?;`
	row := db.QueryRow(query, uuid)
	var picture string
	err := row.Scan(&picture)
	if err != nil {
		return ""
	}
	return picture
}

func deleteProfile(uuid string) error {
	query := `DELETE FROM Likes WHERE user_id = (SELECT user_id FROM Users WHERE uuid = ?);`
	if _, err := db.Exec(query, uuid); err != nil {
		return err
	}

	query = `DELETE FROM Dislikes WHERE user_id = (SELECT user_id FROM Users WHERE uuid = ?);`
	if _, err := db.Exec(query, uuid); err != nil {
		return err
	}

	query = `DELETE FROM Comments WHERE username = (SELECT username FROM Users WHERE uuid = ?);`
	if _, err := db.Exec(query, uuid); err != nil {
		return err
	}

	query = `DELETE FROM Posts WHERE username = (SELECT username FROM Users WHERE uuid = ?);`
	if _, err := db.Exec(query, uuid); err != nil {
		return err
	}

	query = `DELETE FROM Users WHERE uuid = ?;`
	if _, err := db.Exec(query, uuid); err != nil {
		return err
	}

	picture := getProfilePictureByUUID(uuid)
	if picture != "Default.png" {
		os.Remove("./img/profile/" + picture)
	}

	return nil
}

// Rank functions
func createRank(name string) error {
	query := `INSERT INTO Ranks (rank_name) VALUES (?);`
	_, err := db.Exec(query, name)
	if err != nil {
		return err
	}
	return nil
}

func getRankByUUID(uuid string) string {
	query := `SELECT rank_name FROM Ranks WHERE rank_id = (SELECT rank_id FROM Users WHERE uuid = ?);`
	row := db.QueryRow(query, uuid)
	var rank string
	err := row.Scan(&rank)
	if err != nil {
		return ""
	}
	return rank
}

func promoteUser(uuid string) error {
	query := `UPDATE Users SET rank_id = (SELECT rank_id FROM Users WHERE uuid = ?) + 1 WHERE uuid = ? AND (SELECT rank_id FROM Users WHERE uuid = ?) < 3;`
	_, err := db.Exec(query, uuid, uuid, uuid)
	if err != nil {
		return err
	}
	return nil
}

func demoteUser(uuid string) error {
	query := `UPDATE Users SET rank_id = (SELECT rank_id FROM Users WHERE uuid = ?) - 1 WHERE uuid = ? AND (SELECT rank_id FROM Users WHERE uuid = ?) > 1;`
	_, err := db.Exec(query, uuid, uuid, uuid)
	if err != nil {
		return err
	}
	return nil
}

func createBasicRanks() {
	createRank("User")
	createRank("Moderator")
	createRank("Administrator")
}

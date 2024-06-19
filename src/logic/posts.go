// Code by MisterIdle
// Try Pistol Hand on Itch.io: https://misteridle.itch.io/pistol-hand :D

package logic

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// getPostsData retrieves the post data by ID along with its comments
func getPostsData(id int, session Session) (Posts, error) {
	post, err := fetchPostByID(id)
	if err != nil {
		return Posts{}, err
	}

	comments, err := getCommentsDataByPostID(id, session)
	if err != nil {
		return Posts{}, err
	}

	data := Posts{
		PostID:       id,
		Title:        post.Title,
		Content:      post.Content,
		Timestamp:    post.Timestamp,
		Username:     getUsernameByPostID(id),
		LikesPost:    getLikesByPostID(id),
		DislikesPost: getDislikesByPostID(id),
		Images:       getImagesByPostID(id),
		CategoryName: getCategoryNameByPostID(id),
		CategoryID:   getCategoryIDByPostID(id),
		Comments:     comments,
	}

	return data, nil
}

// PostsHandler handles requests for displaying posts
func PostsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		errorPage(w, r) // Show error page if post ID is invalid
		return
	}

	session := getActiveSession(r) // Get the current session

	data, err := getPostsData(id, session) // Fetch post data
	if err != nil {
		errorPage(w, r) // Show error page if data fetching fails
		return
	}

	RenderTemplateGlobal(w, r, "templates/posts.html", data) // Render the post template with the fetched data
}

// CreatePostHandler handles the creation of a new post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	categoryID := r.FormValue("category_id")

	id, err := strconv.Atoi(categoryID)
	if err != nil {
		errorPage(w, r) // Show error page if category ID is invalid
		return
	}

	if checkPostTitle(title) {
		reloadPageWithError(w, r, "Post title already exists") // Reload page with error if title already exists
		return
	}

	if containsAllHtmlTags(title) {
		reloadPageWithError(w, r, "HTML tags are not allowed") // Reload page with error if title contains HTML tags
		return
	}

	if containsHTMLTags(content) {
		reloadPageWithError(w, r, "Limited HTML tags are allowed") // Reload page with error if content contains HTML tags
		return
	}

	postID, err := newPost(id, title, content, getUsernameByUUID(getSessionUUID(r))) // Create a new post
	if err != nil {
		reloadPageWithError(w, r, "Error creating post") // Reload page with error if post creation fails
		return
	}

	err = r.ParseMultipartForm(10 << 20) // Parse the multipart form data
	if err != nil {
		reloadPageWithError(w, r, "Error parsing form") // Reload page with error if form parsing fails
		return
	}

	files := r.MultipartForm.File["image"] // Retrieve the uploaded files
	for _, fileHandler := range files {
		file, err := fileHandler.Open()
		if err != nil {
			reloadPageWithError(w, r, "Error retrieving file") // Reload page with error if file retrieval fails
			return
		}
		defer file.Close()

		if fileHandler.Size > MaxImageSize {
			reloadPageWithError(w, r, "File size too large") // Reload page with error if file size is too large
			return
		}

		if !isValidType(fileHandler.Header.Get("Content-Type")) {
			reloadPageWithError(w, r, "Invalid file type") // Reload page with error if file type is invalid
			return
		}

		dst, err := os.Create(fmt.Sprintf("./img/upload/%d_%s", postID, fileHandler.Filename))
		if err != nil {
			reloadPageWithError(w, r, "Error saving file") // Reload page with error if file saving fails
			return
		}
		defer dst.Close()

		if _, err = io.Copy(dst, file); err != nil {
			reloadPageWithError(w, r, "Error saving file") // Reload page with error if file copy fails
			return
		}

		if err = uploadImage(postID, fmt.Sprintf("%d_%s", postID, fileHandler.Filename)); err != nil {
			reloadPageWithError(w, r, "Error saving image") // Reload page with error if image upload fails
			return
		}
	}

	createLogs("Post created: " + title)                                                                   // Log the post creation
	http.Redirect(w, r, fmt.Sprintf("/categories/post?name=%s&id=%d", title, postID), http.StatusSeeOther) // Redirect to the new post
}

// DeletePostHandler handles the deletion of a post
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")

	id, err := strconv.Atoi(postID)
	if err != nil {
		errorPage(w, r) // Show error page if post ID is invalid
		return
	}

	deletePost(id)               // Delete the post
	createLogs("Post deleted")   // Log the post deletion
	reloadPageWithoutError(w, r) // Reload the page without error
}

// LikePostHandler handles liking a post
func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")

	id, err := strconv.Atoi(postID)
	if err != nil {
		errorPage(w, r) // Show error page if post ID is invalid
		return
	}

	if !isUserLoggedIn(r) {
		logginPage(w, r) // Redirect to login page if user is not logged in
		return
	}

	userID := getIDByUUID(getSessionUUID(r)) // Get user ID from session

	if hasUserDislikedPost(id, userID) {
		removeDislikePost(id, userID) // Remove dislike if user has disliked the post
	}

	if !hasUserLikedPost(id, userID) {
		newLikePost(id, userID) // Add like if user hasn't liked the post
	} else {
		removeLikePost(id, userID) // Remove like if user has already liked the post
	}

	createLogs("Post liked")     // Log the post like
	reloadPageWithoutError(w, r) // Reload the page without error
}

// DislikePostHandler handles disliking a post
func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")

	id, err := strconv.Atoi(postID)
	if err != nil {
		errorPage(w, r) // Show error page if post ID is invalid
		return
	}

	if !isUserLoggedIn(r) {
		logginPage(w, r) // Redirect to login page if user is not logged in
		return
	}

	userID := getIDByUUID(getSessionUUID(r)) // Get user ID from session

	if hasUserLikedPost(id, userID) {
		removeLikePost(id, userID) // Remove like if user has liked the post
	}

	if !hasUserDislikedPost(id, userID) {
		newDislikePost(id, userID) // Add dislike if user hasn't disliked the post
	} else {
		removeDislikePost(id, userID) // Remove dislike if user has already disliked the post
	}

	createLogs("Post disliked")  // Log the post dislike
	reloadPageWithoutError(w, r) // Reload the page without error
}

// isValidType checks if the file type is valid for image uploads
func isValidType(fileType string) bool {
	switch fileType {
	case "image/png", "image/jpg", "image/jpeg", "image/gif", "image/svg+xml", "image/webp":
		return true
	default:
		return false
	}
}

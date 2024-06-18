package logic

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// Posts retrieves the post data by ID
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
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	session := getActiveSession(r)

	data, err := getPostsData(id, session)
	if err != nil {
		http.Error(w, "Error retrieving post data", http.StatusInternalServerError)
		return
	}

	RenderTemplateGlobal(w, r, "templates/posts.html", data)
}

// CreatePostHandler handles the creation of a new post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	categoryID := r.FormValue("category_id")

	id, err := strconv.Atoi(categoryID)
	if err != nil {
		RenderTemplateError(w, r, "templates/categories.html", ErrorMessage{Error: "Invalid category ID"}, nil)
		return
	}

	data, _ := getCategoryData(id)

	if checkPostTitle(title) {
		RenderTemplateError(w, r, "templates/categories.html", ErrorMessage{Error: "Post already exists"}, data)
		return
	}

	postID, err := newPost(id, title, content, getUsernameByUUID(getSessionUUID(r)))
	if err != nil {
		RenderTemplateError(w, r, "templates/categories.html", ErrorMessage{Error: "Error creating post"}, data)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		RenderTemplateError(w, r, "templates/categories.html", ErrorMessage{Error: "Error parsing multipart form"}, data)
		return
	}

	files := r.MultipartForm.File["image"]
	for _, fileHandler := range files {
		file, err := fileHandler.Open()
		if err != nil {
			RenderTemplateError(w, r, "templates/categories.html", ErrorMessage{Error: "Error uploading file"}, data)
			return
		}
		defer file.Close()

		if fileHandler.Size > MaxImageSize {
			RenderTemplateError(w, r, "templates/categories.html", ErrorMessage{Error: "Image is too large. Maximum allowed size is 20 MB."}, data)
			return
		}

		if !isValidType(fileHandler.Header.Get("Content-Type")) {
			RenderTemplateError(w, r, "templates/categories.html", ErrorMessage{Error: "Invalid file type"}, data)
			return
		}

		dst, err := os.Create(fmt.Sprintf("./img/upload/%d_%s", postID, fileHandler.Filename))
		if err != nil {
			RenderTemplateError(w, r, "templates/categories.html", ErrorMessage{Error: "Error creating file"}, data)
			return
		}
		defer dst.Close()

		if _, err = io.Copy(dst, file); err != nil {
			RenderTemplateError(w, r, "templates/categories.html", ErrorMessage{Error: "Error saving file"}, data)
			return
		}

		if err = uploadImage(postID, fmt.Sprintf("%d_%s", postID, fileHandler.Filename)); err != nil {
			RenderTemplateError(w, r, "templates/categories.html", ErrorMessage{Error: "Error uploading image"}, data)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/categories/post?name=%s&id=%d", title, postID), http.StatusSeeOther)
}

// DeletePostHandler handles the deletion of a post
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")

	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	deletePost(id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LikePostHandler handles liking a post
func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")

	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	if !isUserLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID := getIDByUUID(getSessionUUID(r))

	if hasUserDislikedPost(id, userID) {
		removeDislikePost(id, userID)
	}

	if !hasUserLikedPost(id, userID) {
		newLikePost(id, userID)
	} else {
		removeLikePost(id, userID)
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

// DislikePostHandler handles disliking a post
func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")

	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	if !isUserLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID := getIDByUUID(getSessionUUID(r))

	if hasUserLikedPost(id, userID) {
		removeLikePost(id, userID)
	}

	if !hasUserDislikedPost(id, userID) {
		newDislikePost(id, userID)
	} else {
		removeDislikePost(id, userID)
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
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

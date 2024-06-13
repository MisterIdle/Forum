package logic

import (
	"fmt"
	"net/http"
	"strconv"
)

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := fetchCommentsByID(id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := Posts{
		PostID:       id,
		Title:        post.Title,
		Content:      post.Content,
		Timestamp:    post.Timestamp,
		Username:     getUsernameByPostID(id),
		Likes:        getLikesByPostID(id),
		Dislikes:     getDislikesByPostID(id),
		Images:       getImagesByPostID(id),
		Comments:     getCommentsByPostID(id),
		CategoryName: getCategoryNameByPostID(id),
		CategoryID:   getCategoryIDByPostID(id),
	}

	RenderTemplateGlobal(w, r, "templates/posts.html", data)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")
	title := r.FormValue("title")

	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/categories/post?name=%s&id=%d", title, id), http.StatusSeeOther)

	deletePost(id)
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")
	title := r.FormValue("title")

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

	http.Redirect(w, r, fmt.Sprintf("/categories/post?name=%s&id=%d", title, id), http.StatusSeeOther)
}

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")
	title := r.FormValue("title")

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

	http.Redirect(w, r, fmt.Sprintf("/categories/post?name=%s&id=%d", title, id), http.StatusSeeOther)
}

// Comment

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")
	content := r.FormValue("content")
	title := r.FormValue("title")

	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	newComment(id, content, getUsernameByUUID(getSessionUUID(r)))

	http.Redirect(w, r, fmt.Sprintf("/categories/post?name=%s&id=%d", title, id), http.StatusSeeOther)
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID := r.FormValue("comment_id")
	postID := r.FormValue("post_id")
	title := r.FormValue("title")

	id, err := strconv.Atoi(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	deleteComment(id)

	http.Redirect(w, r, fmt.Sprintf("/categories/post?name=%s&id=%s", title, postID), http.StatusSeeOther)
}

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID := r.FormValue("comment_id")
	postID := r.FormValue("post_id")
	title := r.FormValue("title")

	id, err := strconv.Atoi(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	if !isUserLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID := getIDByUUID(getSessionUUID(r))

	if hasUserDislikedComment(id, userID) {
		removeDislikeComment(id, userID)
	}

	if !hasUserLikedComment(id, userID) {
		newLikeComment(id, userID)
	} else {
		removeLikeComment(id, userID)
	}

	http.Redirect(w, r, fmt.Sprintf("/categories/post?name=%s&id=%s", title, postID), http.StatusSeeOther)
}

func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID := r.FormValue("comment_id")
	postID := r.FormValue("post_id")
	title := r.FormValue("title")

	id, err := strconv.Atoi(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	if !isUserLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID := getIDByUUID(getSessionUUID(r))

	if hasUserLikedComment(id, userID) {
		removeLikeComment(id, userID)
	}

	if !hasUserDislikedComment(id, userID) {
		newDislikeComment(id, userID)
	} else {
		removeDislikeComment(id, userID)
	}

	http.Redirect(w, r, fmt.Sprintf("/categories/post?name=%s&id=%s", title, postID), http.StatusSeeOther)
}

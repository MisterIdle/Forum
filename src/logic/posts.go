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

	post, err := fetchPostByID(id)
	if err != nil {
		http.Error(w, "Error retrieving post", http.StatusInternalServerError)
		return
	}

	data := Posts{
		PostID:    id,
		Title:     post.Title,
		Content:   post.Content,
		Timestamp: post.Timestamp,
		Likes:     getLikesByPostID(id),
		Dislikes:  getDislikesByPostID(id),
		Comments:  getCommentsByPostID(id),
	}

	RenderTemplateGlobal(w, r, "templates/posts.html", data)
}

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")
	title := r.FormValue("title")

	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID := getIDByUUID(getSessionUUID(r))

	if hasUserDislikedPost(id, userID) {
		removeDislike(id, userID)
	}

	if !hasUserLikedPost(id, userID) {
		newLike(id, userID)
	} else {
		removeLike(id, userID)
	}

	http.Redirect(w, r, fmt.Sprintf("/categories/post?name=%s&id=%d", title, id), http.StatusSeeOther)
}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")
	title := r.FormValue("title")

	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID := getIDByUUID(getSessionUUID(r))

	if hasUserLikedPost(id, userID) {
		removeLike(id, userID)
	}

	if !hasUserDislikedPost(id, userID) {
		newDislike(id, userID)
	} else {
		removeDislike(id, userID)
	}

	http.Redirect(w, r, fmt.Sprintf("/categories/post?name=%s&id=%d", title, id), http.StatusSeeOther)
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")
	content := r.FormValue("content")
	title := r.FormValue("title")

	id, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	newComment(id, content)

	http.Redirect(w, r, fmt.Sprintf("/categories/post?name=%s&id=%d", title, id), http.StatusSeeOther)
}

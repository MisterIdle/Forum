package logic

import (
	"fmt"
	"net/http"
	"strconv"
)

type Post struct {
	PostID    int
	Title     string
	Content   string
	Timestamp string
}

type Comment struct {
	CommentID int
	Content   string
	Timestamp string
}

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

	data := struct {
		PostID    int
		Title     string
		Content   string
		Timestamp string
		Comments  []Comment
	}{
		PostID:    id,
		Title:     post.Title,
		Content:   post.Content,
		Timestamp: post.Timestamp,
		Comments:  getCommentsByPostID(id),
	}

	RenderTemplateGlobal(w, "templates/posts.html", data)
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

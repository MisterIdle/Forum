package logic

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	posts := getPostsByCategoryID(id)
	categoryName := getCategoryName(id)

	data := Category{
		CategoryID: id,
		Name:       categoryName,
		Posts:      posts,
	}

	RenderTemplateGlobal(w, r, "templates/categories.html", data)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	categoryID := r.FormValue("category_id")

	id, err := strconv.Atoi(categoryID)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	title = strings.Replace(title, " ", "-", -1)

	postID, err := newPost(id, title, content, getUsernameByUUID(getSessionUUID(r)))
	if err != nil {
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/categories/post?name=%s&id=%d", title, postID), http.StatusSeeOther)
}

package logic

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const MaxImageSize = 20 * 1024 * 1024

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	data := Category{
		CategoryID:  id,
		Name:        getCategoryName(id),
		Description: getCategoryDescription(id),
		TotalPosts:  getPostTotalsByCategoryID(id),
		Posts:       getPostsByCategoryID(id),
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

	postID, err := newPost(id, title, content, getUsernameByUUID(getSessionUUID(r)))
	if err != nil {
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	r.ParseMultipartForm(10 << 20)

	files := r.MultipartForm.File["image"]

	for i, fileHandler := range files {
		file, err := fileHandler.Open()
		if err != nil {
			http.Error(w, "Error uploading file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fileSize := fileHandler.Size
		if fileSize > MaxImageSize {
			http.Error(w, "Image is too large. Maximum allowed size is 20 MB.", http.StatusBadRequest)
			return
		}

		if !isValideType(fileHandler.Header.Get("Content-Type")) {
			http.Error(w, "Invalid file type", http.StatusBadRequest)
			return
		}

		dst, _ := os.Create(fmt.Sprintf("./img/upload/%d_%s", postID, files[i].Filename))
		defer dst.Close()

		io.Copy(dst, file)

		err = uploadImage(postID, fmt.Sprintf("%d_%s", postID, files[i].Filename))
		if err != nil {
			http.Error(w, "Error uploading image", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/categories/post?name=%s&id=%d", title, postID), http.StatusSeeOther)
}

func isValideType(fileType string) bool {
	switch fileType {
	case "image/png", "image/jpg", "image/jpeg", "image/gif", "image/svg", "image/webp":
		return true
	}
	return false
}

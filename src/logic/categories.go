package logic

import (
	"net/http"
	"strconv"
)

func getCategoryData(id int) (Category, error) {
	category := Category{
		CategoryID:  id,
		Name:        getCategoryName(id),
		Description: getCategoryDescription(id),
		TotalPosts:  getPostTotalsByCategoryID(id),
		Posts:       getPostsByCategoryID(id),
	}
	return category, nil
}

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		errorPage(w, r)
		return
	}

	data, err := getCategoryData(id)
	if err != nil {
		errorPage(w, r)
		return
	}

	RenderTemplateGlobal(w, r, "templates/categories.html", data)
}

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	description := r.FormValue("description")
	global := r.FormValue("global")

	if checkCategoryName(name) {
		reloadPageWithError(w, r, "Category already exists")
		return
	}

	err := createCategory(name, description, global)
	if err != nil {
		reloadPageWithError(w, r, "Error creating category")
		return
	}

	reloadPageWithoutError(w, r)
}

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryName := r.FormValue("categories")

	err := deleteCategory(categoryName)
	if err != nil {
		return
	}

	reloadPageWithoutError(w, r)
}

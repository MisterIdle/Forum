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
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	data, err := getCategoryData(id)
	if err != nil {
		http.Error(w, "Error retrieving category data", http.StatusInternalServerError)
		return
	}

	RenderTemplateGlobal(w, r, "templates/categories.html", data)
}

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	description := r.FormValue("description")
	global := r.FormValue("global")

	if checkCategoryName(name) {
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)

		data, _ := getCategoryData(id)
		RenderTemplateError(w, r, "templates/categories.html", ErrorMessage{Error: "Category already exists"}, data)
		return
	}

	err := createCategory(name, description, global)
	if err != nil {
		http.Error(w, "Error creating category", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryName := r.FormValue("categories")

	err := deleteCategory(categoryName)
	if err != nil {
		data, _ := getCategoryData(getCategoryIDByName(categoryName))
		RenderTemplateError(w, r, "templates/categories.html", ErrorMessage{Error: "Error deleting category"}, data)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

package logic

import (
	"net/http"
	"strconv"
)

// Fetch category data based on the given ID
func getCategoryData(id int) (Category, error) {
	category := Category{
		CategoryID:  id,
		Name:        getCategoryName(id),           // Get category name by ID
		Description: getCategoryDescription(id),    // Get category description by ID
		TotalPosts:  getPostTotalsByCategoryID(id), // Get total number of posts in the category
		Posts:       getPostsByCategoryID(id),      // Get posts in the category
	}
	return category, nil
}

// Handler for displaying categories
func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id") // Get the category ID from the URL query parameter

	id, err := strconv.Atoi(idStr) // Convert the ID string to an integer
	if err != nil || id <= 0 {
		errorPage(w, r) // If conversion fails or ID is invalid, show error page
		return
	}

	data, err := getCategoryData(id) // Fetch category data
	if err != nil {
		errorPage(w, r) // If fetching data fails, show error page
		return
	}

	RenderTemplateGlobal(w, r, "templates/categories.html", data) // Render the category template with data
}

// Handler for creating a new category
func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	description := r.FormValue("description")
	global := r.FormValue("global")

	// Check if category name already exists
	if checkCategoryName(name) {
		reloadPageWithError(w, r, "Category already exists") // Reload page with error message
		return
	}

	// Create the new category
	err := createCategory(name, description, global)
	if err != nil {
		reloadPageWithError(w, r, "Error creating category") // Reload page with error message if creation fails
		return
	}

	createLogs("Category created: " + name) // Log the category creation
	reloadPageWithoutError(w, r)            // Reload page without error message after successful creation
}

// Handler for deleting a category
func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryName := r.FormValue("categories") // Get the category name from form data

	// Delete the category
	err := deleteCategory(categoryName)
	if err != nil {
		return
	}

	createLogs("Category deleted: " + categoryName) // Log the category deletion
	reloadPageWithoutError(w, r)                    // Reload page without error message after successful deletion
}

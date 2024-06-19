// Code by MisterIdle
// Try Pistol Hand on Itch.io: https://misteridle.itch.io/pistol-hand :D

package logic

import (
	"fmt"
	"net/http"
)

// MaxImageSize defines the maximum size for image uploads (20 MB)
const MaxImageSize = 20 * 1024 * 1024

// LaunchApp initializes the server and starts listening on port 3030
func LaunchApp() {
	HandleAll() // Set up the route handlers
	fmt.Println("")
	fmt.Println("######################")
	fmt.Println("Welcome to the forum!")
	fmt.Println("######################")
	fmt.Println("")

	fmt.Println("Server is running on port 3030 🌐")
	fmt.Println("Visit http://localhost:3030 to access the forum")

	createLogs("Server started") // Log the server start

	err := http.ListenAndServe(":3030", nil) // Start the server
	if err != nil {
		fmt.Println("Error starting server: ", err) // Print error if server fails to start
	}
}

// HandleAll sets up the route handlers for the server
func HandleAll() {
	// Serve static files
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))
	http.Handle("/javascript/", http.StripPrefix("/javascript/", http.FileServer(http.Dir("javascript"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	// Route handlers for various functionalities
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/categories/", CategoriesHandler)
	http.HandleFunc("/categories/post/", PostsHandler)
	http.HandleFunc("/profile/", ProfileHandler)
	http.HandleFunc("/dashboard/", DashboardHandler)
	http.HandleFunc("/auth/", AuthHandler)
	http.HandleFunc("/error", ErrorHandler)
	http.HandleFunc("/back", BackHandler)

	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)

	http.HandleFunc("/create-category", CreateCategoryHandler)
	http.HandleFunc("/delete-category", DeleteCategoryHandler)

	http.HandleFunc("/create-post", CreatePostHandler)
	http.HandleFunc("/delete-post", DeletePostHandler)

	http.HandleFunc("/create-comment", CreateCommentHandler)
	http.HandleFunc("/delete-comment", DeleteCommentHandler)

	http.HandleFunc("/change-username", ChangeProfileUsernameHandler)
	http.HandleFunc("/change-password", ChangeProfilePasswordHandler)
	http.HandleFunc("/change-email", ChangeProfileEmailHandler)
	http.HandleFunc("/change-picture", ChangeProfilePictureHandler)
	http.HandleFunc("/delete-account", DeleteProfileHandler)

	http.HandleFunc("/promote", PromoteUserHandler)
	http.HandleFunc("/demote", DemoteUserHandler)
	http.HandleFunc("/delete", DeleteUserHandler)

	http.HandleFunc("/like-post", LikePostHandler)
	http.HandleFunc("/dislike-post", DislikePostHandler)

	http.HandleFunc("/like-comment", LikeCommentHandler)
	http.HandleFunc("/dislike-comment", DislikeCommentHandler)

	http.HandleFunc("/reload", ReloadHandler)
	http.HandleFunc("/reload-auth", AuthReloadHandler)

	http.HandleFunc("/logout", LogoutHandler)
}

// IndexHandler handles the requests to the homepage
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	global := r.URL.Query().Get("global") // Get the 'global' query parameter
	var globals map[string][]Category
	var err error

	// Fetch global categories based on the 'global' query parameter
	if global == "" || global == "all" {
		globals, err = fetchGlobalCategories()
		if err != nil {
			errorPage(w, r) // Show error page if fetching global categories fails
			return
		}
	} else {
		globals, err = fetchGlobalCategoriesByName(global)
		if err != nil {
			errorPage(w, r) // Show error page if fetching global categories by name fails
			return
		}
	}

	// Prepare data for the index page
	data := Categories{
		Globals:       globals,
		AllCategories: fetchCategoriesName(),
		AllGlobals:    fetchGlobalCategoriesName(),
	}

	RenderTemplateGlobal(w, r, "templates/index.html", data) // Render the index template with the data
}

// ErrorHandler handles the error page requests
func ErrorHandler(w http.ResponseWriter, r *http.Request) { // Log the error page access
	RenderTemplateWithoutData(w, "templates/error.html") // Render the error template
}

// BackHandler handles requests to go back to the previous page and clear errors
func BackHandler(w http.ResponseWriter, r *http.Request) {
	setHasError(r, false, "")                     // Clear error status
	http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect to the homepage
}

// ReloadHandler reloads the current page without any error
func ReloadHandler(w http.ResponseWriter, r *http.Request) {
	reloadPageWithoutError(w, r) // Reload the page without error
}

// AuthReloadHandler reloads the authentication page
func AuthReloadHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplateWithoutData(w, "templates/auth.html") // Render the authentication template
}

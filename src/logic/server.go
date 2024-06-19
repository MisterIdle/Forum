package logic

import (
	"fmt"
	"net/http"
)

const MaxImageSize = 20 * 1024 * 1024

func LaunchApp() {
	HandleAll()

	fmt.Println("Server is running on port 3030")
	err := http.ListenAndServe(":3030", nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}

func HandleAll() {
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))
	http.Handle("/javascript/", http.StripPrefix("/javascript/", http.FileServer(http.Dir("javascript"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

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

	http.HandleFunc("/logout", LogoutHandler)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	global := r.URL.Query().Get("global")
	var globals map[string][]Category
	var err error

	if global == "" || global == "all" {
		globals, err = fetchGlobalCategories()
		if err != nil {
			errorPage(w, r)
			return
		}
	} else {
		globals, err = fetchGlobalCategoriesByName(global)
		if err != nil {
			errorPage(w, r)
			return
		}
	}

	data := Categories{
		Globals:       globals,
		AllCategories: fetchCategoriesName(),
		AllGlobals:    fetchGlobalCategoriesName(),
	}

	RenderTemplateGlobal(w, r, "templates/index.html", data)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplateWithoutData(w, "templates/error.html")
}

func BackHandler(w http.ResponseWriter, r *http.Request) {
	mainPage(w, r)
}

func ReloadHandler(w http.ResponseWriter, r *http.Request) {
	reloadPageWithoutError(w, r)
}

package logic

import (
	"fmt"
	"net/http"
	"time"
)

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

	http.HandleFunc("/", IsAuth(IndexHandler))
	http.HandleFunc("/categories/", IsAuth(CategoriesHandler))
	http.HandleFunc("/categories/post/", IsAuth(PostsHandler))

	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)

	http.HandleFunc("/create-post", CreatePostHandler)
	http.HandleFunc("/create-comment", CreateCommentHandler)

	http.HandleFunc("/like-post", LikeHandler)
	http.HandleFunc("/dislike-post", DislikeHandler)

	http.HandleFunc("/logout", LogoutHandler)
}

func IsAuth(hander http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkCookie(r) {
			sessionToken := time.Now().Format("2006-01-02 15:04:05")

			sessions[sessionToken] = Session{
				LoggedIn: false,
			}

			http.SetCookie(w, &http.Cookie{
				Name:  "session_token",
				Value: sessionToken,
			})
		}

		hander(w, r)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := fetchCategories()
	if err != nil {
		http.Error(w, "Error retrieving categories", http.StatusInternalServerError)
		return
	}

	data := Categories{
		Categories: categories,
	}

	RenderTemplateGlobal(w, r, "templates/index.html", data)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if isUserLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		Register(w, r)
		return
	}
	RenderTemplateWithoutData(w, "templates/register.html")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if isUserLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		Login(w, r)
		return
	}
	RenderTemplateWithoutData(w, "templates/register.html")
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sessionToken := c.Value

	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/register", http.StatusSeeOther)
}

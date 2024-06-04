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

	http.HandleFunc("/", WelcomeHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/forgot", ForgotHandler)
	http.HandleFunc("/reset", ResetHandler)

	http.HandleFunc("/login/github/", githubLoginHandler)
	http.HandleFunc("/login/github/callback", githubCallbackHandler)

	http.HandleFunc("/login/google/", googleLoginHandler)
	http.HandleFunc("/login/google/callback", googleCallbackHandler)

	http.HandleFunc("/logout", LogoutHandler)

	http.HandleFunc("/view/categories", ViewCategoryHandler)
	http.HandleFunc("/view/category", ViewCategoryPostsHandler)

	http.HandleFunc("/create/category", CreateCategoryHandler)
	http.HandleFunc("/create/post", CreatePostHandler)
	http.HandleFunc("/create/comment", CreateCommentHandler)
	http.HandleFunc("/post", ViewPostHandler)


}
/*
func getActualSession(r *http.Request) Session {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return Session{}
		}
		return Session{}
	}

	sessionToken := c.Value
	userSession, ok := sessions[sessionToken]
	if !ok {
		return Session{}
	}

	return userSession
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	userSession := getActualSession(r)

	if userSession.Username == "" {
		userSession.Username = "Guest"
		userSession.IsLoggedIn = false
	}

	RenderTemplateGlobal(w, "templates/index.html", userSession)
}*/

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

func ForgotHandler(w http.ResponseWriter, r *http.Request) {
	if isUserLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		Forgot(w, r)
		return
	}

	RenderTemplateWithoutData(w, "templates/register.html")
}

func ResetHandler(w http.ResponseWriter, r *http.Request) {
	if isUserLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		Reset(w, r)
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

package logic

import (
	"fmt"
	"net/http"
)

const (
	GITHUB_CLIENT_ID     = "Ov23liWELhSqACpxuAnb"
	GITHUB_CLIENT_SECRET = "a6764689efbf7cb3f02e844ad5c18215a1eedc36"
	GOOGLE_CLIENT_ID     = "881937808313-8a95bvir63s8ceku4s9f4jmmf6omd3ij.apps.googleusercontent.com"
	GOOGLE_CLIENT_SECRET = "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"
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

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)

	http.HandleFunc("/login/github/", githubLoginHandler)
	http.HandleFunc("/login/github/callback", githubCallbackHandler)

	http.HandleFunc("/login/google/", googleLoginHandler)
	http.HandleFunc("/login/google/callback", googleCallbackHandler)

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	userSession, exists := sessions[sessionToken]
	if !exists {
		fmt.Println("Session not found")
	}

	RenderTemplateGlobal(w, "templates/index.html", userSession)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		Register(w, r)
		return
	}
	RenderTemplateWithoutData(w, "templates/register.html")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		Login(w, r)
		return
	}
	RenderTemplateWithoutData(w, "templates/register.html")
}

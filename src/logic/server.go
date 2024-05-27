package logic

import (
	"fmt"
	"net/http"
)

type Data struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ErrorMessage struct {
	Message string
}

type Session struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	ID       int    `json:"id"`
}

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
	RenderTemplateForum(w, "templates/forum.html", nil, nil)
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

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
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplateGlobal(w, "templates/index.html", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm")

		data := ErrorMessage{
			Message: "",
		}

		if username == "" || email == "" || password == "" || confirmPassword == "" {
			data.Message = "Please fill in all fields"
			RenderTemplateGlobal(w, "templates/register.html", data)
			return
		}

		if password != confirmPassword {
			data.Message = "Passwords do not match"
			RenderTemplateGlobal(w, "templates/register.html", data)
			return
		}

		if checkUser(username, email) {
			data.Message = "User already exists"
			RenderTemplateGlobal(w, "templates/register.html", data)
			return
		}

		insertUser(username, email, password)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	RenderTemplateWithoutData(w, "templates/register.html")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		identifier := r.FormValue("identifier")
		password := r.FormValue("password")

		data := ErrorMessage{
			Message: "",
		}

		if identifier == "" || password == "" {
			data.Message = "Please fill in all fields"
			RenderTemplateGlobal(w, "templates/register.html", data)
			return
		}

		if !checkUser(identifier, identifier) {
			data.Message = "User does not exist"
			RenderTemplateGlobal(w, "templates/register.html", data)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	RenderTemplateWithoutData(w, "templates/register.html")
}

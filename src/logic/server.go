package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	CLIENT_ID     = "Ov23liWELhSqACpxuAnb"
	CLIENT_SECRET = "a6764689efbf7cb3f02e844ad5c18215a1eedc36"
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
	http.HandleFunc("/loggedin", func(w http.ResponseWriter, r *http.Request) {
		loggedinHandler(w, r.FormValue("data"))
	})
}

// GIT HUB LOGIC

func loggedinHandler(w http.ResponseWriter, githubData string) {
	if githubData == "" {
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}
	w.Header().Set("Content-type", "application/json")

	var prettyJSON bytes.Buffer
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}
}

func getGithubClientID() string {
	return CLIENT_ID
}

func getGithubClientSecret() string {
	return CLIENT_SECRET
}

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	githubClientID := getGithubClientID()

	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		githubClientID,
		"http://localhost:3030/login/github/callback",
	)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	githubAccessToken := getGithubAccessToken(code)

	githubData := getGithubData(githubAccessToken)

	var githubUser struct {
		Login     string `json:"login"`
		ID        int    `json:"id"`
		NodeID    string `json:"node_id"`
		AvatarURL string `json:"avatar_url"`
		Name      string `json:"name"`
		Email     string `json:"email"`
	}

	err := json.Unmarshal([]byte(githubData), &githubUser)
	if err != nil {
		log.Panic("Failed to parse GitHub user data")
	}

	if !checkUser(githubUser.Login, githubUser.Email) {
		insertUser(githubUser.Login, githubUser.Email, "", githubUser.AvatarURL, "GITHUB")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getGithubAccessToken(code string) string {

	clientID := getGithubClientID()
	clientSecret := getGithubClientSecret()

	requestBodyMap := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	// POST request to set URL
	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := io.ReadAll(resp.Body)

	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)

	return ghresp.AccessToken
}

func getGithubData(accessToken string) string {
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}
	respbody, _ := io.ReadAll(resp.Body)

	return string(respbody)
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

		insertUser(username, email, password, "default.jpg", "NORMAL")
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

		if !checkPassword(identifier, password) {
			data.Message = "Incorrect password"
			RenderTemplateGlobal(w, "templates/register.html", data)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	RenderTemplateWithoutData(w, "templates/register.html")
}

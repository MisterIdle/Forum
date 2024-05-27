package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// LOGIN LOGIC

func Register(w http.ResponseWriter, r *http.Request) {
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
}

func Login(w http.ResponseWriter, r *http.Request) {
	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	data := ErrorMessage{
		Message: "",
	}

	if identifier == "" || password == "" {
		data.Message = "Please fill in all fields"
		RenderTemplateGlobal(w, "templates/login.html", data)
		return
	}

	if !checkUser(identifier, identifier) {
		data.Message = "User does not exist"
		RenderTemplateGlobal(w, "templates/login.html", data)
		return
	}

	if !checkPassword(identifier, password) {
		data.Message = "Incorrect password"
		RenderTemplateGlobal(w, "templates/login.html", data)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GITHUB LOGIC

func getGithubClientID() string {
	return GITHUB_CLIENT_ID
}

func getGithubClientSecret() string {
	return GITHUB_CLIENT_SECRET
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

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	githubClientID := getGithubClientID()

	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		githubClientID,
		"http://localhost:3030/login/github/callback",
	)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// GOOGLE LOGIC

func getGoogleClientID() string {
	return GOOGLE_CLIENT_ID
}

func getGoogleClientSecret() string {
	return GOOGLE_CLIENT_SECRET
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	googleClientID := getGoogleClientID()
	redirectURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email&state=pseudo-random",
		googleClientID,
		"http://localhost:3030/login/google/callback",
	)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	googleAccessToken := getGoogleAccessToken(code)

	googleUserData := getGoogleUserData(googleAccessToken)

	fmt.Println(googleUserData)

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	err := json.Unmarshal([]byte(googleUserData), &googleUser)
	if err != nil {
		log.Panic("Failed to parse Google user data")
	}

	if !checkUser(googleUser.Name, googleUser.Email) {
		insertUser(googleUser.Name, googleUser.Email, "", googleUser.Picture, "GOOGLE")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getGoogleAccessToken(code string) string {
	clientID := getGoogleClientID()
	clientSecret := getGoogleClientSecret()

	requestBodyMap := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"redirect_uri":  "http://localhost:3030/login/google/callback",
		"grant_type":    "authorization_code",
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, reqerr := http.NewRequest(
		"POST",
		"https://oauth2.googleapis.com/token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := io.ReadAll(resp.Body)

	type googleAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	var gresp googleAccessTokenResponse
	json.Unmarshal(respbody, &gresp)

	return gresp.AccessToken
}

func getGoogleUserData(accessToken string) string {
	req, reqerr := http.NewRequest(
		"GET",
		"https://www.googleapis.com/oauth2/v2/userinfo",
		nil,
	)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("Bearer %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}
	respbody, _ := io.ReadAll(resp.Body)

	return string(respbody)
}

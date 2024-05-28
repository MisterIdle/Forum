package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var sessions = map[string]Session{}

type ErrorMessage struct {
	Message string
}

type Session struct {
	Username string
	Method   string
	expiry   time.Time
}

type Credentials struct {
	Username string
}

var creds Credentials

////////////////////////
// USER AUTH LOGIC /////
////////////////////////

func Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm")

	data := ErrorMessage{
		Message: "",
	}

	if password != confirmPassword {
		data.Message = "Passwords do not match"
		RenderTemplateGlobal(w, "templates/register.html", data)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic("Failed to generate hashed password")
	}

	if checkUser(username, email) {
		data.Message = "Username or email already exists"
		RenderTemplateGlobal(w, "templates/register.html", data)
		return
	}

	creds.Username = username
	createSession(w, username, "LOCAL")

	newUser(username, email, string(hashedPassword), "Default.jpg", "LOCAL")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	password := r.FormValue("password")

	data := ErrorMessage{
		Message: "",
	}

	hashedPassword := getPassword(user)

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		data.Message = "Username or password is incorrect"
		RenderTemplateGlobal(w, "templates/register.html", data)
		return
	}

	creds.Username = user
	createSession(w, user, "LOCAL")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
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

func isUserLoggedIn(r *http.Request) bool {
	c, err := r.Cookie("session_token")
	if err != nil || err == http.ErrNoCookie {
		return false
	}

	sessionToken := c.Value
	_, exists := sessions[sessionToken]
	return exists
}

////////////////////////
// GITHUB LOGIC ////////
////////////////////////

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
		Name      string `json:"login"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	err := json.Unmarshal([]byte(githubData), &githubUser)
	if err != nil {
		log.Panic("Failed to parse Github user data")
	}

	creds.Username = githubUser.Name

	createSession(w, githubUser.Name, "GITHUB")

	if !checkUser(githubUser.Name, githubUser.Email) {
		newUser(githubUser.Name, githubUser.Email, "", githubUser.AvatarURL, "GITHUB")
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

////////////////////////
// GOOGLE LOGIC ////////
////////////////////////

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

	var googleUser struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	err := json.Unmarshal([]byte(googleUserData), &googleUser)
	if err != nil {
		log.Panic("Failed to parse Google user data")
	}

	creds.Username = googleUser.Name

	createSession(w, googleUser.Name, "GOOGLE")

	if !checkUser(googleUser.Name, googleUser.Email) {
		newUser(googleUser.Name, googleUser.Email, "", googleUser.Picture, "GOOGLE")
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

func createSession(w http.ResponseWriter, username, method string) {
	sessionToken := time.Now().Format("2006-01-02 15:04:05")
	expiresAt := time.Now().Add(120 * time.Second)

	sessions[sessionToken] = Session{
		Username: username,
		Method:   method,
		expiry:   expiresAt,
	}

	cookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}

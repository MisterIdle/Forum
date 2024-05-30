package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var sessions = map[string]Session{}

const (
	GITHUB_CLIENT_ID     = "Ov23liWELhSqACpxuAnb"
	GITHUB_CLIENT_SECRET = "a6764689efbf7cb3f02e844ad5c18215a1eedc36"
	GOOGLE_CLIENT_ID     = "881937808313-8a95bvir63s8ceku4s9f4jmmf6omd3ij.apps.googleusercontent.com"
	GOOGLE_CLIENT_SECRET = "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"
	FORGOT_EMAIL         = "noreplyforumtest@gmail.com"
	FORGOT_PASSWORD      = "lnkqxjttfyrzyoju"
	SMTP_ADDRESS         = "smtp.gmail.com"
	SMTP_PORT            = "587"
)

type ErrorMessage struct {
	Message string
}

type Session struct {
	Username string
	expiry   time.Time
}

type Credentials struct {
	Username string
}

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

	hashedPassword := hashedPassword(password)

	if checkUserMailOrIdentidy(email) {
		data.Message = "Email already exists"
		RenderTemplateGlobal(w, "templates/register.html", data)
		return
	}

	newUser(username, email, string(hashedPassword), "LOCAL", "Default.png")
	createSession(w, username)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	password := r.FormValue("password")

	data := ErrorMessage{Message: ""}

	hashedPassword, username := getCredentialsByUsernameOrEmail(user)

	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		data.Message = "Invalid password or username/email"
		RenderTemplateGlobal(w, "templates/register.html", data)
		return
	}

	createSession(w, username)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Forgot(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	code := generateRandomCode()

	data := ErrorMessage{
		Message: "",
	}

	if !checkUserMailOrIdentidy(email) {
		data.Message = "Email does not exist"
		RenderTemplateGlobal(w, "templates/register.html", data)
		return
	}

	sendEmail("Password Reset", code, []string{email})
	setCodeByEmail(email, code)

	http.Redirect(w, r, "/reset", http.StatusSeeOther)

}

func Reset(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
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

	email := getEmailFromCode(code)

	hashedPassword := hashedPassword(password)

	resetPassword(email, hashedPassword)
	resetCode(email)

	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func hashedPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic("Failed to generate hashed password")
	}
	return string(hashedPassword)
}

func generateRandomCode() string {
	code := fmt.Sprintf("%d", rand.Intn(999999))
	return code
}

func sendEmail(subject string, body string, to []string) {
	auth := smtp.PlainAuth(
		"",
		FORGOT_EMAIL,
		FORGOT_PASSWORD,
		SMTP_ADDRESS,
	)

	msg := "From: " + FORGOT_EMAIL + "\n" + "Subject: " + subject + "\n" + body

	err := smtp.SendMail(
		SMTP_ADDRESS+":"+SMTP_PORT,
		auth,
		FORGOT_EMAIL,
		to,
		[]byte(msg),
	)

	if err != nil {
		log.Panic("Email failed to send")
	}

	log.Print("Email sent")
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
		ID        int    `json:"id"`
		AvatarURL string `json:"avatar_url"`
	}

	err := json.Unmarshal([]byte(githubData), &githubUser)
	if err != nil {
		log.Panicf("Failed to parse Github user data: %v", err) // Include error details
	}

	if checkUserMailOrIdentidy(fmt.Sprint(githubUser.ID)) {
		fmt.Println("Welcome Back: " + githubUser.Name)
	} else {
		newUser(githubUser.Name, "GITHUB", "GITHUB", fmt.Sprint(githubUser.ID), githubUser.AvatarURL)
	}

	createSession(w, githubUser.Name)

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
		Name      string `json:"name"`
		ID        string `json:"id"`
		AvatarURL string `json:"picture"`
	}

	err := json.Unmarshal([]byte(googleUserData), &googleUser)
	if err != nil {
		log.Panic("Failed to parse Google user data")
	}

	if checkUserMailOrIdentidy(googleUser.ID) {
		fmt.Println("Welcome Back: " + googleUser.Name)
	} else {
		newUser(googleUser.Name, "GOOGLE", "GOOGLE", googleUser.ID, googleUser.AvatarURL)
	}

	createSession(w, googleUser.Name)

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

func createSession(w http.ResponseWriter, username string) {
	sessionToken := time.Now().Format("2006-01-02 15:04:05")
	expiresAt := time.Now().Add(120 * time.Second)

	sessions[sessionToken] = Session{
		Username: username,
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

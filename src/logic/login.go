package logic

import (
	"log"
	"net/http"
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

// Basic Auth logic

func Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm")

	data := ErrorMessage{Error: ""}

	if password != confirmPassword {
		data.Error = "Passwords do not match"
		RenderTemplateGlobal(w, r, "templates/register.html", data)
		return
	}

	hashedPassword := hashedPassword(password)

	if checkUserMailOrIdentidy(email) {
		data.Error = "Email already exists"
		RenderTemplateGlobal(w, r, "templates/register.html", data)
		return
	}

	newUser(username, email, string(hashedPassword), "LOCAL", "Default.png")
	createSession(w, username)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	password := r.FormValue("password")

	data := ErrorMessage{Error: ""}

	hashedPassword, username := getCredentialsByUsernameOrEmail(user)

	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		data.Error = "Invalid password or username/email"
		RenderTemplateGlobal(w, r, "templates/register.html", data)
		return
	}

	createSession(w, username)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func hashedPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic("Failed to generate hashed password")
	}
	return string(hashedPassword)
}

// Session logic

func isUserLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false
	}

	session, ok := sessions[cookie.Value]
	if !ok {
		return false
	}

	return session.LoggedIn
}

func checkCookie(r *http.Request) bool {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false
	}

	_, ok := sessions[cookie.Value]
	return ok
}

func createSession(w http.ResponseWriter, username string) {
	sessionToken := time.Now().Format("2006-01-02 15:04:05")

	sessions[sessionToken] = Session{
		LoggedIn: true,
		Username: username,
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
}

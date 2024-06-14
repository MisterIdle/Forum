package logic

import (
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var sessions = map[string]Session{}

// Basic Auth logic

func Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	data := ErrorMessage{Error: ""}

	if isUserLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if password != confirmPassword {
		data.Error = "Passwords do not match"
		RenderTemplateGlobal(w, r, "templates/register.html", data)
		return
	}

	hashedPassword := hashedPassword(password)

	if checkUserEmail(email) {
		data.Error = "Email already exists"
		RenderTemplateGlobal(w, r, "templates/register.html", data)
		return
	}

	if checkUserUsername(username) {
		data.Error = "Username already exists"
		RenderTemplateGlobal(w, r, "templates/register.html", data)
		return
	}

	newUser(username, email, string(hashedPassword), "Default.png", 1)
	createSession(w, username)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	password := r.FormValue("password")

	data := ErrorMessage{Error: ""}

	if isUserLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var hashedPassword, username string
	if strings.Contains(user, "@") {
		hashedPassword, username = getCredentialsByEmail(user)
	} else {
		hashedPassword, username = getCredentialsByUsername(user)
	}

	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		data.Error = "Invalid password"
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

func createSession(w http.ResponseWriter, username string) {
	sessionToken := getUUIDByUsername(username)

	sessions[sessionToken] = Session{
		LoggedIn: true,
		Username: username,
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: sessionToken,
	})
}

func forceLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return
	}

	sessionToken := cookie.Value

	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
}

func getSessionUUID(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return ""
	}

	return cookie.Value
}

package logic

import (
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var sessions = map[string]Session{}

// Authentication handler
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if isUserLoggedIn(r) {
		mainPage(w, r)
		return
	}

	data := getNoSessionData()

	RenderTemplateGlobal(w, r, "templates/auth.html", data)
}

// Registration handler
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	if password != confirmPassword {
		reloadPageWithError(w, r, "Passwords do not match")
		return
	}

	hashedPassword := hashedPassword(password)

	if checkUserEmail(email) {
		reloadPageWithError(w, r, "Email already exists")
		return
	}

	if checkUserUsername(username) {
		reloadPageWithError(w, r, "Username already exists")
		return
	}

	newUser(username, email, string(hashedPassword), "Default.png", 1)
	createSession(w, username)
	mainPage(w, r)

}

// Login handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	password := r.FormValue("password")

	var hashedPassword, username string
	if strings.Contains(user, "@") {
		hashedPassword, username = getCredentialsByEmail(user)
	} else {
		hashedPassword, username = getCredentialsByUsername(user)
	}

	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		reloadPageWithError(w, r, "Invalid username or password")
		return
	}

	createSession(w, username)
	mainPage(w, r)
}

// Hash password
func hashedPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPassword)
}

// Logout handler
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}
		errorPage(w, r)
		return
	}

	sessionToken := c.Value

	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
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
		Rank:     getRankByUUID(sessionToken),
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

func getActiveSession(r *http.Request) Session {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return Session{}
	}

	session, ok := sessions[cookie.Value]
	if !ok {
		return Session{}
	}

	return session
}

func getSessionUUID(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return ""
	}

	return cookie.Value
}

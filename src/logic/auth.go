package logic

import (
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Global map to store sessions
var sessions = map[string]Session{}

// Authentication handler
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if isUserLoggedIn(r) {
		mainPage(w, r) // If user is logged in, redirect to the main page
		return
	}

	data := getNoSessionData(false, "")                     // Prepare data for rendering the authentication page
	RenderTemplateGlobal(w, r, "templates/auth.html", data) // Render the authentication template
}

// Registration handler
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	// Check if passwords match
	if password != confirmPassword {
		data := getNoSessionData(true, "Passwords do not match")
		RenderTemplateGlobal(w, r, "templates/auth.html", data)
		return
	}

	hashedPassword := hashedPassword(password) // Hash the password

	// Check if email already exists
	if checkUserEmail(email) {
		data := getNoSessionData(true, "Email already exists")
		RenderTemplateGlobal(w, r, "templates/auth.html", data)
		return
	}

	// Check if username already exists
	if checkUserUsername(username) {
		data := getNoSessionData(true, "Username already exists")
		RenderTemplateGlobal(w, r, "templates/auth.html", data)
		return
	}

	// Create new user and session
	newUser(username, email, string(hashedPassword), "Default.png", 1)
	createSession(w, username)
	createLogs("User " + username + " registered") // Log the registration
	mainPage(w, r)                                 // Redirect to the main page after successful registration
}

// Login handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	password := r.FormValue("password")

	var hashedPassword, username string
	// Check if user is logging in with email or username
	if strings.Contains(user, "@") {
		hashedPassword, username = getCredentialsByEmail(user)
	} else {
		hashedPassword, username = getCredentialsByUsername(user)
	}

	// Verify the password
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		data := getNoSessionData(true, "Invalid username or password")
		RenderTemplateGlobal(w, r, "templates/auth.html", data)
		return
	}

	createSession(w, username)                    // Create a session after successful login
	createLogs("User " + username + " logged in") // Log the login
	mainPage(w, r)                                // Redirect to the main page
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

	delete(sessions, sessionToken) // Remove the session from the map
	createLogs("User " + getActiveSession(r).Username + " logged out")

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(), // Expire the session token cookie
	})

	http.Redirect(w, r, "/auth/login", http.StatusSeeOther) // Redirect to the login page
}

// Session logic

// Check if user is logged in
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

// Create a new session
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

// Force logout by deleting the session and expiring the cookie
func forceLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return
	}

	sessionToken := cookie.Value

	delete(sessions, sessionToken)
	createLogs("User " + getActiveSession(r).Username + " logged out")

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
}

// Get active session from request
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

// Get session UUID from request
func getSessionUUID(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return ""
	}

	return cookie.Value
}

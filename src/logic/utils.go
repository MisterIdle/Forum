package logic

import (
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// Function to log events
func createLogs(message string) {
	f, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	logMessage := time.Now().Format("2006-01-02 15:04:05") + " - " + message + "\n"
	if _, err := f.WriteString(logMessage); err != nil {
		return
	}
}

// Session-related functions

// setHasError sets the error status and message in the current session
func setHasError(r *http.Request, hasError bool, message string) {
	session := getActiveSession(r)        // Get the active session
	session.HasError = hasError           // Set the error status
	session.Message = message             // Set the error message
	sessions[getSessionUUID(r)] = session // Update the session in the session map
}

// reloadPageWithError reloads the current page with an error message
func reloadPageWithError(w http.ResponseWriter, r *http.Request, message string) {
	setHasError(r, true, message)                                     // Set the error in the session
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther) // Redirect to the referring page
}

// reloadPageWithoutError reloads the current page without an error message
func reloadPageWithoutError(w http.ResponseWriter, r *http.Request) {
	setHasError(r, false, "")                                         // Clear the error in the session
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther) // Redirect to the referring page
}

// mainPage redirects to the main page and clears any errors
func mainPage(w http.ResponseWriter, r *http.Request) {
	setHasError(r, false, "")                     // Clear the error in the session
	http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect to the main page
}

// errorPage redirects to the error page with a generic error message
func errorPage(w http.ResponseWriter, r *http.Request) {
	setHasError(r, true, "Error")                      // Set a generic error in the session
	http.Redirect(w, r, "/error", http.StatusSeeOther) // Redirect to the error page
}

// logginPage redirects to the login page with a prompt to login
func logginPage(w http.ResponseWriter, r *http.Request) {
	setHasError(r, true, "Please login")                    // Set a login prompt in the session
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther) // Redirect to the login page
}

// getNoSessionData creates a Data struct without session-specific data, but with error information
func getNoSessionData(hasError bool, message string) Data {
	data := Data{}
	data.Data = nil
	data.Session = Session{
		HasError: hasError,
		Message:  message,
	}
	return data
}

// HTML Tag-related functions

// allowedTags defines a map of HTML tags that are allowed
var allowedTags = map[string]bool{
	"b": true, "i": true, "u": true, "s": true, "ol": true, "li": true, "ul": true,
}

// containsHTMLTags checks if a string contains any HTML tags that are not allowed
func containsHTMLTags(s string) bool {
	re := regexp.MustCompile(`<[^>]+>`) // Regular expression to find HTML tags
	matches := re.FindAllString(s, -1)
	for _, tag := range matches {
		tagName := strings.Trim(tag, "</>") // Extract tag name
		if _, allowed := allowedTags[tagName]; !allowed {
			return true // Return true if tag is not allowed
		}
	}
	return false // Return false if all tags are allowed
}

// containsAllHtmlTags checks if a string contains any HTML tags
func containsAllHtmlTags(s string) bool {
	re := regexp.MustCompile(`<[^>]+>`) // Regular expression to find HTML tags
	hasTags := re.MatchString(s)
	return hasTags // Return true if any HTML tags are found
}

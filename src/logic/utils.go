package logic

import (
	"net/http"
	"regexp"
	"strings"
)

// Session function

func setHasError(r *http.Request, hasError bool, message string) {
	session := getActiveSession(r)
	session.HasError = hasError
	session.Message = message
	sessions[getSessionUUID(r)] = session
}

func reloadPageWithError(w http.ResponseWriter, r *http.Request, message string) {
	setHasError(r, true, message)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func reloadPageWithoutError(w http.ResponseWriter, r *http.Request) {
	setHasError(r, false, "")
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	setHasError(r, false, "")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func errorPage(w http.ResponseWriter, r *http.Request) {
	setHasError(r, true, "Error")
	http.Redirect(w, r, "/error", http.StatusSeeOther)
}

func logginPage(w http.ResponseWriter, r *http.Request) {
	setHasError(r, true, "Please login")
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

func getNoSessionData(hasError bool, message string) Data {
	data := Data{}
	data.Data = nil
	data.Session = Session{
		HasError: hasError,
		Message:  message,
	}

	return data
}

var allowedTags = map[string]bool{
	"b": true, "i": true, "u": true, "s": true, "ol": true, "li": true, "ul": true,
}

func containsHTMLTags(s string) bool {
	re := regexp.MustCompile(`<[^>]+>`)
	matches := re.FindAllString(s, -1)
	for _, tag := range matches {
		tagName := strings.Trim(tag, "</>")
		if _, allowed := allowedTags[tagName]; !allowed {
			return true
		}
	}
	return false
}

func containsAllHtmlTags(s string) bool {
	re := regexp.MustCompile(`<[^>]+>`)
	return re.MatchString(s)
}

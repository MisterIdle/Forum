package logic

import (
	"net/http"
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

func getNoSessionData() Data {
	data := Data{}
	data.Data = nil
	data.Session = Session{}
	data.Session.HasError = Error{}.HasError
	data.Session.Message = Error{}.Message

	return data
}

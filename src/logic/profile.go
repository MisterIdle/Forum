package logic

import (
	"fmt"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	profile, err := fetchProfile(name)
	if err != nil {
		http.Error(w, "Error retrieving profile", http.StatusInternalServerError)
		return
	}

	data := Profile{
		Username:      profile.Username,
		UUID:          profile.UUID,
		Picture:       profile.Picture,
		Rank:          profile.Rank,
		Timestamp:     profile.Timestamp,
		TotalPosts:    profile.TotalPosts,
		TotalComments: profile.TotalComments,
		TotalLikes:    profile.TotalLikes,
		TotalDislikes: profile.TotalDislikes,
		Posts:         fetchProfilePosts(name),
		Comments:      fetchProfileComments(name),
	}

	RenderTemplateGlobal(w, r, "templates/profile.html", data)
}

func ChangeProfileUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	uuid := r.FormValue("uuid")

	if checkUserUsername(username) {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	changeProfileUsername(username, uuid)

	fmt.Println("Username changed to: ", username)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ChangeProfilePasswordHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	uuid := r.FormValue("uuid")

	changeProfilePassword(hashedPassword(password), uuid)

	fmt.Println("Password changed to: ", password)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ChangeProfileEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	uuid := r.FormValue("uuid")

	if checkUserEmail(email) {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}

	changeProfileEmail(email, uuid)

	fmt.Println("Email changed to: ", email)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

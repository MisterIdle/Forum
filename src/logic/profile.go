package logic

import (
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
		Username:  profile.Username,
		UUID:      profile.UUID,
		Picture:   profile.Picture,
		Rank:      profile.Rank,
		Timestamp: profile.Timestamp,
		Posts:     fetchProfilePosts(name),
		Comments:  fetchProfileComments(name),
	}

	RenderTemplateGlobal(w, r, "templates/profile.html", data)
}

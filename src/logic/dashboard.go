package logic

import "net/http"

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	if !isUserLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// get le view data
	name := r.URL.Query().Get("view")

	profile, err := fetchProfile(name)
	if err != nil {
		RenderTemplateWithoutData(w, "templates/dashboard.html")
		return
	}

	data := Profile{
		Username:      profile.Username,
		UUID:          profile.UUID,
		Picture:       getProfilePictureByUUID(profile.UUID),
		Rank:          profile.Rank,
		Timestamp:     profile.Timestamp,
		TotalPosts:    profile.TotalPosts,
		TotalComments: profile.TotalComments,
		TotalLikes:    profile.TotalLikes,
		TotalDislikes: profile.TotalDislikes,
		Posts:         fetchProfilePosts(name),
		Comments:      fetchProfileComments(name),
	}

	RenderTemplateGlobal(w, r, "templates/dashboard.html", data)
}

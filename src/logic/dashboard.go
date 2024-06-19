package logic

import (
	"net/http"
)

// DashboardHandler handles requests to the user dashboard
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	if !isUserLoggedIn(r) {
		// If not logged in, redirect to the login page
		logginPage(w, r)
		return
	}

	// Get the username or profile name to view from the URL query parameter
	name := r.URL.Query().Get("view")

	// Fetch the profile data for the specified username
	profile, err := fetchProfile(name)
	if err != nil {
		// If an error occurs while fetching the profile, display the error page
		errorPage(w, r)
	}

	// Prepare the data to be rendered in the dashboard template
	data := DashBoard{
		Users:      fetchAllUsernames(),   // Fetch all usernames
		Categories: fetchCategoriesName(), // Fetch all category names
		Profile: Profile{
			Username:      profile.Username,
			UUID:          profile.UUID,
			Picture:       getProfilePictureByUUID(profile.UUID), // Get the profile picture
			Rank:          profile.Rank,
			Timestamp:     profile.Timestamp,
			TotalPosts:    profile.TotalPosts,
			TotalComments: profile.TotalComments,
			TotalLikes:    profile.TotalLikes,
			TotalDislikes: profile.TotalDislikes,
			Posts:         fetchProfilePosts(name),    // Fetch posts made by the user
			Comments:      fetchProfileComments(name), // Fetch comments made by the user
		},
	}

	// Render the dashboard template with the prepared data
	RenderTemplateGlobal(w, r, "templates/dashboard.html", data)
}

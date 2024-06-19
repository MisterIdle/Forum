package logic

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// getProfileData fetches profile data and returns a Profile struct.
func getProfileData(name string) (Profile, error) {
	profile, err := fetchProfile(name)
	if err != nil {
		return Profile{}, err
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

	return data, nil
}

// ProfileHandler handles the profile page request.
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")

	data, err := getProfileData(name)
	if err != nil {
		errorPage(w, r) // Show error page if profile data fetching fails
		return
	}

	createLogs("Profile page accessed: " + name)               // Log the profile page access
	RenderTemplateGlobal(w, r, "templates/profile.html", data) // Render the profile template with fetched data
}

// ChangeProfileUsernameHandler handles the username change request.
func ChangeProfileUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	uuid := r.FormValue("uuid")

	if checkUserUsername(username) {
		reloadPageWithError(w, r, "Username already exists") // Reload page with error if username already exists
		return
	}

	changeProfileUsername(username, uuid) // Change the username
	createLogs("Username changed to: " + username)
	mainPage(w, r) // Redirect to main page
}

// ChangeProfilePasswordHandler handles the password change request.
func ChangeProfilePasswordHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	uuid := r.FormValue("uuid")

	changeProfilePassword(hashedPassword(password), uuid) // Change the password
	createLogs("Password changed")                        // Log the password change
	reloadPageWithoutError(w, r)                          // Reload the page without error
}

// ChangeProfileEmailHandler handles the email change request.
func ChangeProfileEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	uuid := r.FormValue("uuid")

	if checkUserEmail(email) {
		reloadPageWithError(w, r, "Email already exists") // Reload page with error if email already exists
		return
	}

	changeProfileEmail(email, uuid)          // Change the email
	createLogs("Email changed to: " + email) // Log the email change
	reloadPageWithoutError(w, r)             // Reload the page without error
}

// ChangeProfilePictureHandler handles the profile picture change request.
func ChangeProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // Parse the multipart form with a 10 MB limit

	file, handler, err := r.FormFile("picture")
	if err != nil {
		errorPage(w, r) // Show error page if file retrieval fails
		return
	}
	defer file.Close()

	fileSize := handler.Size
	if fileSize > MaxImageSize {
		reloadPageWithError(w, r, "File size too large") // Reload page with error if file size is too large
		return
	}

	if !isValidType(handler.Header.Get("Content-Type")) {
		reloadPageWithError(w, r, "Invalid file type") // Reload page with error if file type is invalid
		return
	}

	uuid := r.FormValue("uuid")
	oldPicture := getProfilePictureByUUID(uuid)
	if oldPicture != "Default.png" {
		err := os.Remove(fmt.Sprintf("./img/profile/%s", oldPicture)) // Remove old profile picture if it's not the default
		if err != nil {
			reloadPageWithError(w, r, "Error deleting old picture") // Reload page with error if old picture deletion fails
			return
		}
	}

	dst, err := os.Create(fmt.Sprintf("./img/profile/%s_%s", uuid, handler.Filename))
	if err != nil {
		reloadPageWithError(w, r, "Error saving file") // Reload page with error if file saving fails
		return
	}
	defer dst.Close()

	io.Copy(dst, file) // Save the new profile picture

	changeProfilePicture(fmt.Sprintf("%s_%s", uuid, handler.Filename), uuid) // Update profile picture in the database
	createLogs("Profile picture changed + UUID: " + uuid)                    // Log the profile picture change
	reloadPageWithoutError(w, r)                                             // Reload the page without error
}

// DeleteProfileHandler handles the profile deletion request.
func DeleteProfileHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")

	deleteProfile(uuid)                           // Delete the profile
	forceLogout(w, r)                             // Force logout the user
	createLogs("Profile deleted + UUID: " + uuid) // Log the profile deletion
	reloadPageWithoutError(w, r)                  // Reload the page without error
}

// PromoteUserHandler handles the user promotion request.
func PromoteUserHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")

	promoteUser(uuid) // Promote the user
	createLogs("User promoted + UUID: " + uuid)
	reloadPageWithoutError(w, r) // Reload the page without error
}

// DemoteUserHandler handles the user demotion request.
func DemoteUserHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")

	demoteUser(uuid) // Demote the user
	createLogs("User demoted + UUID: " + uuid)
	reloadPageWithoutError(w, r) // Reload the page without error
}

// DeleteUserHandler handles the user deletion request.
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")

	deleteProfile(uuid) // Delete the user's profile
	createLogs("User deleted + UUID: " + uuid)
	reloadPageWithoutError(w, r) // Reload the page without error
}

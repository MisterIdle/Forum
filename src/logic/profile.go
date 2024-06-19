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
		errorPage(w, r)
		return
	}

	RenderTemplateGlobal(w, r, "templates/profile.html", data)
}

// ChangeProfileUsernameHandler handles the username change request.
func ChangeProfileUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	uuid := r.FormValue("uuid")

	if checkUserUsername(username) {
		reloadPageWithError(w, r, "Username already exists")
		return
	}

	changeProfileUsername(username, uuid)
	mainPage(w, r)
}

// ChangeProfilePasswordHandler handles the password change request.
func ChangeProfilePasswordHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	uuid := r.FormValue("uuid")

	changeProfilePassword(hashedPassword(password), uuid)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

// ChangeProfileEmailHandler handles the email change request.
func ChangeProfileEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	uuid := r.FormValue("uuid")

	if checkUserEmail(email) {
		reloadPageWithError(w, r, "Email already exists")
		return
	}

	changeProfileEmail(email, uuid)
	reloadPageWithoutError(w, r)
}

// ChangeProfilePictureHandler handles the profile picture change request.
func ChangeProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("picture")
	if err != nil {
		errorPage(w, r)
		return
	}
	defer file.Close()

	fileSize := handler.Size
	if fileSize > MaxImageSize {
		reloadPageWithError(w, r, "File size too large")
		return
	}

	if !isValidType(handler.Header.Get("Content-Type")) {
		reloadPageWithError(w, r, "Invalid file type")
		return
	}

	uuid := r.FormValue("uuid")
	oldPicture := getProfilePictureByUUID(uuid)
	if oldPicture != "Default.png" {
		err := os.Remove(fmt.Sprintf("./img/profile/%s", oldPicture))
		if err != nil {
			reloadPageWithError(w, r, "Error deleting old picture")
			return
		}
	}

	dst, err := os.Create(fmt.Sprintf("./img/profile/%s_%s", uuid, handler.Filename))
	if err != nil {
		reloadPageWithError(w, r, "Error saving file")
		return
	}
	defer dst.Close()

	io.Copy(dst, file)

	changeProfilePicture(fmt.Sprintf("%s_%s", uuid, handler.Filename), uuid)
	reloadPageWithoutError(w, r)
}

// DeleteProfileHandler handles the profile deletion request.
func DeleteProfileHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")

	deleteProfile(uuid)
	forceLogout(w, r)
	reloadPageWithoutError(w, r)
}

// PromoteUserHandler handles the user promotion request.
func PromoteUserHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")

	promoteUser(uuid)
	reloadPageWithoutError(w, r)
}

// DemoteUserHandler handles the user demotion request.
func DemoteUserHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")

	demoteUser(uuid)
	reloadPageWithoutError(w, r)
}

// DeleteUserHandler handles the user deletion request.
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")

	deleteProfile(uuid)
	reloadPageWithoutError(w, r)
}

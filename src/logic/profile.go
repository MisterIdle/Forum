package logic

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

	http.Redirect(w, r, fmt.Sprintf("/profile/?name=%s", username), http.StatusSeeOther)
}

func ChangeProfilePasswordHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	uuid := r.FormValue("uuid")

	changeProfilePassword(hashedPassword(password), uuid)

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func ChangeProfileEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	uuid := r.FormValue("uuid")

	if checkUserEmail(email) {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}

	changeProfileEmail(email, uuid)

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func ChangeProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("picture")
	if err != nil {
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	fileSize := handler.Size
	if fileSize > MaxImageSize {
		http.Error(w, "Image is too large. Maximum allowed size is 20 MB.", http.StatusBadRequest)
		return
	}

	if !isValidType(handler.Header.Get("Content-Type")) {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	uuid := r.FormValue("uuid")

	oldPicture := getProfilePictureByUUID(uuid)
	if oldPicture != "Default.png" {
		err := os.Remove(fmt.Sprintf("./img/profile/%s", oldPicture))
		if err != nil {
			http.Error(w, "Error removing old picture", http.StatusInternalServerError)
			return
		}
	}

	dst, _ := os.Create(fmt.Sprintf("./img/profile/%s_%s", uuid, handler.Filename))
	defer dst.Close()

	io.Copy(dst, file)

	changeProfilePicture(fmt.Sprintf("%s_%s", uuid, handler.Filename), uuid)

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func DeleteProfileHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")

	deleteProfile(uuid)

	forceLogout(w, r)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func PromoteUserHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")

	promoteUser(uuid)

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func DemoteUserHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")

	demoteUser(uuid)

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")

	deleteProfile(uuid)

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

// Code by MisterIdle
// Try Pistol Hand on Itch.io: https://misteridle.itch.io/pistol-hand :D

package logic

import (
	"fmt"
	"net/http"
	"text/template"
)

// RenderTemplateGlobal renders a template with global data including session information.
func RenderTemplateGlobal(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	// Parse the specified template file
	tmpt, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Print("Error parsing template: ", err)
		return
	}

	// Prepare data with session information to pass to the template
	dataWithSession := Data{
		Data: data,
		Session: Session{
			Username: getUsernameByUUID(getSessionUUID(r)), // Get the username from the session UUID
			LoggedIn: isUserLoggedIn(r),                    // Check if the user is logged in
			Rank:     getRankByUUID(getSessionUUID(r)),     // Get the user's rank from the session UUID
			HasError: getActiveSession(r).HasError,         // Get the error status from the active session
			Message:  getActiveSession(r).Message,          // Get the message from the active session
		},
	}

	// Execute the template with the data including session information
	err = tmpt.Execute(w, dataWithSession)
	if err != nil {
		fmt.Print("Error executing template: ", err)
		return
	}
}

// RenderTemplateWithoutData renders a template without any additional data.
func RenderTemplateWithoutData(w http.ResponseWriter, tmpl string) {
	// Parse the specified template file
	tmpt, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Print("Error parsing template: ", err)
		return
	}

	// Execute the template without any data
	err = tmpt.Execute(w, nil)
	if err != nil {
		fmt.Print("Error executing template: ", err)
		return
	}
}

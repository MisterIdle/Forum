package logic

import (
	"fmt"
	"net/http"
	"text/template"
)

// RenderTemplateGlobal renders a template with global data.
func RenderTemplateGlobal(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	tmpt, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Print("Error parsing template: ", err)
		return
	}

	dataWithSession := Data{
		Data: data,
		Session: Session{
			Username: getUsernameByUUID(getSessionUUID(r)),
			LoggedIn: isUserLoggedIn(r),
			Rank:     getRankByUUID(getSessionUUID(r)),
			HasError: getActiveSession(r).HasError,
			Message:  getActiveSession(r).Message,
		},
	}

	err = tmpt.Execute(w, dataWithSession)
	if err != nil {
		fmt.Print("Error executing template: ", err)
		return
	}
}

// RenderTemplateWithoutData renders a template without any additional data.
func RenderTemplateWithoutData(w http.ResponseWriter, tmpl string) {
	tmpt, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Print("Error parsing template: ", err)
		return
	}

	err = tmpt.Execute(w, nil)
	if err != nil {
		fmt.Print("Error executing template: ", err)
		return
	}
}

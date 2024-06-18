package logic

import (
	"fmt"
	"net/http"
	"text/template"
)

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
		},
		ErrorMessage: ErrorMessage{
			Error: "",
		},
	}

	err = tmpt.Execute(w, dataWithSession)
	if err != nil {
		fmt.Print("Error executing template: ", err)
		return
	}
}

func RenderTemplateError(w http.ResponseWriter, r *http.Request, templateName string, errMsg ErrorMessage, data interface{}) {
	tmpl, err := template.ParseFiles(templateName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	combinedData := struct {
		Error ErrorMessage
		Data  interface{}
	}{
		Error: errMsg,
		Data:  data,
	}

	tmpl.Execute(w, combinedData)
}

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

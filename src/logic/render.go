package logic

import (
	"fmt"
	"net/http"
	"text/template"
)

func RenderTemplateForum(w http.ResponseWriter, tmpl string, data interface{}, session *Session) {
	tmpt, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Print("Error parsing template: ", err)
		return
	}

	dataMap := map[string]interface{}{}

	err = tmpt.Execute(w, dataMap)
	if err != nil {
		fmt.Print("Error executing template: ", err)
		return
	}
}

func RenderTemplateGlobal(w http.ResponseWriter, tmpl string, data interface{}) {
	tmpt, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Print("Error parsing template: ", err)
		return
	}

	err = tmpt.Execute(w, data)
	if err != nil {
		fmt.Print("Error executing template: ", err)
		return
	}
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

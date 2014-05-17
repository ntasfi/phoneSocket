package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type DesktopNoPostTemplate struct {
	FormAction string
}

type DesktopTemplate struct {
	LobbyID       string
	ServerAddress string
	MobileAddress string
}

type MobileNoCodeTemplate struct {
	FormAction string
}

type MobileTemplate struct {
	LobbyID   string
	ServerURL string
	Frequency int
}

//reads the template file and returns a string of it.
func loadTemplate(templateName string) string {
	templateFileLocation := fmt.Sprintf("%s%s.%s", configuration.Template.RootLocation, templateName, configuration.Template.Extension)
	contents, err := ioutil.ReadFile(templateFileLocation)
	if err != nil {
		log.Fatal(err)
	}
	return string(contents)
}

func createExecuteTemplate(w http.ResponseWriter, templateIdentifier string, data interface{}) error {
	//need to update this template to have javascript etc. that connects and listens to the right channels
	var templateString = loadTemplate(templateIdentifier)
	templ := template.Must(template.New(templateIdentifier).Parse(templateString)) //render original page

	return templ.ExecuteTemplate(w, templateIdentifier, data)

	// if data != nil {
	// 	return templ.ExecuteTemplate(w, templateIdentifier, data)
	// } else {
	// 	return templ.Execute(w, nil)
	// }

}

package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	//"net/url"
)

//reads the template file and returns a string of it.
func loadTemplate(templateName string) string {
	contents, err := ioutil.ReadFile("templates/" + templateName + ".tmpl")
	if err != nil {
		log.Fatal(err)
	}
	return string(contents)
}

func handleMobile(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	fmt.Println(req)
	if code == "" {
		fmt.Println("No code given...")
		var templateString = loadTemplate("mobile-noCode")
		templ := template.Must(template.New("mobile-noCode").Parse(templateString))

		templ.Execute(w, req.FormValue("lobbyCode"))
	} else {
		var templateString = loadTemplate("mobile")
		templ := template.Must(template.New("mobile").Parse(templateString))

		templ.Execute(w, nil)
		fmt.Println("Code given...")
		fmt.Println(code)
	}

	fmt.Println("Hello from mobile (=")
}

func handleDesktop(w http.ResponseWriter, req *http.Request) {
	var templateString = loadTemplate("desktop")
	templ := template.Must(template.New("desktop").Parse(templateString))

	templ.Execute(w, req.FormValue("s"))
	fmt.Println("Hello from desktop =)")
}

func handleSocket(ws *websocket.Conn) {
	fmt.Println("Hello from ProcessSocket")
}

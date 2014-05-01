package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
)

func main() {

	http.Handle("/desktop", http.HandlerFunc(handleDesktop))

	http.Handle("/mobile", http.HandlerFunc(handleMobile))

	http.Handle("/websocket/", websocket.Handler(handleSocket))
	//if desktop send here
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
)

func main() {
	//channels!
	var newLobbyChan = make(chan *lobby)         //all the new lobbies are sent over this to our hub to be registered
	var socketConnChan = make(chan socketEntity) //the sockets send a channel to the hub with ID for their connection.

	go hub(newLobbyChan, socketConnChan) //spawn hub to keep track of the lobbies and answer queries about lobbies

	http.Handle("/desktop", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			handleDesktop(newLobbyChan, w, r)
		}))

	http.Handle("/mobile", http.HandlerFunc(handleMobile))

	http.Handle("/javascript/", http.StripPrefix("/javascript/", http.FileServer(http.Dir("../../client/"))))

	http.Handle("/websocket", websocket.Handler(
		func(ws *websocket.Conn) {
			handleSocket(socketConnChan, ws)
		}))

	//if desktop send here
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

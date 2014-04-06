package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http"
)

func webHandler(ws *websocket.Conn) {
	var in []byte
	if err := websocket.Message.Receive(ws, &in); err != nil {
		break
	}
	fmt.Printf("Received: %s\n", string(in))
	websocket.Message.Send(ws, in)
}

func main() {
	fmt.Println("Starting websock server: ")
	http.Handle("/echo", websocket.Handler(webHandler))
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

var myAddress = findIpAddress() //find a better way to find this...maybe create a function that finds out external ip?
var configuration = readConfiguration("configuration.json")

func main() {
	fmt.Println(myAddress)
	var serverURI = fmt.Sprintf("%s:%d", configuration.ServerAddress, configuration.ServerPort)
	var isMobile = regexp.MustCompile(`(M|m)obile|(I|i)P(hone|od|ad)|(A|a)ndroid|(B|b)lackBerry|(I|i)EMobile|(K|k)indle|(N|n)etFront|(S|s)ilk-Accelerated|(hpw|web)OS|(F|f)ennec|(M|m)inimo|(O|o)pera (M|m)(obi|ini)|(B|b)lazer|(D|d)olfin|(D|d)olphin|(S|s)kyfire|(Z|z)une`)

	//communcation settings
	var newLobbyChan = make(chan *lobby)         //all the new lobbies are sent over this to our hub to be registered
	var socketConnChan = make(chan socketEntity) //the sockets send a channel to the hub with ID for their connection.
	var killHubChan = make(chan entity)          //used to kill lobbies etc off. This is sockets to lobby ONLY.

	go hub(newLobbyChan, socketConnChan, killHubChan) //spawn hub to keep track of the lobbies and answer queries about lobbies

	http.Handle(configuration.HTTPRoutes.Root, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if isMobile.MatchString(r.UserAgent()) { //if its a mobile user
				handleMobile(w, r)
			} else {
				handleDesktop(newLobbyChan, w, r)
			}

		}))

	http.Handle(configuration.HTTPRoutes.Javascript.Route, http.StripPrefix(configuration.HTTPRoutes.Javascript.Route, http.FileServer(http.Dir(configuration.HTTPRoutes.Javascript.RootLocation))))

	http.Handle(configuration.HTTPRoutes.Images.Route, http.StripPrefix(configuration.HTTPRoutes.Images.Route, http.FileServer(http.Dir(configuration.HTTPRoutes.Images.RootLocation))))

	http.Handle(configuration.HTTPRoutes.Websocket, websocket.Handler(
		func(ws *websocket.Conn) {
			handleSocket(socketConnChan, killHubChan, ws)
		}))

	fmt.Println("Binding and listening on", serverURI)
	err := http.ListenAndServe(serverURI, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

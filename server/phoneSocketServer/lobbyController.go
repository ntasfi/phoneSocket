package main

import (
	"fmt"
)

func lobbyController(thisLobby *lobby, deviceInChan chan deviceUpdate, desktopOutChan chan deviceUpdate, killChan chan bool) {
	//TODO: add filters to run on incoming data...
	fmt.Println("lobbyController: New lobby running. ", thisLobby)
	var isDead bool

	for {
		select {
		case recievedData := <-deviceInChan:
			//doo something with recievedData here maybe smoothing etc...
			desktopOutChan <- recievedData //send it out
		case isDead = <-killChan:
		}
		if isDead { //is there a better way?
			//clean up in here...close the out channel only.
			break
		}
	} //end for
	fmt.Println("lobbyController: Lobby", thisLobby.id, "has closed.")
} //end lobbyController

package main

import (
	"fmt"
)

func lobbyController(thisLobby *lobby, deviceInChan chan deviceUpdate, desktopOutChan chan deviceUpdate, killChan chan bool) {
	//TODO: add filters to run on incoming data...
	fmt.Println("lobbyController: New lobby running. ", thisLobby)

	for {
		select {
		case recievedData := <-deviceInChan:
			//doo something with recievedData here maybe smoothing etc...
			desktopOutChan <- recievedData //send it out
		case isDead := <-killChan:
			if isDead {
				fmt.Println("lobbyController: Killing lobby", thisLobby.id)
				return
			}
		}
	} //end for
} //end lobbyController

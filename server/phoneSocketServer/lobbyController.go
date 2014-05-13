package main

import (
	"fmt"
)

func lobbyController(thisLobby *lobby, deviceInChan chan deviceUpdate, desktopOutChan chan deviceUpdate) {
	//TODO: add filters to run on incoming data...
	fmt.Println("lobbyController: New lobby running. ", thisLobby)

	go func() {
		for {
			temp := <-desktopOutChan
			fmt.Println("SIMULATED SEND:", temp)
		}
	}()

	for {
		select {
		case recievedData := <-deviceInChan:
			//doo something with recievedData here maybe smoothing etc...
			desktopOutChan <- recievedData //send it out
		}
	}
	fmt.Println("lobbyID", thisLobby.id, "closed.")
}

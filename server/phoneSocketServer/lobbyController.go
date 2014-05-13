package main

import (
	"fmt"
)

func lobbyController(thisLobby *lobby, deviceInChan chan deviceUpdate, desktopOutChan chan deviceUpdate) {
	//TODO: add filters to run on incoming data...
	fmt.Println("lobbyController: New lobby running. ", thisLobby)

	go func() {
		for {
			<-desktopOutChan
			//fmt.Println("SIMULATED SEND:", <-desktopOutChan)
		}
	}()

	for {
		select {
		case recievedData := <-deviceInChan:
			fmt.Println("lobbyController: Recieved data - ", recievedData)
			//doo something with recievedData here maybe smoothing etc...
			desktopOutChan <- recievedData //send it out
		}
	}
	fmt.Println("lobbyID", thisLobby.id, "closed.")
}

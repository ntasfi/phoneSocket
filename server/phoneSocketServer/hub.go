package main

import (
	"fmt"
)

func hub(newLobbyChan chan *lobby, socketConnChan chan socketEntity) {
	fmt.Println("Hub: Started.")
	currentLobbies := make(map[string]*lobby)

	for {
		select {
		case newLobby := <-newLobbyChan: //a new lobby was created
			fmt.Println("Hub: New lobby.")
			currentLobbies[newLobby.id] = newLobby //add it in!
		case newSocketConn := <-socketConnChan: //a new socket has connected
			if _, ok := currentLobbies[newSocketConn.LobbyID]; ok { //make sure the lobby exists
				if newSocketConn.IsDesktop {
					fmt.Println("Hub: New desktop socket connection.")
					newSocketConn.socketChan <- currentLobbies[newSocketConn.LobbyID].desktopOutChan //send the desktop chan back
					fmt.Print("Hub: Sent channel back to handleSocket.")
				} else if newSocketConn.IsDesktop == false { //is device
					fmt.Println("Hub: New device socket connection.")
					//check if the lobby is full
					if currentLobbies[newSocketConn.LobbyID].currentNumberUsers+1 <= currentLobbies[newSocketConn.LobbyID].numberAllowedUsers {
						newSocketConn.socketChan <- currentLobbies[newSocketConn.LobbyID].deviceInputChan //send the device chan back
						newSocketConn.lobbySettings <- currentLobbies[newSocketConn.LobbyID].settings     //send the settings

						currentLobbies[newSocketConn.LobbyID].currentNumberUsers += 1 //increment the lobbies number of devices.
						fmt.Println("Hub: There are", currentLobbies[newSocketConn.LobbyID].currentNumberUsers, "users connected in lobby", newSocketConn.LobbyID)
					} else {
						newSocketConn.errorChan <- "Error: The lobby is full."
					}
				} else {
					//no idea how you got here
				}
			} else { //lobby doesnt exist
				newSocketConn.errorChan <- "Error: The supplied lobbyID was not found."
			}
			//DO A CLIENT CLOSED CONNECTION. use entity as channel type
		} //end select
	} //end for
} //end hub

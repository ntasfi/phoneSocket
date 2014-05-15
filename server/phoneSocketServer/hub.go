package main

import (
	"fmt"
)

func hub(newLobbyChan chan *lobby, socketConnChan chan socketEntity, killHubChan chan entity) {
	fmt.Println("Hub: Started.")
	currentLobbies := make(map[string]*lobby)

	for {
		select {
		case newLobby := <-newLobbyChan: //a new lobby was created
			currentLobbies[newLobby.id] = newLobby //add it in!
			fmt.Println("Hub: New lobby added. There are now", len(currentLobbies), "lobbies.")
		case newSocketConn := <-socketConnChan: //a new socket has connected
			if _, ok := currentLobbies[newSocketConn.LobbyID]; ok { //make sure the lobby exists
				if newSocketConn.IsDesktop {
					fmt.Println("Hub: New desktop socket connection.")
					newSocketConn.socketChan <- currentLobbies[newSocketConn.LobbyID].desktopOutChan //send the desktop chan back
					fmt.Println("Hub: Sent channel back to handleSocket.")
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
		case killThis := <-killHubChan:
			if killThis.IsDesktop { //if its the desktop that died...
				currentLobbies[killThis.LobbyID].killChan <- true //we are killing the lobby
				//close the killChan for this lobby
				delete(currentLobbies, killThis.LobbyID) //remove the lobby from memory
				fmt.Println("Hub: Lobby removed. There are now", len(currentLobbies), "lobbies.")
			} else if _, ok := currentLobbies[killThis.LobbyID]; killThis.IsDesktop == false && ok { //a client has been killed
				currentLobbies[killThis.LobbyID].currentNumberUsers -= 1
				fmt.Println("Hub: There are", currentLobbies[killThis.LobbyID].currentNumberUsers, "users connected in lobby", killThis.LobbyID)
			} else {
				//no idea how you got here
			}
		} //end select
	} //end for
} //end hub

package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func handleMobile(w http.ResponseWriter, req *http.Request) {
	fmt.Println("handleMobile: Request recieved.")
	code := strings.TrimPrefix(req.RequestURI, "/") //remove the '/' from the start. //TODO: sanitize
	//TODO: check if code is valid against current lobbys
	if code == "" { //TODO: better what to check this...maybe check if not valid.
		fmt.Println("handleMobile: No code given.")
		templateString := loadTemplate("mobile-no-code")
		templ := template.Must(template.New("mobile-no-code").Parse(templateString))

		templ.Execute(w, req.FormValue("code"))

		if err := createExecuteTemplate(w, "mobile-no-code", MobileNoCodeTemplate{FormAction: configuration.HTTPRoutes.Root}); err != nil {
			fmt.Println(err)
		}

	} else { //TODO: insert the code into the template.

		serverURL := fmt.Sprintf("%s:%d%s", myAddress, configuration.ServerPort, configuration.HTTPRoutes.Websocket)
		mobileTemplate := MobileTemplate{LobbyID: code, ServerURL: serverURL, Frequency: configuration.Mobile.DefaultUpdateFrequencyMilliseconds}

		if err := createExecuteTemplate(w, "mobile", mobileTemplate); err != nil {
			fmt.Println(err)
		}

		fmt.Println("handleMobile: Code -", code)
	}
}

func handleDesktop(newLobbyChan chan *lobby, w http.ResponseWriter, req *http.Request) {
	fmt.Println("handleDesktop: Request recieved.")
	if req.Method == "POST" { //if post
		req.ParseForm()      //parse it beforehand
		settings := req.Form //get the map

		newLobby := &lobby{}                      //make a new lobby object
		newLobby.settings = make(map[string]bool) //initalize
		//newLobby.id = strconv.FormatInt(time.Now().Unix(), 10)                           //take the unix current time and use it as the lobby id
		newLobby.id = "dator"
		newLobby.numberAllowedUsers, _ = strconv.Atoi(settings["numberAllowedUsers"][0]) //TODO: make sure its an int, validation!

		delete(settings, "numberAllowedUsers") //remove it

		for k, _ := range settings {
			newLobby.settings[k] = true
		}

		newLobby.desktopOutChan = make(chan deviceUpdate)
		newLobby.deviceInputChan = make(chan deviceUpdate)
		newLobby.killChan = make(chan bool)

		newLobbyChan <- newLobby                                                                           //register the new lobby with the hub.
		go lobbyController(newLobby, newLobby.deviceInputChan, newLobby.desktopOutChan, newLobby.killChan) //start the new go process with the channels passed.

		serverURL := fmt.Sprintf("%s:%d%s", myAddress, configuration.ServerPort, configuration.HTTPRoutes.Websocket)
		mobileURL := fmt.Sprintf("%s:%d/%s", myAddress, configuration.ServerPort, newLobby.id)
		if err := createExecuteTemplate(w, "desktop", DesktopTemplate{LobbyID: newLobby.id, ServerAddress: serverURL, MobileAddress: mobileURL}); err != nil {
			fmt.Print(err)
		}
	} else { //else no post
		if err := createExecuteTemplate(w, "desktop-no-post", DesktopNoPostTemplate{FormAction: configuration.HTTPRoutes.Root}); err != nil {
			fmt.Println(err)
		}
	}
}

func handleSocket(socketConnChan chan socketEntity, killHubChan chan entity, ws *websocket.Conn) {
	var message string
	var initialRecievedEntity entity
	var handlerChan (chan deviceUpdate)

	//get the first message
	if err := websocket.Message.Receive(ws, &message); err != nil { //receive the message here
		fmt.Println("handleSocket:", err)
		return
	}

	//its the first connect

	fmt.Println("handleSocket: Message -", message)

	//TODO: make this less sloppy.
	if match, _ := regexp.MatchString("IsDesktop", message); match == true {
		fmt.Println("handleSocket: Handling the first message from client.")
		err := json.Unmarshal([]byte(message), &initialRecievedEntity)
		if err != nil {
			fmt.Println(err)
		}

		tempErrorChan := make(chan string)
		tempLobbySettingsChan := make(chan map[string]bool)
		tempSocketChan := make(chan (chan deviceUpdate)) //this is what we will get our chan deviceUpdate on.

		//send our channel to the hub which will send us the channel to the lobby.
		socketConnChan <- socketEntity{entity: initialRecievedEntity, errorChan: tempErrorChan, socketChan: tempSocketChan, lobbySettings: tempLobbySettingsChan}

		select { //TODO: check if this just runs once and quits out. UPDATE: appears that it just runs once.
		case tempChan := <-tempSocketChan: //recieve the correct channel
			handlerChan = tempChan //set it
			fmt.Println("handleSocket: Recieved the handlerChan from the hub.")

			//update the client here if not desktop
			if initialRecievedEntity.IsDesktop == false { //its a client
				clientSettings := <-tempLobbySettingsChan //this is where we would get the settings for our lobby and send it to the client.
				if err := websocket.JSON.Send(ws, deviceLobbySettings{Settings: clientSettings}); err != nil {
					fmt.Println("Error sending settings. Error:", err)
				}
			}

		case errMsgChan := <-tempErrorChan: //an error occured
			deviceErrMsg := deviceErrorMessage{Error: errMsgChan, Timestamp: time.Now().Unix()}

			if err := websocket.JSON.Send(ws, deviceErrMsg); err != nil {
				fmt.Println("handleSocket", err)
			}

			//close the connection
			if err := ws.Close(); err != nil {
				fmt.Println(err)
			}
		} //end select

	} //end check

	//we test here to see if the websocket was closed.
	if err := websocket.Message.Send(ws, []byte("{ping:true}")); err != nil {
		fmt.Println("handleSocket: Websocket is closed.")
		killHubChan <- initialRecievedEntity //kill this
		return                               //get the hell out of dodge
	}

	if initialRecievedEntity.IsDesktop {
		timer := time.After(configuration.Lobby.InactivityNoActivityMinutes * time.Minute) //Lobbies are destroyed after 5 minutes of inactivity initially.
		var isDead = false
		for {
			select {
			case clientData := <-handlerChan:
				if err := websocket.JSON.Send(ws, clientData); err != nil {
					fmt.Println("Error occured:", err)
					killHubChan <- initialRecievedEntity //kill this
					isDead = true                        //not sure if needed
				}
			case <-timer:
				killHubChan <- initialRecievedEntity
				isDead = true //not sure if needed
			} //end select
			if isDead {
				if err := ws.Close(); err != nil {
					fmt.Print(err)
				}
				fmt.Println("handlers: Desktop socket killed.")
				break
			}
			timer = time.After(configuration.Lobby.InactivityAfterActivitySeconds * time.Second) //after we get a single update this is our new timer.
		} //end for
	} else if initialRecievedEntity.IsDesktop == false { //its a device
		for {
			var recievedUpdate deviceUpdate //setup holder vars
			rawJSONObjectMap := make(map[string]*json.RawMessage)

			recievedUpdate.Measurments = make(map[string]float64)
			recievedUpdate.DeviceID = initialRecievedEntity.DeviceID
			recievedUpdate.Timestamp = time.Now().Unix()

			if err := websocket.Message.Receive(ws, &message); err != nil { //receive the message here from the device
				fmt.Println("Error occured:", err)
				killHubChan <- initialRecievedEntity //kill this
				return
			}

			err := json.Unmarshal([]byte(message), &rawJSONObjectMap) //unmarshal the message into the recievedUpdate var
			if err != nil {
				fmt.Println(err)
			}

			for key, _ := range rawJSONObjectMap {
				var value float64
				err := json.Unmarshal(*rawJSONObjectMap[key], &value)
				if err != nil {
					fmt.Println(err)
				}
				recievedUpdate.Measurments[key] = value //set it
			}

			handlerChan <- recievedUpdate //send the update
		}
	} else {
		fmt.Println("no idea how you got here")
	}

} //end handleSocket

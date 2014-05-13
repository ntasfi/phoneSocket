package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

//reads the template file and returns a string of it.
func loadTemplate(templateName string) string {
	contents, err := ioutil.ReadFile("templates/" + templateName + ".tmpl")
	if err != nil {
		log.Fatal(err)
	}
	return string(contents)
}

func handleMobile(w http.ResponseWriter, req *http.Request) {
	fmt.Println("handleMobile: Request recieved.")
	code := req.URL.Query().Get("code") //TODO: sanitize
	//TODO: check if code is valid against current lobbys
	if code == "" { //TODO: better what to check this...maybe check if not valid.
		fmt.Println("handleMobile: No code given.")
		templateString := loadTemplate("mobile-no-code")
		templ := template.Must(template.New("mobile-no-code").Parse(templateString))

		templ.Execute(w, req.FormValue("lobbyCode"))
	} else { //TODO: insert the code into the template.

		templateString := loadTemplate("mobile")
		templ := template.Must(template.New("mobile").Parse(templateString))
		err := templ.ExecuteTemplate(w, "mobile", code) //TODO: look into other ways of doing this.
		if err != nil {
			//log.Fatal(err)
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

		newLobbyChan <- newLobby                                                        //register the new lobby with the hub.
		go lobbyController(newLobby, newLobby.deviceInputChan, newLobby.desktopOutChan) //start the new go process with the channels passed.

		//need to update this template to have javascript etc. that connects and listens to the right channels
		var templateString = loadTemplate("desktop")
		templ := template.Must(template.New("desktop").Parse(templateString)) //render original page

		templ.Execute(w, req.FormValue("nothing-like-me-exists"))
	} else { //else no post
		var templateString = loadTemplate("desktop-no-post")
		templ := template.Must(template.New("desktop-no-post").Parse(templateString)) //render original page

		templ.Execute(w, req.FormValue("nothing-like-me-exists"))
	}
}

func handleSocket(socketConnChan chan socketEntity, ws *websocket.Conn) {
	var message string
	var initialRecievedEntity entity
	var handlerChan (chan deviceUpdate)

	//get the first message
	if err := websocket.Message.Receive(ws, &message); err != nil { //receive the message here
		fmt.Println("Error occured:", err)
		return
	}

	//its the first connect
	//TODO: make this less sloppy.
	fmt.Println("handleSocket: Message -", message)

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

			clientSettings := <-tempLobbySettingsChan //this is where we would get the settings for our lobby and send it to the client.
			fmt.Println(clientSettings)
			// jsonData, err := json.Marshal(deviceLobbySettings{Settings: clientSettings})
			// if err != nil {
			// 	fmt.Println(err)
			// }

			if err := websocket.JSON.Send(ws, deviceLobbySettings{Settings: clientSettings}); err != nil {
				fmt.Println("Error sending settings. Error:", err)
			}

		case errMsgChan := <-tempErrorChan: //an error occured
			deviceErrMsg := deviceErrorMessage{Error: errMsgChan, Timestamp: time.Now().Unix()}

			if err := websocket.JSON.Send(ws, deviceErrMsg); err != nil {
				fmt.Println("Error occured:", err)
			}

			//close the connection
			if err := ws.Close(); err != nil {
				fmt.Println(err)
			}
		} //end select

		//update the client here if not desktop
		if initialRecievedEntity.IsDesktop == false { //its a client

		}

	} //end check

	//we test here to see if the websocket was closed.
	if err := websocket.Message.Send(ws, []byte("{ping:true}")); err != nil {
		fmt.Println("handleSocket: Websocket is closed.")
		return //get the hell out of dodge
	}

	for {
		if initialRecievedEntity.IsDesktop {
			data := <-handlerChan //take recieved values from channel

			//send the data over to the desktop socket. Change the dataJson to a string type.
			if err := websocket.JSON.Send(ws, data); err != nil {
				fmt.Println("Error occured:", err)
				return
			}

		} else if initialRecievedEntity.IsDesktop == false { //its a device
			var recievedUpdate deviceUpdate //setup holder vars
			var rawJSONObjectMap map[string]*json.RawMessage

			recievedUpdate.measurments = make(map[string]float64)
			recievedUpdate.deviceID = initialRecievedEntity.DeviceID
			recievedUpdate.timestamp = time.Now().Unix()

			if err := websocket.Message.Receive(ws, &message); err != nil { //receive the message here from the device
				fmt.Println("Error occured:", err)
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
				recievedUpdate.measurments[key] = value //set it
			}

			fmt.Printf("handleSocket: recievedUpdate - %v\n", recievedUpdate)
			handlerChan <- recievedUpdate //send the update
		} else {
			fmt.Println("no idea how you got here")
		}

	} //end for
}

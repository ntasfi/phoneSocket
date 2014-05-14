package main

type entity struct {
	LobbyID   string
	DeviceID  string
	IsDesktop bool
}

type socketEntity struct {
	entity
	lobbySettings (chan map[string]bool)
	errorChan     (chan string)
	socketChan    (chan (chan deviceUpdate)) //returns the channel from hub that the handler needs
}

type deviceLobbySettings struct {
	Settings map[string]bool
}

type deviceErrorMessage struct {
	Error     string
	Timestamp int64
}

type deviceUpdate struct {
	DeviceID    string
	Measurments map[string]float64
	Timestamp   int64
}

type lobby struct {
	id                 string
	numberAllowedUsers int
	currentNumberUsers int
	settings           map[string]bool
	desktopOutChan     (chan deviceUpdate)
	deviceInputChan    (chan deviceUpdate)
	killChan           (chan bool)
}

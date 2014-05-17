package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func readConfiguration(configurationLocation string) Configuration {
	file, err := os.Open(configurationLocation)
	if err != nil {
		fmt.Println(err)
	}
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	if err := decoder.Decode(&configuration); err != nil {
		fmt.Println(err)
	}
	return configuration
}

type Configuration struct {
	ServerAddress string
	ServerPort    int
	HTTPRoutes    HTTPRouteConfiguration
	Template      TemplateConfiguration
	Lobby         LobbyConfiguration
	Mobile        MobileConfiguration
}

type HTTPRouteConfiguration struct {
	Root       string
	Websocket  string
	Javascript FileRouteConfiguration
	Images     FileRouteConfiguration
}

type FileRouteConfiguration struct {
	Route        string
	RootLocation string
}

type TemplateConfiguration struct {
	RootLocation string
	Extension    string
}

type LobbyConfiguration struct {
	InactivityNoActivityMinutes    time.Duration
	InactivityAfterActivitySeconds time.Duration
}

type MobileConfiguration struct {
	DefaultUpdateFrequencyMilliseconds int
}

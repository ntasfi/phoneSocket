package main

type device struct {
	measurments map[string]float32
}

type lobby struct {
	id                 string
	numberAllowedUsers int
	listenFor          map[string]bool
	devices            []device
}

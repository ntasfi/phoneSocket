phoneSocket
===========

An experiment in adding mobile inputs to web games. 

This client & server combination basically renders any (recently manufactured) mobile device into a wiimote-like controller. 


Build:
go build phoneSocketServer

Run:
./phoneSocketServer 


Visit localhost:8080/desktop to setup your lobby
Then visit localhost:8080/mobile?code=<here> to join the lobby on a mobile device

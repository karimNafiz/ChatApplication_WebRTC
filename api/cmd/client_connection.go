package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{

	/*
		this is mainly for development
		need to decide if I want this to be within the app struct

	*/
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (app *application) clientRegisterHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) clientLoginHandler(w http.ResponseWriter, r *http.Request) {

}

/*

	clients need to be authenticated by jwt tokens

*/

func (app *application) clientEstablishWebSocket(w http.ResponseWriter, r *http.Request) {
	_, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		app.logError(r, fmt.Errorf("could not establish web socket connection %w ", err))
		return
	}
	/*
		need to add more information about the client in this log
	*/
	app.logger.PrintInfo("established webSocket connection with client", nil)

}

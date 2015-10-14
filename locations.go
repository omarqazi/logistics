package main

import (
	"fmt"
	"github.com/omarqazi/logistics/auth"
	"github.com/omarqazi/logistics/datastore"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

var webSocketHandler = websocket.Handler(WebSocketServer)

type LocationsController struct {
}

func (l LocationsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := datastore.GetUser(r.URL.Path)
	if err != nil {
		log.Println("Error getting user:", err)
		http.Error(w, "Error getting user", 500)
		return
	}
	if user == nil {
		http.Error(w, "User does not exist", 404)
		return
	}

	rsaKey, err := user.RSAKey()
	if err != nil {
		http.Error(w, "Error authenticating", 500)
		return
	}

	if ok := auth.Request(w, r, rsaKey); !ok {
		return
	}

	webSocketHandler.ServeHTTP(w, r)
}

func WebSocketServer(ws *websocket.Conn) {
	r := ws.Request()
	user, err := datastore.GetUser(r.URL.Path)
	if err != nil || user == nil {
		return
	}

	changeChannel := user.ChangeJSON()
	for {
		payload, ok := <-changeChannel
		if !ok {
			return
		}
		if _, err := fmt.Fprintln(ws, payload); err != nil {
			return // client probably closed socket
		}
	}
}

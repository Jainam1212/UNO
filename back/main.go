package main

import (
	"log"
	"net/http"

	"example.com/models"
	"example.com/routes"
	"example.com/utils"
	"github.com/fasthttp/router"
	"github.com/gorilla/websocket"
	"github.com/valyala/fasthttp"
)

var InitialCards []models.Card

type Message struct {
	Type string `json:"type"`
	User string `json:"user"`
	Text string `json:"text"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all (for dev)
	},
}

func main() {
	InitialCards = utils.GenerateUnoDeck()
	InitialCards = utils.ShuffleCards(InitialCards)
	r := router.New()
	routes.InitializeV1Routes(r)
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()
	fasthttp.ListenAndServe(":8080", r.Handler)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	broadcast <- Message{
		Type: "join",
		User: "New User",
	}

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}

		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

package main

import (
	"fmt"
	"log"

	"example.com/api"
	"example.com/internals/utils"
	"example.com/store"
	"github.com/fasthttp/router"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

func main() {
	store.InitialCards = utils.GenerateUnoDeck()
	store.InitialCards = utils.ShuffleCards(store.InitialCards)
	r := router.New()

	api.InitializeV1Routes(r)

	r.GET("/playerHandler", handleConnections)
	go handleChats()
	go alertHandler()
	go sendPlayerListInfo()
	go inGamePlayersListHandler()
	// go gameStateUpdateHandler()

	fmt.Println("server listening on port 8080")
	handler := CORS(r.Handler)
	fasthttp.ListenAndServe(":8080", handler)
}

var upgrader = websocket.FastHTTPUpgrader{
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

func handleConnections(ctx *fasthttp.RequestCtx) {
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		addClient(conn)
		log.Println("New client connected:")
		for {
			var msg interface{}
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println("read error:", err)
				removeClient(conn)
				break
			}
			if data, ok := msg.(map[string]interface{}); ok {
				channelHandler(data)
			} else {
				log.Println("Received non-object JSON, ignoring.")
			}
		}
	})

	if err != nil {
		log.Println("upgrade error:", err)
	}
}

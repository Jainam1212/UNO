package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"sync"

	"example.com/internals/models"
	"example.com/store"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

var (
	GameData    = make(chan models.GameInfo)
	alerts      = make(chan models.ChannelSendAlert)
	broadcastMu sync.Mutex
)

func handleChats() {
	for {
		Msg := <-store.Broadcast
		broadcastMu.Lock()
		for _, client := range store.Clients {
			client.Send <- Msg
		}
		broadcastMu.Unlock()
	}
}

func channelHandler(data map[string]interface{}) {
	msgType, ok := data["messageType"].(string)
	if !ok {
		return
	}
	switch msgType {
	case "message":
		name := data["userName"].(string)
		text := data["userMessage"].(string)
		store.Broadcast <- models.ChannelSendMessage{
			Type:       "players_chat",
			PlayerName: name,
			Message:    text,
		}
	case "startGame":

	}
}

func alertHandler() {
	for {
		Msg := <-alerts
		broadcastMu.Lock()
		for _, client := range store.Clients {
			client.Send <- Msg
		}
		broadcastMu.Unlock()
	}
}
func sendPlayerListInfo() {
	for {
		Msg := <-store.NewPlayer
		broadcastMu.Lock()
		for _, client := range store.Clients {
			client.Send <- struct {
				Type    string   `json:"type"`
				Players []string `json:"newPlayerIds"`
			}{
				Type:    "add_player_to_list",
				Players: Msg,
			}
		}
		broadcastMu.Unlock()
	}
}

func inGamePlayersListHandler() {
	for {
		Msg := <-store.InGamePlayersList
		broadcastMu.Lock()
		for _, client := range store.Clients {
			client.Send <- struct {
				Type    string                     `json:"type"`
				Players []models.InGamePlayersInfo `json:"inGamePlayerIds"`
			}{
				Type:    "ingame_player_list",
				Players: Msg,
			}
		}
		broadcastMu.Unlock()
	}
}

// func gameStateUpdateHandler() {
// 	for {
// 		Msg := <-store.GameStateUpdater
// 		broadcastMu.Lock()
// 		for _, client := range store.Clients {
// 			client.Send <- struct {
// 				Type    string                     `json:"type"`
// 				Players []models.InGamePlayersInfo `json:"inGamePlayerIds"`
// 			}{
// 				Type:    "ingame_player_list",
// 				Players: Msg,
// 			}
// 		}
// 		broadcastMu.Unlock()
// 	}
// }

func addClient(conn *websocket.Conn) {
	min := 1000000000
	max := 9999999999
	randomNumber := rand.IntN(max-min+1) + min
	client := &store.Client{
		Conn:     conn,
		PlayerId: randomNumber,
		Send:     make(chan interface{}, 10),
	}
	store.Mu.Lock()
	fmt.Println("new client joined - ", randomNumber)
	store.Clients[conn] = client
	store.Mu.Unlock()

	go func(c *store.Client) {
		defer fmt.Println("Writer stopped for:", c.PlayerId)
		for msg := range c.Send {
			err := c.Conn.WriteJSON(msg)
			if err != nil {
				fmt.Println("read error:", err)
				return
			}
		}
	}(client)
	client.Send <- models.ChannelSendUserInfo{
		PlayerId: randomNumber,
		Type:     "user_info",
	}

	store.Mu.Lock()
	var players []string
	for _, v := range store.Clients {
		players = append(players, strconv.Itoa(v.PlayerId))
	}
	store.NewPlayer <- players
	store.Mu.Unlock()
}

func removeClient(conn *websocket.Conn) {
	store.Mu.Lock()
	client, exists := store.Clients[conn]
	if !exists {
		store.Mu.Unlock()
		return
	}
	delete(store.Clients, conn)
	store.Mu.Unlock()
	close(client.Send)

	conn.Close()

	alerts <- models.ChannelSendAlert{
		Type:     "player_left",
		Message:  strconv.Itoa(client.PlayerId) + " left",
		PlayerId: client.PlayerId,
	}

	store.GameStateMutex.Lock()

	index := -1

	for i, v := range store.GameState.Players {
		if v.Pid == client.PlayerId {
			index = i
		}
	}
	if index != -1 {
		store.GameState.Players = append(store.GameState.Players[:index], store.GameState.Players[index+1:]...)
	}
	store.GameStateMutex.Unlock()
	store.Mu.Lock()
	var players []string
	for _, v := range store.Clients {
		players = append(players, strconv.Itoa(v.PlayerId))
	}
	store.NewPlayer <- players
	store.Mu.Unlock()
}

func CORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		origin := string(ctx.Request.Header.Peek("Origin"))
		if origin != "" {
			ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
		}

		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		ctx.Response.Header.Set("Access-Control-Expose-Headers", "*")

		if string(ctx.Method()) == fasthttp.MethodOptions {
			ctx.SetStatusCode(fasthttp.StatusOK)
			return
		}

		next(ctx)
	}
}

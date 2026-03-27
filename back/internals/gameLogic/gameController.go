package gamelogic

import (
	"encoding/json"
	"fmt"

	"example.com/internals/models"
	"example.com/internals/utils"
	"example.com/store"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

func InitGame(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) != fasthttp.MethodPost {
		ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	var data struct {
		PlayerId   int    `json:"playerId"`
		GameId     int    `json:"gameId"`
		PlayerName string `json:"playerName"`
		MaxPlayers int    `json:"maxPlayers"`
	}

	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		utils.JSONResponseWrite(ctx, fasthttp.StatusBadRequest, struct {
			Status  string
			Message string
		}{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	if data.MaxPlayers > 6 {
		utils.JSONResponseWrite(ctx, fasthttp.StatusBadRequest, struct {
			Status  string
			Message string
		}{
			Status:  "fail",
			Message: "only maximum 6 players are allowed",
		})
		return
	}
	store.GameStateMutex.Lock()
	fmt.Println("tryin to create game", data)
	defer store.GameStateMutex.Unlock()
	store.GameState.RoomId = data.GameId
	store.GameState.MaxPlayers = data.MaxPlayers
	store.GameState.TurnInfo.TurnFlowIsReverse = false
	store.GameState.PileCards = utils.ShuffleCards(utils.GenerateUnoDeck())
	store.GameState.Players = append(store.GameState.Players, models.PlayerInfo{Pid: data.PlayerId, Name: data.PlayerName, Role: "owner", CardsInHand: []models.Card{}})

	utils.JSONResponseWrite(ctx, fasthttp.StatusOK, struct {
		Status  string
		Message string
	}{
		Status:  "success",
		Message: "game created successfully",
	})
	player := models.InGamePlayersInfo{
		GamePlayersInfo: models.GamePlayersInfo{
			PlayerName: data.PlayerName,
			PlayerId:   data.PlayerId,
		},
		PlayerRole: "owner",
	}
	store.InGamePlayersList <- []models.InGamePlayersInfo{player}
	store.Broadcast <- models.ChannelSendMessage{
		Type:       "players_chat",
		PlayerName: data.PlayerName,
		Message:    "joined the game",
	}
}

func JoinGame(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) != fasthttp.MethodPost {
		ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
		return
	}
	var data struct {
		PlayerId   int    `json:"playerId"`
		GameId     int    `json:"gameId"`
		PlayerName string `json:"playerName"`
	}

	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		utils.JSONResponseWrite(ctx, fasthttp.StatusInternalServerError, struct {
			Status  string
			Message string
		}{
			Status:  "fail",
			Message: "some error occured while reading body",
		})
		return
	}

	store.GameStateMutex.Lock()
	if store.GameState.RoomId != data.GameId {
		store.GameStateMutex.Unlock()
		utils.JSONResponseWrite(ctx, fasthttp.StatusBadRequest, struct {
			Status  string
			Message string
		}{
			Status:  "fail",
			Message: "no game found with provided id",
		})
		return
	}
	store.GameState.Players = append(store.GameState.Players, models.PlayerInfo{Pid: data.PlayerId, Name: data.PlayerName, Role: "player", CardsInHand: []models.Card{}})
	player := models.InGamePlayersInfo{
		GamePlayersInfo: models.GamePlayersInfo{
			PlayerName: data.PlayerName,
			PlayerId:   data.PlayerId,
		},
		PlayerRole: "player",
	}
	store.InGamePlayersList <- []models.InGamePlayersInfo{player}
	fmt.Println(store.GameState.Players)
	store.GameStateMutex.Unlock()
	store.Broadcast <- models.ChannelSendMessage{
		Type:       "players_chat",
		PlayerName: data.PlayerName,
		Message:    "joined the game",
	}
	utils.JSONResponseWrite(ctx, fasthttp.StatusOK, struct {
		Status  string
		Message string
	}{
		Status:  "success",
		Message: "joined game successfully",
	})

}

func StartGame(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) != fasthttp.MethodPost {
		ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
		return
	}
	var PlayerList []models.GamePlayersInfo
	store.GameStateMutex.Lock()
	if len(store.GameState.Players) <= 1 {
		store.GameStateMutex.Unlock()
		utils.JSONResponseWrite(ctx, fasthttp.StatusBadRequest, struct {
			Status  string
			Message string
		}{
			Status:  "fail",
			Message: "no enough players",
		})
		return
	}
	for _, v := range store.GameState.Players {
		PlayerList = append(PlayerList, models.GamePlayersInfo{
			PlayerName: v.Name,
			PlayerId:   v.Pid,
		})
		subCards := store.GameState.PileCards[0:7]
		v.CardsInHand = subCards
		store.GameState.PileCards = store.GameState.PileCards[7:]
	}
	store.GameState.TurnInfo.Players = PlayerList
	store.GameState.TurnInfo.Winners = []models.GamePlayersInfo{}
	store.GameState.TurnInfo.CurrentTurnPosition = 0
	store.GameState.TurnInfo.CurrentTurn = store.GameState.TurnInfo.Players[0].PlayerId
	store.GameStateMutex.Unlock()
	utils.InitGameInfoHandler()
}

func LeaveGame(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) != fasthttp.MethodPost {
		ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
		return
	}
	var data struct {
		PlayerId int `json:"playerId"`
		GameId   int `json:"gameId"`
	}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		utils.JSONResponseWrite(ctx, fasthttp.StatusInternalServerError, struct {
			Status  string
			Message string
		}{
			Status:  "fail",
			Message: "some error occured while reading body",
		})
		return
	}
	if store.GameState.RoomId != data.GameId {
		utils.JSONResponseWrite(ctx, fasthttp.StatusBadRequest, struct {
			Status  string
			Message string
		}{
			Status:  "fail",
			Message: "invalid game id",
		})
		return
	}
	store.GameStateMutex.Lock()

	index := -1
	// var playerToLeave models.PlayerInfo
	for i, v := range store.GameState.Players {
		if v.Pid == data.PlayerId {
			index = i
		}
		// playerjToLeave = v
	}

	if index == -1 {
		utils.JSONResponseWrite(ctx, fasthttp.StatusBadRequest, struct {
			Status  string
			Message string
		}{
			Status:  "fail",
			Message: "invalid player id",
		})
		return
	}

	store.GameState.Players = append(store.GameState.Players[:index], store.GameState.Players[index+1:]...)
	// if  {

	// }
	store.GameStateMutex.Unlock()

	store.Mu.Lock()
	var connStr *websocket.Conn
	for _, v := range store.Clients {
		if v.PlayerId == data.PlayerId {
			connStr = v.Conn
		}
	}
	delete(store.Clients, connStr)
	store.Mu.Unlock()
	utils.JSONResponseWrite(ctx, fasthttp.StatusOK, struct {
		Status  string
		Message string
	}{
		Status:  "success",
		Message: "left successfully",
	})
}

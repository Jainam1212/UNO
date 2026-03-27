package socketroutes

import (
	"fmt"

	"example.com/models"
	"example.com/store"
	"example.com/utils"
)

func InitGame(playerId int, gameId int, playerName string, maxPlayer int) {

	if maxPlayer > 6 {
		return
	}
	store.GameStateMutex.Lock()
	defer store.GameStateMutex.Unlock()
	store.GameState.RoomId = gameId
	store.GameState.MaxPlayers = maxPlayer
	store.GameState.TurnInfo.TurnFlow = true
	store.GameState.PileCards = utils.ShuffleCards(utils.GenerateUnoDeck())
	store.GameState.Players = append(store.GameState.Players, models.PlayerInfo{Pid: playerId, Name: playerName, Role: "owner", CardsInHand: []models.Card{}})

	player := models.InGamePlayersInfo{
		GamePlayersInfo: models.GamePlayersInfo{
			PlayerName: playerName,
			PlayerId:   playerId,
		},
		PlayerRole: "owner",
	}
	fmt.Println("initing game", player)
	store.InGamePlayersList <- []models.InGamePlayersInfo{player}
	store.Broadcast <- models.ChannelSendMessage{
		Type:       "players_chat",
		PlayerName: playerName,
		Message:    "joined the game",
	}
}

func JoinGame(playerId int, gameId int, playerName string) {
	store.GameStateMutex.Lock()
	if store.GameState.RoomId != gameId {
		store.GameStateMutex.Unlock()
		return
	}
	store.GameState.Players = append(store.GameState.Players, models.PlayerInfo{Pid: playerId, Name: playerName, Role: "player", CardsInHand: []models.Card{}})
	player := models.InGamePlayersInfo{
		GamePlayersInfo: models.GamePlayersInfo{
			PlayerName: playerName,
			PlayerId:   playerId,
		},
		PlayerRole: "player",
	}

	store.InGamePlayersList <- []models.InGamePlayersInfo{player}
	fmt.Println(store.GameState.Players)
	store.GameStateMutex.Unlock()
	store.Broadcast <- models.ChannelSendMessage{
		Type:       "players_chat",
		PlayerName: playerName,
		Message:    "joined the game",
	}
}

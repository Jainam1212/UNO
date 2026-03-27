package store

import (
	"sync"

	"example.com/internals/models"
	"github.com/fasthttp/websocket"
)

var InitialCards []models.Card
var (
	GameStateMutex sync.Mutex
	GameState      models.GameInfo
)

type Client struct {
	Conn     *websocket.Conn
	PlayerId int
	Send     chan interface{}
}

var (
	Mu                sync.Mutex
	Clients           = make(map[*websocket.Conn]*Client)
	NewPlayer         = make(chan []string)
	InGamePlayersList = make(chan []models.InGamePlayersInfo)
	GameStateUpdater  = make(chan models.GameStateUpdater)
	Broadcast         = make(chan models.ChannelSendMessage)
)

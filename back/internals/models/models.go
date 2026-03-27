package models

type Card struct {
	Value string `json:"cardValue"`
	Color string `json:"cardColor"`
}

type GamePlayersInfo struct {
	PlayerName string `json:"playerName"`
	PlayerId   int    `json:"playerId"`
}

type InGamePlayersInfo struct {
	GamePlayersInfo
	PlayerRole string `json:"playerRole"`
}

type GameInfo struct {
	RoomId   int  `json:"roomId"`
	TopCard  Card `json:"topCard"`
	TurnInfo struct {
		CurrentTurn         int
		CurrentTurnPosition int
		TurnFlowIsReverse   bool
		Players             []GamePlayersInfo
		Winners             []GamePlayersInfo
	}
	MaxPlayers int          `json:"maxPlayersLimit"`
	PileCards  []Card       `json:"cardStack"`
	Players    []PlayerInfo `json:"playersList"`
}

type GameStateUpdater struct {
	RoomId   int  `json:"roomId"`
	TopCard  Card `json:"topCard"`
	TurnInfo struct {
		CurrentTurn string
	}
}

type PlayerInfo struct {
	Name        string `json:"playerName"`
	Role        string `json:"playerRole"`
	Pid         int    `json:"playerId"`
	CardsInHand []Card `json:"playerCards"`
}

type ChannelSendMessage struct {
	Type       string `json:"type"`
	PlayerName string `json:"playerName"`
	Message    string `json:"playerMessage"`
}

type ChannelSendAlert struct {
	Type     string `json:"type"`
	Message  string `json:"alertMessage"`
	PlayerId int    `json:"playerId"`
}

type ChannelSendUserInfo struct {
	Type     string `json:"type"`
	PlayerId int    `json:"playerId"`
}

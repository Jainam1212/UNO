package utils

import (
	"math/rand"

	"example.com/internals/models"
	"example.com/store"
)

func GenerateUnoDeck() []models.Card {
	colors := []string{"red", "yellow", "green", "blue"}
	values := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "Skip", "Reverse", "Draw 2"}
	var deck []models.Card

	for _, vc := range colors {
		for _, vv := range values {
			var count int
			if vv == "0" {
				count = 1
			} else {
				count = 2
			}
			for i := 0; i < count; i++ {
				deck = append(deck, models.Card{Value: vv, Color: vc})
			}
		}
	}
	for range 4 {
		deck = append(deck, models.Card{Value: "Wild", Color: "black"})
		deck = append(deck, models.Card{Value: "Wild Draw 4", Color: "black"})
	}
	return deck
}

func ShuffleCards(array []models.Card) []models.Card {
	for i := len(array) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
	return array
}

func InitGameInfoHandler() {
	type Body struct {
		Type        string                   `json:"type"`
		PlayerList  []models.GamePlayersInfo `json:"playerList"`
		CurrentTurn int                      `json:"currentTurn"`
		CardsInHand []models.Card            `json:"cardsInHand"`
	}
	store.GameStateMutex.Lock()
	var InHand = make(map[int][]models.Card)
	for _, v := range store.GameState.Players {
		InHand[v.Pid] = v.CardsInHand
	}

	for _, clientDetails := range store.Clients {
		clientDetails.Send <- Body{
			Type:        "game_init_info",
			PlayerList:  store.GameState.TurnInfo.Players,
			CurrentTurn: store.GameState.TurnInfo.CurrentTurn,
			CardsInHand: InHand[clientDetails.PlayerId],
		}
	}
	store.GameStateMutex.Unlock()
}

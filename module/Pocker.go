package module

import (
	"fmt"
	"math/rand"

	"github.com/gorilla/websocket"
)

type Card struct {
	Value int
	Suit  string
}

var suits = []string{"黑桃", "紅心", "方塊", "梅花"}

var suitRank = map[string]int{
	"黑桃": 1,
	"紅心": 2,
	"方塊": 3,
	"梅花": 4,
}

func GetCard() Card {
	value := rand.Intn(13) + 1
	suit := suits[rand.Intn(4)]
	return Card{Value: value, Suit: suit}
}

func ComparePoker(player1 Message, player2 Message) Message {
	if player1.Card.Value > player2.Card.Value {
		return player1
	}
	if player1.Card.Value < player2.Card.Value {
		return player2
	}
	// 以上是數字有結果情況

	if suitRank[player1.Card.Suit] > suitRank[player2.Card.Suit] {
		return player1
	}
	if suitRank[player1.Card.Suit] < suitRank[player2.Card.Suit] {
		return player2
	}
	// 以上是數字平 花色不平結果情況
	return Message{Name: "平手"} // 平手
}

func handleGame(ws *websocket.Conn, PlayerName string) {
	player := &PlayerInfo{Name: PlayerName}

	fmt.Println("player", player)
}

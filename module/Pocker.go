package module

import (
	"fmt"
	"math/rand"

	"github.com/gorilla/websocket"
)

type Card struct {
	value int
	suit  string
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
	return Card{value: value, suit: suit}
}

func ComparePoker(player1 *PlayerInfo, player2 *PlayerInfo) *PlayerInfo {
	if player1.Poker.value > player2.Poker.value {
		return player1
	}
	if player1.Poker.value < player2.Poker.value {
		return player2
	}
	// 以上是數字有結果情況

	if suitRank[player1.Poker.suit] > suitRank[player2.Poker.suit] {
		return player1
	}
	if suitRank[player1.Poker.suit] < suitRank[player2.Poker.suit] {
		return player2
	}
	// 以上是數字平 花色不平結果情況
	return nil // 平手
}

func handleGame(ws *websocket.Conn, PlayerName string) {
	player := &PlayerInfo{Name: PlayerName, Money: 1000, Poker: GetCard()}
	fmt.Println("player", player)
}

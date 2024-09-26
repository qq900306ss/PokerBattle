package module

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade:", err)
		return
	}
	defer conn.Close()
	_, data, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Failed to read message:", err)
		return
	}
	msg := PlayerInfo{}
	err = json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println("Failed to unmarshal message:", err)
		return
	}
	PlayerName := msg.Name
	handleGame(conn, PlayerName)

}

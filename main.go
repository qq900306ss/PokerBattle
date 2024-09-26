package main

import (
	"fmt"
	"github/qq900306ss-PokerBattle/router"
	"github/qq900306ss-PokerBattle/utils"
)

func main() {

	utils.InitMySQL()
	r := router.Router()

	// http.HandleFunc("/ws", module.WsHandler)
	fmt.Println("WebSocket server started on :8080")

	r.Run(":8080")
}

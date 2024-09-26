package module

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Node struct {
	Conn      *websocket.Conn
	Addr      string // 地址
	DataQueue chan []byte
}

var clientMap map[string]*Node = make(map[string]*Node, 0) // 映射關係

var rwLocker sync.RWMutex

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允許所有來源
		return true
	},
}

func WsHandler(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}
	defer conn.Close()

	_, msg, err := conn.ReadMessage() // 讀取從前端發送的資料

	var receiveData PlayerInfo
	if err := json.Unmarshal(msg, &receiveData); err != nil {
		fmt.Println("unmarshal:", err)
		return
	}

	node := &Node{
		Conn:      conn,
		Addr:      c.Request.RemoteAddr,
		DataQueue: make(chan []byte, 50),
	}

	rwLocker.Lock()
	clientMap[receiveData.Name] = node
	fmt.Println("add client:", receiveData.Name)
	rwLocker.Unlock()

	ConnectOpponent(receiveData, conn)

}

func ListClients(c *gin.Context) {
	rwLocker.RLock()
	defer rwLocker.RUnlock()

	clientCount := len(clientMap)
	clientNames := make([]string, 0, clientCount)
	for name := range clientMap {
		clientNames = append(clientNames, name)
	}

	c.JSON(http.StatusOK, gin.H{
		"client_count": clientCount,
		"clients":      clientNames,
	})
}

func ConnectOpponent(receiveData PlayerInfo, conn *websocket.Conn) {

	for {
		if len(clientMap) >= 2 {
			fmt.Println("配對對手")

			for opponent := range clientMap { //遍歷clientMap 找到對手 處理完之後退出
				if opponent != receiveData.Name {
					fmt.Println("找到對手")

					opponentData := Message{
						Event: "對手ID",
						Name:  opponent,
					}
					opponentInfo, err := json.Marshal(opponentData)
					err = conn.WriteMessage(websocket.TextMessage, opponentInfo)
					if err != nil {
						fmt.Println("Websocket write:", err)
						return
					}
					msg := []byte("yes")
					clientMap[opponent].DataQueue <- msg

				loop:
					for {
						select {
						case data := <-clientMap[receiveData.Name].DataQueue:
							{
								fmt.Println("收到對手傳來的資料:", string(data), receiveData.Name, "可以解放了")
								break loop
							}
						}
					}

					delete(clientMap, receiveData.Name) // 刪除自己
					fmt.Println("刪除自己", receiveData.Name)

					return
					// c.JSON(http.StatusOK, gin.H{
					// 	"opponentName": opponent,
					// })

				}
			}
		}

	}
}

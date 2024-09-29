package module

import (
	"encoding/json"
	"fmt"
	"github/qq900306ss-PokerBattle/utils"
	"net/http"
	"sync"
	"time"

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

	// defer conn.Close()

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

	opponent := ConnectOpponent(receiveData, conn) // 連線對手

	time.Sleep(1 * time.Second) // 刪除等待一陣子以免衝突

	CocosGetCard(receiveData, conn, opponent, node)

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

func ConnectOpponent(receiveData PlayerInfo, conn *websocket.Conn) string {

	for {
		if len(clientMap) >= 2 {
			fmt.Println("配對對手")

			for opponent := range clientMap { //遍歷clientMap 找到對手 處理完之後退出  //感覺要考慮 三個仁 有兩個仁 進行配對時候
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

					return opponent
					// c.JSON(http.StatusOK, gin.H{
					// 	"opponentName": opponent,
					// })

				}
			}
		}

	}
}

func CocosGetCard(receiveData PlayerInfo, conn *websocket.Conn, opponent string, node *Node) {
	rwLocker.Lock()
	clientMap[receiveData.Name] = node
	fmt.Println("在連線完成之後add client:", receiveData.Name)
	rwLocker.Unlock()

	status := Message{}

	MyCard := GetCard()

	status.Event = "可以配卡了"
	status.Name = receiveData.Name
	status.Card = MyCard

	MyCardData := status      //封裝自己卡片資料
	OpponentData := Message{} //封裝對手卡片資料

	statusJson, err := json.Marshal(status)

	if err != nil {
		fmt.Println("json marshal:", err)
		return
	}

	fmt.Println("傳送自己牌資訊給自己", MyCard)

	err = conn.WriteMessage(websocket.TextMessage, statusJson) //寄送自己資料

	Bet := PlayerInfo{}

	for { //讀取下注金額

		_, msg, err := conn.ReadMessage() // 讀取從前端發送的資料

		if err != nil {
			fmt.Println(" 讀取下注資料錯誤Websocket read:", err)
			return
		}

		if err := json.Unmarshal(msg, &Bet); err != nil {
			fmt.Println(" 在下注解析json出錯unmarshal:", err)
			return
		}

		if Bet.Name == "" {
			fmt.Println(receiveData.Name, "放棄下注")
			status.Event = "放棄下注"
			status.Name = ""

			statusJson, err = json.Marshal(status)

			clientMap[opponent].DataQueue <- statusJson

			error := utils.DB.Exec("UPDATE user_basics SET money = money - 10 WHERE username = ?", receiveData.Name)
			if error.Error != nil {
				fmt.Println(receiveData.Name, "在放棄下的Mysql更新錯誤:", error)
			}

		} else {
			fmt.Println("收到下注資料:", Bet.Name, Bet.Money)

		}

		break

	}
	if clientMap[opponent] != nil {

		clientMap[opponent].DataQueue <- statusJson

	}

loop: //後端部分拿到對手資料跟如果收到放棄下注通知前端
	for {
		select {
		case data := <-clientMap[receiveData.Name].DataQueue:
			{

				opponentData := Message{}
				err = json.Unmarshal(data, &opponentData)
				if err != nil {
					fmt.Println("在封裝對手卡片資料json unmarshal:", err)
					return
				}

				if opponentData.Event == "放棄下注" {
					fmt.Println("對手放棄下注")
					status.Event = "對手放棄下注了"
					statusJson, err = json.Marshal(status)
					if err != nil {
						fmt.Println(" 對手放棄下注的封裝json marshal錯誤:", err)
						return
					}
					err = conn.WriteMessage(websocket.TextMessage, statusJson)
					if err != nil {
						fmt.Println(" 對手放棄下注的Websocket write錯誤:", err)
						return
					}

					delete(clientMap, receiveData.Name)
					return //如果對手放棄下注 就直接結束
				}

				OpponentData = opponentData //封裝對手卡片資料 這是由dataqueue 收到的資料 所以是對手的

				status.Event = "收到對手傳送卡片資料"
				status.Card.Value = opponentData.Card.Value
				status.Card.Suit = opponentData.Card.Suit

				statusJson, err = json.Marshal(status)
				if err != nil {
					fmt.Println(" 對手卡片資料封裝json marshal:", err)
					return
				}
				fmt.Println("收到對手:", string(data), receiveData.Name, "可以解放了")
				delete(clientMap, receiveData.Name)
				break loop
			}

		}
	}
	fmt.Println("寄送對手資料給前端", string(statusJson))

	if Bet.Name == "" {
		status.Event = "你放棄下注"
		statusJson, err = json.Marshal(status)
		if err != nil {
			fmt.Println(" 你放棄下注的封裝json marshal錯誤:", err)
			return
		}
		err = conn.WriteMessage(websocket.TextMessage, statusJson) //寄送對手資料給前端

		if err != nil {
			fmt.Println("你放棄下注Websocket write錯誤:", err)
			return
		}

		delete(clientMap, receiveData.Name)

		return
	}

	err = conn.WriteMessage(websocket.TextMessage, statusJson) //寄送對手資料給前端

	if err != nil {
		fmt.Println("Websocket write:", err)
		return
	}

	WhoWin := ComparePoker(MyCardData, OpponentData)

	if WhoWin.Name == receiveData.Name {
		status.Event = "你贏了"
		Bet.Money *= 2

		utils.DB.Where("username = ?").First(&Bet.Name)

		fmt.Println("檢查BET.name", Bet.Name, "尾巴", Bet.Money)
		error := utils.DB.Exec("UPDATE user_basics SET money = money + ? WHERE username = ?", Bet.Money, Bet.Name)
		if error.Error != nil {
			fmt.Println("有啥錯:", error)
		}

		statusJson, err = json.Marshal(status)
		if err != nil {
			fmt.Println("輸贏判定的封裝json marshal:", err)
			return
		}
		err = conn.WriteMessage(websocket.TextMessage, statusJson) //寄送結果給前端
		if err != nil {
			fmt.Println("輸贏判定的傳送錯誤Websocket write:", err)
			return
		}
	} else if WhoWin.Name == opponent {
		status.Event = "你輸了"
		error := utils.DB.Exec("UPDATE user_basics SET money = money - ? WHERE username = ?", Bet.Money, Bet.Name)
		if error.Error != nil {
			fmt.Println("有啥錯:", error)
		}
		statusJson, err = json.Marshal(status)
		if err != nil {
			fmt.Println("輸贏判定的封裝json marshal:", err)
			return
		}
		err = conn.WriteMessage(websocket.TextMessage, statusJson) //寄送結果給前端
		if err != nil {
			fmt.Println("輸贏判定的傳送錯誤Websocket write:", err)
			return
		}

	} else {
		status.Event = "平手"

		error := utils.DB.Exec("UPDATE user_basics SET money = money - 10 WHERE username = ?", Bet.Name)
		if error.Error != nil {
			fmt.Println("有啥錯:", error)
		}
		statusJson, err = json.Marshal(status)
		if err != nil {
			fmt.Println("輸贏判定的封裝json marshal:", err)
			return
		}
		err = conn.WriteMessage(websocket.TextMessage, statusJson) //寄送結果給前端
		if err != nil {
			fmt.Println("輸贏判定的傳送錯誤Websocket write:", err)
			return
		}

	}

	time.Sleep(2 * time.Second)         //最後等個兩秒緩衝
	delete(clientMap, receiveData.Name) // 刪除自己

	return

}

package module

import (
	"fmt"
	"github/qq900306ss-PokerBattle/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserBasic struct {
	ID       uint `gorm:"primarykey"`
	Username string
	Password string
	Sweet    string //加密
	Money    int
}

type PlayerInfo struct {
	Name  string
	Money int
}

type Message struct {
	Event string
	Name  string
}

func LoginHandler(c *gin.Context) { //方法
	data := UserBasic{}

	// 解析 JSON 請求數據
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username := data.Username
	password := data.Password

	user := UserBasic{}

	utils.DB.Where("username = ? and password = ?", username, password).First(&user) // First 只有一個ｆｉｎｄ集合

	fmt.Println("有東西?", user)

	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusOK, gin.H{ //創造一個map json格式 gin.H就是回應
			"code":    -1, // 0是成功 -1是失敗
			"message": "該用戶不存在 或者 密碼錯誤",
			"data":    user,
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{ //創造一個map json格式 gin.H就是回應
		"code":    0, // 0是成功 -1是失敗
		"message": "登入成功",
		"data":    user,
	})

}

func FindUserByName(username string) UserBasic {

	user := UserBasic{}
	utils.DB.Where("username = ?", username).First(&user)
	return user
}

// func GetUserInfoHandler(c *gin.Context) {

// 	data := PlayerInfo{}
// 	if err := c.ShouldBindJSON(&data); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// }

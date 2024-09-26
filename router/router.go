package router

import (
	"github/qq900306ss-PokerBattle/module"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允許所有來源（可根據需要調整）
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.POST("/login", module.LoginHandler)
	// r.POST("/GetUserInfo", module.GetUserInfoHandler)
	r.GET("/ws", module.WsHandler)
	r.GET("/clients", module.ListClients)

	return r
}

package ws

import (
	"chatapp/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterWS(router *gin.Engine) {
	router.GET("/ws", handlers.HandleWebSocket)
}

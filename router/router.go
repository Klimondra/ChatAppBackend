package router

import (
	"chatapp/router/api"
	"chatapp/router/ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
)

func Init() *gin.Engine {
	router := gin.Default()

	frontendUrl := os.Getenv("FRONTEND_URL")
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendUrl},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	api.RegisterContactApi(router)
	ws.RegisterWS(router)

	return router
}

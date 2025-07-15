package api

import (
	"chatapp/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterContactApi(router *gin.Engine) {
	contact := router.Group("/api/contact")
	{
		contact.POST("/list", handlers.GetContactList)
		contact.POST("/new", handlers.NewContact)
	}
}

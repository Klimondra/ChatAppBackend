package handlers

import (
	"chatapp/database"
	"chatapp/models"
	"github.com/gin-gonic/gin"
	"log"
)

type NewContactRequest struct {
	ContactId int `json:"idToContact"`
	UserID    int `json:"userId"`
}

func NewContact(context *gin.Context) {
	var request NewContactRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	var roomIDs []int
	err := database.DB.
		Table("chat_members").
		Select("room_id").
		Where("user_id IN ?", []int{request.UserID, request.ContactId}).
		Group("room_id").
		Having("COUNT(DISTINCT user_id) = 2").
		Pluck("room_id", &roomIDs).Error

	if err != nil {
		log.Fatal(err)
	}

	if len(roomIDs) > 0 {
		context.JSON(200, gin.H{"roomId": roomIDs[0]})
		return
	}

	newRoom := models.ChatRoom{
		Type: "private",
	}

	result := database.DB.Create(&newRoom)
	if result.Error != nil {
		context.JSON(500, gin.H{"error": "Failed to create new room"})
		return
	}

	newMember1 := models.ChatMember{
		UserID: request.UserID,
		RoomID: newRoom.ID,
	}
	result = database.DB.Create(&newMember1)
	if result.Error != nil {
		context.JSON(500, gin.H{"error": "Failed to add first member to room"})
		return
	}

	newMember2 := models.ChatMember{
		UserID: request.ContactId,
		RoomID: newRoom.ID,
	}
	result = database.DB.Create(&newMember2)
	if result.Error != nil {
		context.JSON(500, gin.H{"error": "Failed to add second member to room"})
		return
	}

	context.JSON(200, gin.H{
		"roomId": newRoom.ID,
	})
}

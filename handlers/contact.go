package handlers

import (
	"chatapp/database"
	"chatapp/models"
	"chatapp/models/dto"
	"github.com/gin-gonic/gin"
	"sort"
	"sync"
	"time"
)

type userId struct {
	Id int `json:"userId"`
}

func GetContactList(context *gin.Context) {
	var id userId
	if err := context.ShouldBindJSON(&id); err != nil {
		context.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var roomsWithUser []models.ChatMember
	result := database.DB.Where("user_id = ?", id.Id).Find(&roomsWithUser)
	if result.Error != nil {
		context.JSON(500, gin.H{"error": "No rooms found"})
	}

	var contactList dto.ContactList

	var wg sync.WaitGroup
	for _, room := range roomsWithUser {
		wg.Add(1)

		go func(room models.ChatMember) {
			defer wg.Done()

			var roomDetails models.ChatRoom
			result = database.DB.Where("id = ?", room.RoomID).First(&roomDetails)
			if result.Error != nil {
				context.JSON(500, gin.H{"error": "Room not found"})
				return
			}

			if roomDetails.Type == "private" {
				var wgContacts sync.WaitGroup

				wgContacts.Add(1)
				var secondMember models.ChatMember
				go func() {
					defer wgContacts.Done()

					result = database.DB.Where("room_id = ? AND user_id != ?", room.RoomID, id.Id).First(&secondMember)
					if result.Error != nil {
						context.JSON(500, gin.H{"error": "Second member not found"})
						return
					}
					database.DB.Preload("User").Find(&secondMember)
				}()

				wgContacts.Add(1)
				var lastMessage models.Message
				go func() {
					defer wgContacts.Done()
					result = database.DB.Where("chat_room_id = ?", room.RoomID).First(&lastMessage)
				}()

				wgContacts.Wait()

				contactList.Contacts = append(contactList.Contacts, dto.ContactForList{
					RoomID:             room.RoomID,
					RecipientName:      secondMember.User.Name,
					LastMessageContent: lastMessage.Content,
					LastMessageTime:    lastMessage.Timestamp.Format(time.RFC3339),
				})

			} else {
				return
			}
		}(room)
	}

	wg.Wait()
	sort.Slice(contactList.Contacts, func(i, j int) bool {
		parsedI, _ := time.Parse(time.RFC3339, contactList.Contacts[i].LastMessageTime)
		parsedJ, _ := time.Parse(time.RFC3339, contactList.Contacts[j].LastMessageTime)
		return parsedJ.Before(parsedI)
	})

	context.JSON(200, contactList)
}

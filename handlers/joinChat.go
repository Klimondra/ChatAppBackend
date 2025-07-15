package handlers

import (
	"chatapp/database"
	"chatapp/models"
	"chatapp/models/dto"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type joinChatRequest struct {
	ChatID string `json:"chatId"`
	UserId int    `json:"userId"`
}

func JoinChat(conn *websocket.Conn, msg dto.IncomingMessage) {
	var request joinChatRequest
	if err := json.Unmarshal(msg.Data, &request); err != nil {
		log.Println("❌ Neplatný JSON:", err)
	}

	var canJoin bool
	err := database.DB.
		Model(&models.ChatMember{}).
		Select("1").
		Where("room_id = ? AND user_id = ?", request.ChatID, request.UserId).
		Limit(1).
		Find(&canJoin).Error

	if err != nil {
		log.Println(err)
	} else if !canJoin {
		return
	}

	var wg sync.WaitGroup

	var messages []models.Message
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = database.DB.
			Where("chat_room_id", request.ChatID).
			Find(&messages).Error

		if err != nil {
			log.Println("Error fetching messages:", err)
			return
		}
	}()

	var secondMember models.ChatMember
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = database.DB.
			Where("user_id != ? AND room_id = ?", request.UserId, request.ChatID).
			First(&secondMember).Error

		database.DB.Preload("User").Find(&secondMember)

		if err != nil {
			log.Println("Error fetching second member:", err)
			return
		}
	}()

	wg.Wait()

	var messagesResponse []dto.MessageResponse

	key, _ := GetKey()

	for _, message := range messages {
		var senderName string

		if request.UserId == message.SenderID {
			senderName = "You"
		} else {
			senderName = secondMember.User.Name
		}

		decryptedMessage, errEncrypt := Decrypt(message.Content, key)
		if errEncrypt != nil {
			log.Println("❌ Chyba při šifrování zprávy:", errEncrypt)
			return
		}

		messagesResponse = append(messagesResponse, dto.MessageResponse{
			MessageId: message.ID,
			Content:   decryptedMessage,
			Sender:    senderName,
			Timestamp: message.Timestamp.Format(time.RFC3339),
		})
	}

	response := dto.MessageList{
		Type:     "joinAccepted",
		ChatName: secondMember.User.Name,
		Messages: messagesResponse,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Println("❌ Chyba při serializaci JSON:", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		log.Println("❌ Chyba při odesílání zprávy:", err)
		return
	}
}

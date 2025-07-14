package handlers

import (
	"chatapp/database"
	"chatapp/models"
	"chatapp/models/dto"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
)

type sendMessageRequest struct {
	ChatID    string `json:"chatId"`
	UserId    int    `json:"userId"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type newMessageResponse struct {
	Type    string              `json:"type"`
	Message dto.MessageResponse `json:"message"`
}

func SendMessage(conn *websocket.Conn, msg dto.IncomingMessage) {
	var request sendMessageRequest
	if err := json.Unmarshal(msg.Data, &request); err != nil {
		log.Println("❌ Neplatný JSON:", err)
	}

	var canSend bool
	err := database.DB.
		Model(&models.ChatMember{}).
		Select("1").
		Where("room_id = ? AND user_id = ?", request.ChatID, request.UserId).
		Limit(1).
		Find(&canSend).Error

	if err != nil {
		log.Println(err)
	} else if !canSend {
		return
	}

	roomId, errConvChatId := strconv.Atoi(request.ChatID)
	if errConvChatId != nil {
		log.Println("❌ Chyba při převodu ID místnosti:", errConvChatId)
		return
	}

	timestamp, errConvTime := time.Parse(time.RFC3339, request.Timestamp)
	if errConvTime != nil {
		log.Println("❌ Chyba při převodu času:", errConvTime)
		return
	}

	newMessage := models.Message{
		ChatRoomID: roomId,
		SenderID:   request.UserId,
		Content:    request.Content,
		Timestamp:  timestamp,
	}

	errDb := database.DB.Create(&newMessage).Error
	database.DB.Preload("Sender").First(&newMessage, "sender_id = ?", request.UserId)

	if errDb != nil {
		log.Println("❌ Chyba při ukládání zprávy:", errDb)
		return
	}

	var otherMembers []models.ChatMember
	result := database.DB.Where("room_id = ? AND user_id != ?", request.ChatID, request.UserId).Find(&otherMembers)
	if result.Error != nil {
		log.Println("❌ Chyba při načítání členů místnosti:", result.Error)
		return
	}

	for _, otherMember := range otherMembers {
		otherConn, exists := Clients[otherMember.UserID]
		if !exists {
			continue
		}

		responseMessage := dto.MessageResponse{
			MessageId: newMessage.ID,
			Content:   newMessage.Content,
			Sender:    newMessage.Sender.Name,
			Timestamp: newMessage.Timestamp.Format(time.RFC3339),
		}

		response := newMessageResponse{
			Type:    "newMessage",
			Message: responseMessage,
		}

		jsonData, err := json.Marshal(response)
		if err != nil {
			log.Println("❌ Chyba při serializaci JSON:", err)
			continue
		}

		if err := otherConn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			log.Println("❌ Chyba při odesílání zprávy:", err)
			return
		}

		log.Println("✅ Zpráva odeslána:", jsonData)
	}

}

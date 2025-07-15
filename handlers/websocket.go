package handlers

import (
	"chatapp/models/dto"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strconv"
)

type LiveSession struct {
	RoomID int `json:"roomId"`
	UserID int `json:"userId"`
}

var Clients = make(map[LiveSession]*websocket.Conn)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == os.Getenv("FRONTEND_URL")
	},
}

func HandleWebSocket(context *gin.Context) {
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	userIdStr := context.Query("user_id")
	userId, errIdToInt := strconv.Atoi(userIdStr)
	if errIdToInt != nil {
		log.Println(errIdToInt)
	}

	roomIdStr := context.Query("room_id")
	roomId, errRoomIdToInt := strconv.Atoi(roomIdStr)
	if errRoomIdToInt != nil {
		log.Println(errRoomIdToInt)
	}

	ls := LiveSession{
		UserID: userId,
		RoomID: roomId,
	}
	Clients[ls] = conn

	defer func() {
		conn.Close()
		delete(Clients, LiveSession{UserID: userId, RoomID: roomId})
		log.Println("WebSocket connection closed:", userId)
	}()

	log.Println("New WebSocket connection established")

	var msg dto.IncomingMessage
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Println("❌ Neplatný JSON:", err)
			continue
		}

		log.Printf("Received message: %+v\n", msg)
		handleWsMessage(conn, msg)
	}

	log.Println("WebSocket connection closed")
}

func handleWsMessage(conn *websocket.Conn, msg dto.IncomingMessage) {
	switch msg.Type {

	case "joinChat":
		JoinChat(conn, msg)

	case "sendMessage":
		SendMessage(conn, msg)
	}
}

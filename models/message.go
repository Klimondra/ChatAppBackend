package models

import "time"

type Message struct {
	ID int `gorm:"primaryKey"`

	ChatRoomID int
	ChatRoom   ChatRoom `gorm:"foreignKey:ChatRoomID;references:ID"`

	SenderID int
	Sender   User `gorm:"foreignKey:SenderID;references:ID"`

	Content string

	Timestamp time.Time
}

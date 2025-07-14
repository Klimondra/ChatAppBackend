package models

type ChatRoom struct {
	ID   int `gorm:"primaryKey"`
	Type string
}

type ChatMember struct {
	ID int `gorm:"primaryKey"`

	UserID int
	User   User `gorm:"foreignKey:UserID;references:ID"`

	RoomID int
	Room   ChatRoom `gorm:"foreignKey:RoomID;references:ID"`
}

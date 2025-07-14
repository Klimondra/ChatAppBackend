package models

type User struct {
	ID    int
	Name  string
	Email string
	Image string
}

func (User) TableName() string {
	return "users"
}

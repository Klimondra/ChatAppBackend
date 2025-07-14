package database

import (
	"chatapp/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func ConnectDB() error {
	var err error

	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_NAME"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_SSLMODE"),
		os.Getenv("POSTGRES_TIMEZONE"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	migrationErr := migrateModels()
	if migrationErr != nil {
		return migrationErr
	}

	return nil
}

func migrateModels() error {
	modelsToMgirate := []interface{}{
		&models.ChatRoom{},
		&models.ChatMember{},
		&models.Message{},
	}

	for _, model := range modelsToMgirate {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("Nepodařilo se migrovat model %T: %w", model, err)
		}
	}

	return nil
}

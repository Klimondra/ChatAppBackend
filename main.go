package main

import (
	"chatapp/database"
	"chatapp/router"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	r := router.Init()

	dbErr := database.ConnectDB()
	if dbErr != nil {
		log.Fatal(dbErr)
		return
	}

	port := ":" + os.Getenv("GIN_ROUTER_PORT")
	err := r.Run(port)
	if err != nil {
		log.Fatal(err)
		return
	}
}

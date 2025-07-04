package chat_identity

import (
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
}

func RunServer() {
	log.Println("Server exiting")
}

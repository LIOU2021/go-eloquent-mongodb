package core

import (
	"github/LIOU2021/go-eloquent-mongodb/logger"
	"log"

	"github.com/joho/godotenv"
)

func Setup() {
	logger.Init()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

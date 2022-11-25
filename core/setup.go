package core

import (
	"log"

	"github.com/LIOU2021/go-eloquent-mongodb/logger"

	"github.com/joho/godotenv"
)

func Setup() {
	logger.Init()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

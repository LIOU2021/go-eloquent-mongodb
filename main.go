package main

import (
	"fmt"
	"github/LIOU2021/go-eloquent-mongodb/logger"
	"github/LIOU2021/go-eloquent-mongodb/models"
	"github/LIOU2021/go-eloquent-mongodb/repositories"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	logger.Init()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	userRepo := repositories.NewUserRepository()
	userAll := []*models.User{}
	userAllOk := userRepo.Orm.All(&userAll)

	if !userAllOk {
		fmt.Println("user all query fail !")
	}

	for i, v := range userAll {
		logger.LogDebug.Infof("index : %d, id : %s, name : %s, age : %d\n", i, v.ID, v.Name, v.Age)
	}

	logger.Close()
}

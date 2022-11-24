package main

import (
	"github/LIOU2021/go-eloquent-mongodb/logger"
	"github/LIOU2021/go-eloquent-mongodb/services"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	initial()

	userService := services.NewUserService()
	userAll, _ := userService.All()
	logger.LogDebug.Info("[userService@All]")
	for i, v := range userAll {
		logger.LogDebug.Infof("index : %d, id : %s, name : %s, age : %d\n", i, v.ID, v.Name, v.Age)
	}

	userFind, _ := userService.Find("637f9af09c9fe8db039b144b")
	logger.LogDebug.Infof("[userService@Find] - id : %s, name : %s, age : %d\n", userFind.ID, userFind.Name, userFind.Age)

	close()
}

func initial() {
	logger.Init()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func close() {
	logger.Close()
}

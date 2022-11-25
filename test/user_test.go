package test

import (
	"github/LIOU2021/go-eloquent-mongodb/logger"
	"github/LIOU2021/go-eloquent-mongodb/services"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setup() {
	logger.Init()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func cleanup() {
	logger.Close()
}

func Test_All(t *testing.T) {
	setup()
	defer cleanup()

	userService := services.NewUserService()
	userAll, _ := userService.All()
	logger.LogDebug.Info("[userService@All]")
	for i, v := range userAll {
		logger.LogDebug.Infof("index : %d, id : %s, name : %s, age : %d\n", i, v.ID, v.Name, v.Age)
		assert.True(t, v.ID != "")
	}
}

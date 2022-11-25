package test

import (
	"github/LIOU2021/go-eloquent-mongodb/logger"
	"github/LIOU2021/go-eloquent-mongodb/test/services"
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
	userAll, ok := userService.All()
	assert.True(t, ok, "all not ok")
	for i, v := range userAll {
		logger.LogDebug.Infof("index : %d, id : %s, name : %s, age : %d\n", i, v.ID, v.Name, v.Age)
		assert.True(t, v.ID != "", "_id is empty")
	}
}

package main

import (
	"fmt"
	"github/LIOU2021/go-eloquent-mongodb/repositories"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	userRepo := repositories.NewUserRepository()
	userAll := userRepo.Orm.All()

	for i, v := range userAll {
		fmt.Printf("index : %d, value : %v\n", i, v)
	}
}

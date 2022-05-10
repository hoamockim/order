package app

import (
	"github.com/joho/godotenv"
	"log"
	"order/db"
	"order/pkg/cache"
	"order/pkg/configs"
)

func InitDomain() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}
	configs.InitConfig()
	db.InitDB()
	cache.InitRedisClient()
}

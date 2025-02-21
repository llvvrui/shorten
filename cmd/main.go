package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/llvvrui/shortener/handlers"
	"github.com/llvvrui/shortener/storage"
	"github.com/redis/go-redis/v9"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	store := storage.NewRedisStore(rdb)
	handler := handlers.NewHandler(store)

	router := gin.Default()
	router.POST("/api/v1/shorten", handler.Shorten)
	router.GET("/:shortCode", handler.Redirect)
	// router.GET("/api/v1/stats/:shortCode", handler.GetStats)

	router.Run(":8080")
}
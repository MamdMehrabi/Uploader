package main

import (
	"log"

	"github.com/MamdMehrabi/Uploader/config"
	"github.com/MamdMehrabi/Uploader/handlers"
	"github.com/MamdMehrabi/Uploader/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.Load()

	r := gin.Default()

	r.Use(middleware.FileSizeLimit(config.MaxFileSizeMB))

	r.Static("/static", "./public")

	r.GET("/", handlers.HomeHandler)

	r.GET("/api/health", handlers.HealthHandler(config.BotToken))

	r.GET("/api/max-file-size", handlers.MaxFileSizeHandler)

	telegramService := handlers.NewTelegramService(config.BotToken, config.ProxyURL)
	uploadHandler := handlers.NewUploadHandler(telegramService, config.DefaultChatID)
	r.POST("/api/upload", uploadHandler.HandleUpload)

	log.Printf("Server starting on http://127.0.0.1:%s", config.Port)
	log.Printf("Telegram Bot Token: %s", map[bool]string{true: "Configured", false: "Missing - please set in .env"}[config.BotToken != ""])
	if config.ProxyURL != "" {
		log.Printf("Proxy configured: %s", config.ProxyURL)
	} else {
		log.Println("No proxy configured. If Telegram is blocked in your country, set HTTP_PROXY in .env")
	}

	if err := r.Run(":" + config.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken      string
	DefaultChatID string
	ProxyURL      string
	Address       string
	Port          string
	MaxFileSizeMB int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg := &Config{
		BotToken:      os.Getenv("TELEGRAM_BOT_TOKEN"),
		DefaultChatID: os.Getenv("DEFAULT_CHAT_ID"),
		ProxyURL:      os.Getenv("HTTP_PROXY"),
		Address:       os.Getenv("ADDRESS"),
		Port:          os.Getenv("PORT"),
		MaxFileSizeMB: 20, // Default 20MB
	}

	if cfg.ProxyURL == "" {
		cfg.ProxyURL = os.Getenv("HTTPS_PROXY")
	}

	if cfg.Port == "" {
		cfg.Port = "3000"
	}

	if maxSizeStr := os.Getenv("MAX_FILE_SIZE_MB"); maxSizeStr != "" {
		if size, err := strconv.Atoi(maxSizeStr); err == nil && size > 0 {
			cfg.MaxFileSizeMB = size
			log.Printf("File size limit set to %d MB", cfg.MaxFileSizeMB)
		}
	}

	if cfg.BotToken == "" {
		log.Println("Warning: TELEGRAM_BOT_TOKEN not set in environment")
	}

	return cfg
}


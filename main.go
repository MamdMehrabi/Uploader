package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	neturl "net/url"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type UploadResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
	FileID    string `json:"fileId,omitempty"`
	MessageID int    `json:"messageId,omitempty"`
}

type HealthResponse struct {
	Status   string `json:"status"`
	BotToken string `json:"botToken"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	defaultChatID := os.Getenv("DEFAULT_CHAT_ID")
	proxyURL := os.Getenv("HTTP_PROXY")
	if proxyURL == "" {
		proxyURL = os.Getenv("HTTPS_PROXY")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	maxFileSizeMB := 20
	if maxSizeStr := os.Getenv("MAX_FILE_SIZE_MB"); maxSizeStr != "" {
		if size, err := strconv.Atoi(maxSizeStr); err == nil && size > 0 {
			maxFileSizeMB = size
			log.Printf("File size limit set to %d MB", maxFileSizeMB)
		}
	}
	maxFileSizeBytes := int64(maxFileSizeMB) * 1024 * 1024

	if botToken == "" {
		log.Println("Warning: TELEGRAM_BOT_TOKEN not set in environment")
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("maxFileSizeBytes", maxFileSizeBytes)
		c.Set("maxFileSizeMB", maxFileSizeMB)
		c.Next()
	})

	r.Static("/static", "./public")

	r.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})

	r.GET("/api/health", func(c *gin.Context) {
		botStatus := "missing"
		if botToken != "" {
			if len(botToken) >= 10 {
				botStatus = "configured"
			} else {
				botStatus = "invalid (too short)"
			}
		}
		c.JSON(http.StatusOK, HealthResponse{
			Status:   "ok",
			BotToken: botStatus,
		})
	})

	r.GET("/api/max-file-size", func(c *gin.Context) {
		maxSizeMB := c.MustGet("maxFileSizeMB").(int)
		c.JSON(http.StatusOK, gin.H{
			"maxFileSizeMB":    maxSizeMB,
			"maxFileSizeBytes": c.MustGet("maxFileSizeBytes").(int64),
		})
	})

	r.POST("/api/upload", func(c *gin.Context) {
		if botToken == "" {
			c.JSON(http.StatusBadRequest, UploadResponse{
				Success: false,
				Error:   "Telegram Bot Token not configured. Please set TELEGRAM_BOT_TOKEN in your .env file",
			})
			return
		}

		chatID := c.PostForm("chatId")
		if chatID == "" {
			chatID = defaultChatID
		}
		chatID = strings.TrimSpace(chatID)
		if chatID == "" {
			c.JSON(http.StatusBadRequest, UploadResponse{
				Success: false,
				Error:   "Chat ID is required. Provide it in the request or set DEFAULT_CHAT_ID in .env",
			})
			return
		}

		if strings.HasPrefix(chatID, "@") {
			username := strings.TrimLeft(chatID, "@")
			username = strings.TrimSpace(username)
			if username == "" {
				c.JSON(http.StatusBadRequest, UploadResponse{
					Success: false,
					Error:   "Invalid username format. Username cannot be empty",
				})
				return
			}
			chatID = "@" + username
		} else {
			isNumeric := true
			for _, r := range chatID {
				if !unicode.IsDigit(r) && r != '-' {
					isNumeric = false
					break
				}
			}
			if !isNumeric && !strings.HasPrefix(chatID, "@") {
				chatID = "@" + chatID
			}
		}

		caption := c.PostForm("caption")

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, UploadResponse{
				Success: false,
				Error:   "No file provided: " + err.Error(),
			})
			return
		}
		defer file.Close()

		maxSizeBytes := c.MustGet("maxFileSizeBytes").(int64)
		maxSizeMB := c.MustGet("maxFileSizeMB").(int)

		fileSize := header.Size
		if fileSize > maxSizeBytes {
			c.JSON(http.StatusBadRequest, UploadResponse{
				Success: false,
				Error:   fmt.Sprintf("File size (%.2f MB) exceeds the limit of %d MB", float64(fileSize)/(1024*1024), maxSizeMB),
			})
			return
		}

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, UploadResponse{
				Success: false,
				Error:   "Failed to read file: " + err.Error(),
			})
			return
		}

		if strings.HasPrefix(chatID, "@") {
			log.Printf("Sending file to username: %s", chatID)
		} else {
			log.Printf("Sending file to chat ID: %s", chatID)
		}

		result, err := sendFileToTelegram(botToken, chatID, header.Filename, fileBytes, caption, proxyURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, UploadResponse{
				Success: false,
				Error:   "Failed to upload to Telegram: " + err.Error(),
			})
			return
		}

		if !result.Success {
			c.JSON(http.StatusInternalServerError, UploadResponse{
				Success: false,
				Error:   result.Error,
			})
			return
		}

		c.JSON(http.StatusOK, UploadResponse{
			Success:   true,
			Message:   "File uploaded successfully",
			FileID:    result.FileID,
			MessageID: result.MessageID,
		})
	})

	log.Printf("Server starting on http://127.0.0.1:%s", port)
	log.Printf("Telegram Bot Token: %s", map[bool]string{true: "Configured", false: "Missing - please set in .env"}[botToken != ""])
	if proxyURL != "" {
		log.Printf("Proxy configured: %s", proxyURL)
	} else {
		log.Println("No proxy configured. If Telegram is blocked in your country, set HTTP_PROXY in .env")
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func sendFileToTelegram(botToken, chatID, filename string, fileBytes []byte, caption, proxyURL string) (*UploadResponse, error) {
	if botToken == "" {
		return &UploadResponse{
			Success: false,
			Error:   "Telegram Bot Token not configured. Please set TELEGRAM_BOT_TOKEN in your .env file",
		}, nil
	}

	if len(botToken) < 10 {
		return &UploadResponse{
			Success: false,
			Error:   "Invalid Telegram Bot Token format. Token seems too short",
		}, nil
	}

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	if err := writer.WriteField("chat_id", chatID); err != nil {
		return nil, fmt.Errorf("failed to write chat_id: %w", err)
	}

	if caption != "" {
		if err := writer.WriteField("caption", caption); err != nil {
			return nil, fmt.Errorf("failed to write caption: %w", err)
		}
	}

	part, err := writer.CreateFormFile("document", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(fileBytes); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", botToken)
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	if proxyURL != "" {
		proxy, err := neturl.Parse(proxyURL)
		if err != nil {
			return &UploadResponse{
				Success: false,
				Error:   fmt.Sprintf("Invalid proxy URL: %s. Error: %v", proxyURL, err),
			}, nil
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		client = &http.Client{
			Transport: transport,
		}
		log.Printf("Using proxy: %s", proxyURL)
	}

	resp, err := client.Do(req)
	if err != nil {
		if err.Error() == "EOF" ||
			err.Error() == "connection refused" ||
			err.Error() == "no such host" ||
			err.Error() == "timeout" ||
			err.Error() == "context deadline exceeded" {
			return &UploadResponse{
				Success: false,
				Error:   fmt.Sprintf("Cannot connect to Telegram API. This might be because Telegram is blocked in your country. Error: %v. Try setting HTTP_PROXY in your .env file", err),
			}, nil
		}
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var telegramError struct {
			OK          bool   `json:"ok"`
			ErrorCode   int    `json:"error_code"`
			Description string `json:"description"`
		}
		if err := json.Unmarshal(body, &telegramError); err == nil {
			errorMsg := telegramError.Description
			if telegramError.ErrorCode == 401 {
				errorMsg = "Unauthorized: Invalid Telegram Bot Token. Please check your TELEGRAM_BOT_TOKEN in .env file"
			} else if telegramError.ErrorCode == 400 {
				if strings.Contains(telegramError.Description, "chat not found") {
					errorMsg = "Chat not found. Please check your Chat ID is correct. For usernames, make sure:\n1. The user has started a conversation with your bot (/start)\n2. The username is correct (e.g., @username)\n3. The user hasn't blocked your bot"
				} else if strings.Contains(telegramError.Description, "username") || strings.Contains(telegramError.Description, "USERNAME") {
					errorMsg = "Username error: " + telegramError.Description + "\n\nTip: Make sure the username starts with @ and the user has started a conversation with your bot"
				} else {
					errorMsg = "Bad Request: " + telegramError.Description
				}
			}
			return &UploadResponse{
				Success: false,
				Error:   errorMsg,
			}, nil
		}
		return &UploadResponse{
			Success: false,
			Error:   fmt.Sprintf("Telegram API error (HTTP %d): %s", resp.StatusCode, string(body)),
		}, nil
	}

	var telegramResponse struct {
		OK     bool `json:"ok"`
		Result struct {
			MessageID int `json:"message_id"`
			Document  struct {
				FileID string `json:"file_id"`
			} `json:"document"`
		} `json:"result"`
		Description string `json:"description,omitempty"`
	}

	if err := json.Unmarshal(body, &telegramResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !telegramResponse.OK {
		return &UploadResponse{
			Success: false,
			Error:   telegramResponse.Description,
		}, nil
	}

	return &UploadResponse{
		Success:   true,
		FileID:    telegramResponse.Result.Document.FileID,
		MessageID: telegramResponse.Result.MessageID,
	}, nil
}

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	neturl "net/url"
	"strings"

	"github.com/MamdMehrabi/Uploader/models"
)

type TelegramService struct {
	BotToken string
	ProxyURL string
}

func NewTelegramService(botToken, proxyURL string) *TelegramService {
	return &TelegramService{
		BotToken: botToken,
		ProxyURL: proxyURL,
	}
}

func (s *TelegramService) SendFile(chatID, filename string, fileBytes []byte, caption string) (*models.UploadResponse, error) {
	if s.BotToken == "" {
		return &models.UploadResponse{
			Success: false,
			Error:   "Telegram Bot Token not configured. Please set TELEGRAM_BOT_TOKEN in your .env file",
		}, nil
	}

	if len(s.BotToken) < 10 {
		return &models.UploadResponse{
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

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", s.BotToken)
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := s.createHTTPClient()
	if client == nil {
		return &models.UploadResponse{
			Success: false,
			Error:   "Failed to create HTTP client",
		}, nil
	}

	resp, err := client.Do(req)
	if err != nil {
		if s.isConnectionError(err.Error()) {
			return &models.UploadResponse{
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
		return s.handleErrorResponse(body)
	}

	return s.handleSuccessResponse(body)
}

func (s *TelegramService) createHTTPClient() *http.Client {
	client := &http.Client{}
	if s.ProxyURL != "" {
		proxy, err := neturl.Parse(s.ProxyURL)
		if err != nil {
			log.Printf("Invalid proxy URL: %s. Error: %v", s.ProxyURL, err)
			return nil
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		client = &http.Client{
			Transport: transport,
		}
		log.Printf("Using proxy: %s", s.ProxyURL)
	}
	return client
}

func (s *TelegramService) isConnectionError(errMsg string) bool {
	connectionErrors := []string{
		"EOF",
		"connection refused",
		"no such host",
		"timeout",
		"context deadline exceeded",
	}
	for _, ce := range connectionErrors {
		if errMsg == ce {
			return true
		}
	}
	return false
}

func (s *TelegramService) handleErrorResponse(body []byte) (*models.UploadResponse, error) {
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
		return &models.UploadResponse{
			Success: false,
			Error:   errorMsg,
		}, nil
	}
	return &models.UploadResponse{
		Success: false,
		Error:   fmt.Sprintf("Telegram API error: %s", string(body)),
	}, nil
}

func (s *TelegramService) handleSuccessResponse(body []byte) (*models.UploadResponse, error) {
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
		return &models.UploadResponse{
			Success: false,
			Error:   telegramResponse.Description,
		}, nil
	}

	return &models.UploadResponse{
		Success:   true,
		FileID:    telegramResponse.Result.Document.FileID,
		MessageID: telegramResponse.Result.MessageID,
	}, nil
}


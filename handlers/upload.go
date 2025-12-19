package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/MamdMehrabi/Uploader/models"
	"github.com/MamdMehrabi/Uploader/utils"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	telegramService *TelegramService
	defaultChatID   string
}

func NewUploadHandler(telegramService *TelegramService, defaultChatID string) *UploadHandler {
	return &UploadHandler{
		telegramService: telegramService,
		defaultChatID:   defaultChatID,
	}
}

func (h *UploadHandler) HandleUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.UploadResponse{
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
		c.JSON(http.StatusBadRequest, models.UploadResponse{
			Success: false,
			Error:   fmt.Sprintf("File size (%.2f MB) exceeds the limit of %d MB", float64(fileSize)/(1024*1024), maxSizeMB),
		})
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.UploadResponse{
			Success: false,
			Error:   "Failed to read file: " + err.Error(),
		})
		return
	}

	chatID := c.PostForm("chatId")
	if chatID == "" {
		chatID = h.defaultChatID
	}

	chatID = utils.NormalizeChatID(chatID)
	if chatID == "" {
		c.JSON(http.StatusBadRequest, models.UploadResponse{
			Success: false,
			Error:   "Chat ID is required. Provide it in the request or set DEFAULT_CHAT_ID in .env",
		})
		return
	}

	caption := c.PostForm("caption")

	if strings.HasPrefix(chatID, "@") {
		log.Printf("Sending file to username: %s", chatID)
	} else {
		log.Printf("Sending file to chat ID: %s", chatID)
	}

	result, err := h.telegramService.SendFile(chatID, header.Filename, fileBytes, caption)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.UploadResponse{
			Success: false,
			Error:   "Failed to upload to Telegram: " + err.Error(),
		})
		return
	}

	if !result.Success {
		c.JSON(http.StatusInternalServerError, models.UploadResponse{
			Success: false,
			Error:   result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, models.UploadResponse{
		Success:   true,
		Message:   "File uploaded successfully",
		FileID:    result.FileID,
		MessageID: result.MessageID,
	})
}

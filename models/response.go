package models

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

type MaxFileSizeResponse struct {
	MaxFileSizeMB    int   `json:"maxFileSizeMB"`
	MaxFileSizeBytes int64 `json:"maxFileSizeBytes"`
}


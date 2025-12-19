package utils

import (
	"strings"
	"unicode"
)

func NormalizeChatID(chatID string) string {
	chatID = strings.TrimSpace(chatID)
	if chatID == "" {
		return ""
	}

	if strings.HasPrefix(chatID, "@") {
		username := strings.TrimLeft(chatID, "@")
		username = strings.TrimSpace(username)
		if username == "" {
			return ""
		}
		return "@" + username
	}

	// Check if it's numeric
	isNumeric := true
	for _, r := range chatID {
		if !unicode.IsDigit(r) && r != '-' {
			isNumeric = false
			break
		}
	}

	// If not numeric, assume it's a username without @
	if !isNumeric && !strings.HasPrefix(chatID, "@") {
		return "@" + chatID
	}

	return chatID
}


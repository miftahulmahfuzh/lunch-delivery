package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GeneratePasswordResetToken() (string, error) {
	uuid := uuid.New()

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	token := fmt.Sprintf("%s-%s-%d",
		uuid.String(),
		hex.EncodeToString(randomBytes),
		time.Now().Unix())

	return token, nil
}
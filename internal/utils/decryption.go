package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func DecryptFile() (string, string, error) {
	const fixedKeyHex = "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6" // 16-byte AES key in hex

	// Read file
	encData, err := ioutil.ReadFile("config.enc")
	if err != nil {
		return "", "", fmt.Errorf("failed to read file: %w", err)
	}

	// Decode key
	key, err := hex.DecodeString(fixedKeyHex)
	if err != nil {
		return "", "", fmt.Errorf("key decode error: %w", err)
	}

	// Extract IV (12 bytes for GCM)
	if len(encData) < 12 {
		return "", "", fmt.Errorf("invalid payload length")
	}
	iv := encData[:12]
	ciphertext := encData[12:]

	// AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	// GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	// Decrypt
	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", "", fmt.Errorf("decryption failed: %w", err)
	}

	// Parse JSON
	var payload struct {
		UID      string `json:"uid"`
		DeviceID string `json:"deviceId"`
	}
	if err := json.Unmarshal(plaintext, &payload); err != nil {
		return "", "", fmt.Errorf("json parse error: %w", err)
	}

	return payload.UID, payload.DeviceID, nil
}

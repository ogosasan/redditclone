package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(inPass string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(inPass))
	if err != nil {
		return "", err
	}
	hashBytes := hash.Sum(nil)
	hashPass := hex.EncodeToString(hashBytes)
	return hashPass, nil
}

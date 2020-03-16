package generator


import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomToken() string {
	return GenerateRandomString(16)
}

func GenerateRandomString(length int) string {
	return hex.EncodeToString(GenerateRandomBytes(length))
}

func GenerateRandomBytes(length int) []byte {
	token := make([]byte, length)

	_, _ = rand.Read(token)
	return token
}

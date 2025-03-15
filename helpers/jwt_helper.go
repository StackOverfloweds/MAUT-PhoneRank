package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
)

func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("JWT_SECRET not found in .env, generating a new one...")
		secret = generateRandomSecret()
		os.Setenv("JWT_SECRET", secret)
	}
	return []byte(secret)
}

func generateRandomSecret() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)

	}
	return base64.StdEncoding.EncodeToString(b)
}

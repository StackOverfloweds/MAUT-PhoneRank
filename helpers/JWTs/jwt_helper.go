package jwts

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
)

/*
GetJWTSecret - Retrieves the JWT secret key.
This function checks if the JWT secret key is set in the environment variable (`JWT_SECRET`).
If the secret key is missing, it generates a new one and sets it in the environment.
The secret key is then returned as a byte slice.
*/
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		log.Println("JWT_SECRET not found in .env, generating a new one...")
		secret = generateRandomSecret()
		os.Setenv("JWT_SECRET", secret) // Set the generated secret in the environment
	}

	return []byte(secret)
}

/*
generateRandomSecret - Generates a secure random secret key.
This function creates a 32-byte random value using the `crypto/rand` package.
The generated random bytes are then encoded using Base64 for safe storage and usage.
If an error occurs while generating the random bytes, the application will log a fatal error and exit.
*/
func generateRandomSecret() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(b)
}

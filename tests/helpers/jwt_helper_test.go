package helpers

import (
	"os"
	"testing"

	"github.com/StackOverfloweds/MAUT-PhoneRank/helpers"
	"github.com/stretchr/testify/assert"
)

/*
TestGetJWTSecret_Existing - Ensures that GetJWTSecret retrieves the existing key from the environment.
*/
func TestGetJWTSecret_Existing(t *testing.T) {
	expectedSecret := "test_secret"
	os.Setenv("JWT_SECRET", expectedSecret)

	secret := helpers.GetJWTSecret()
	assert.Equal(t, expectedSecret, string(secret))
}

/*
TestGetJWTSecret_Generated - Ensures that GetJWTSecret generates a new key if none exists.
*/
func TestGetJWTSecret_Generated(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	secret := helpers.GetJWTSecret()
	assert.NotEmpty(t, secret)
}

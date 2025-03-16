package middleware_test

import (
	"testing"

	"net/http"
	"net/http/httptest"
	"time"

	"github.com/StackOverfloweds/MAUT-PhoneRank/helpers"
	"github.com/StackOverfloweds/MAUT-PhoneRank/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

/*
TestJWTMiddleware_ValidToken - Tests JWTMiddleware with a valid token.
It ensures that the request passes authentication and extracts user info.
*/
func TestJWTMiddleware_ValidToken(t *testing.T) {
	// Retrieve the correct JWT secret
	secret := helpers.GetJWTSecret() // ✅ Ensures correct secret is used for signing & validation

	// Generate valid token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "12345",
		"role":    "user",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(secret) // ✅ Sign token with correct secret

	// Setup Fiber app with middleware
	app := fiber.New()
	app.Use(middleware.JWTMiddleware())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"user_id": c.Locals("user_id"), "role": c.Locals("role")})
	})

	// Simulate request
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

/*
TestJWTMiddleware_InvalidToken - Tests JWTMiddleware with an invalid token.
It ensures that the request is blocked and returns a 401 error.
*/
func TestJWTMiddleware_InvalidToken(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.JWTMiddleware())

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

/*
TestJWTMiddleware_NoToken - Tests JWTMiddleware with no token.
It ensures that the request is blocked and returns a 401 error.
*/
func TestJWTMiddleware_NoToken(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.JWTMiddleware())

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

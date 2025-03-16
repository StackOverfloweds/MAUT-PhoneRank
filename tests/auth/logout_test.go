package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/StackOverfloweds/MAUT-PhoneRank/tests" // Import shared test helpers
	"github.com/stretchr/testify/assert"
)

func TestLogout_Success(t *testing.T) {
	app := tests.SetupApp()

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer valid-jwt-token")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	expectedResponse := `{"message":"Logout successful. Please remove the token on the client side."}`
	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)

	assert.JSONEq(t, expectedResponse, string(body))
}

func TestLogout_Unauthorized(t *testing.T) {
	app := tests.SetupApp()

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	expectedResponse := `{"error":"Unauthorized"}`
	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)

	assert.JSONEq(t, expectedResponse, string(body))
}

func TestLogout_InvalidTokenFormat(t *testing.T) {
	app := tests.SetupApp()

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Authorization", "InvalidTokenFormat")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	expectedResponse := `{"error":"Invalid Authorization header format"}`
	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)

	assert.JSONEq(t, expectedResponse, string(body))
}

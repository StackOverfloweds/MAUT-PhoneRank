package helpers

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

/*
OTPData - Structure to store OTP details.
This struct contains the OTP code and its expiration time.
*/
type OTPData struct {
	Code      string
	ExpiresAt time.Time
}

// Mutex and map to store OTPs in memory (use Redis for production)
var otpStore = make(map[string]OTPData)
var otpMutex = sync.Mutex{}

/*
GenerateOTP - Creates a random 6-digit OTP.
Uses `math/rand` to generate a number between 000000-999999.
This function seeds `rand` with the current time to ensure randomness.
*/
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

/*
SendOTP - Sends an OTP to the provided phone number via the Fonnte API.
Uses `resty` to make a POST request with the OTP message.
If the request fails or the API does not return a 200 status, an error is returned.
The OTP is then stored in `otpStore` with a 60-second expiration time.
*/
func SendOTP(phoneNumber string, otp string) error {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", os.Getenv("FONNTE_TOKEN")). // Get API token from environment
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"target":  phoneNumber,
			"message": fmt.Sprintf("Your verification code is: %s", otp),
		}).
		Post(os.Getenv("API_URL"))

	if err != nil {
		return fmt.Errorf("failed to send OTP request: %w", err)
	}

	// Check if API response is successful
	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to send OTP request: %s", resp.Status())
	}

	// Store OTP with expiration time
	otpMutex.Lock()
	otpStore[phoneNumber] = OTPData{
		Code:      otp,
		ExpiresAt: time.Now().Add(60 * time.Second),
	}
	otpMutex.Unlock()

	// Start a goroutine to delete OTP after 60 seconds
	go func(phone string) {
		time.Sleep(60 * time.Second)
		otpMutex.Lock()
		delete(otpStore, phone)
		otpMutex.Unlock()
	}(phoneNumber)

	return nil
}

/*
FindPhoneByOTP - Retrieves the phone number associated with an OTP.
Iterates through `otpStore` to find a matching OTP.
Returns the phone number if the OTP is found and has not expired.
If the OTP is invalid or expired, an error is returned.
*/
func FindPhoneByOTP(otp string) (string, error) {
	otpMutex.Lock()
	defer otpMutex.Unlock()

	for phone, data := range otpStore {
		if data.Code == otp && time.Now().Before(data.ExpiresAt) {
			return phone, nil
		}
	}
	return "", errors.New("invalid or expired OTP")
}

/*
DeleteOTP - Removes an OTP entry from `otpStore`.
This function is called after a successful OTP verification.
Uses a mutex lock to ensure thread safety during deletion.
*/
func DeleteOTP(phoneNumber string) {
	otpMutex.Lock()
	defer otpMutex.Unlock()

	delete(otpStore, phoneNumber)
}

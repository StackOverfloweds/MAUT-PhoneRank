package helpers_test

import (
	"testing"
	"time"

	"github.com/StackOverfloweds/MAUT-PhoneRank/helpers"
	"github.com/stretchr/testify/assert"
)

// TestGenerateOTP - Ensure generated OTP is always 6 digits.
func TestGenerateOTP(t *testing.T) {
	otp := helpers.GenerateOTP()
	assert.Len(t, otp, 6, "OTP should be exactly 6 digits long")
}

// TestSendOTP - Ensure OTP is stored and expires correctly.
func TestSendOTP(t *testing.T) {
	phone := "081234567890"
	otp := helpers.GenerateOTP()

	err := helpers.SendOTP(phone, otp)
	assert.NoError(t, err, "Sending OTP should not return an error")

	// Ensure OTP is retrievable
	savedPhone, err := helpers.FindPhoneByOTP(otp)
	assert.NoError(t, err, "OTP should be found")
	assert.Equal(t, phone, savedPhone, "Phone number should match")

	// Wait for OTP to expire
	time.Sleep(61 * time.Second)
	savedPhone, err = helpers.FindPhoneByOTP(otp)
	assert.Error(t, err, "OTP should be expired and not found")
}

// TestFindPhoneByOTP - Ensure correct OTP retrieval.
func TestFindPhoneByOTP(t *testing.T) {
	phone := "081234567891"
	otp := helpers.GenerateOTP()
	helpers.SendOTP(phone, otp)

	savedPhone, err := helpers.FindPhoneByOTP(otp)
	assert.NoError(t, err, "OTP should be retrievable")
	assert.Equal(t, phone, savedPhone, "Retrieved phone number should match")
}

// TestDeleteOTP - Ensure OTP is removed after verification.
func TestDeleteOTP(t *testing.T) {
	phone := "081234567892"
	otp := helpers.GenerateOTP()
	helpers.SendOTP(phone, otp)
	helpers.DeleteOTP(phone)

	savedPhone, err := helpers.FindPhoneByOTP(otp)
	assert.Error(t, err, "OTP should not be found after deletion")
	assert.Empty(t, savedPhone, "Phone number should be empty")
}

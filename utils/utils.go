package utils

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateVerificationCode generates a random verification code.
func GenerateVerificationCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func SendVerificationEmail(email, verificationCode string) {
	return
}

func ValidateEmail(email string) bool {
	// Check email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false
	}
	return true
}

func ValidatePassword(password string) bool {
	// Check password requirements
	if len(password) < 6 || len(password) > 16 {
		return false
	}
	hasUpper, hasLower, hasSpecial := false, false, false
	specialChars := "()[]{}<>+-*/?,.:;\"'_\\|~`!@#$%^&="
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}
	if !hasUpper || !hasLower || !hasSpecial {
		return false
	}

	return true
}

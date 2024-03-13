package utils

import (
	"database/sql"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
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

var jwtKey = []byte("will")

func GenerateJWT(email string) (string, error) {
	// Create a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	// Sign the token with a secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func InitializeRecommendItems(db *sql.DB) error {
	// Insert rows into the recommend_item table
	_, err := db.Exec("INSERT INTO recommend_item (item_name, price) VALUES ('Item 1', 10.99), ('Item 2', 20.99)")

	log.Printf("Insert successfully")
	return err
}

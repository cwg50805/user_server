package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"userServer/models"
	"userServer/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate email format and password requirements
	if !utils.ValidateEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}
	if !utils.ValidatePassword(user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be between 6 and 16 characters and contain at least one uppercase letter, one lowercase letter, and one special character"})
		return
	}

	// Generate a verification code
	verificationCode := utils.GenerateVerificationCode()

	db := c.MustGet("db").(*sql.DB)

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Store the user data and verification code in the database
	_, err = db.Exec("INSERT INTO users (email, password, verification_code) VALUES (?, ?, ?)", user.Email, hashedPassword, verificationCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// TODO: Send the verification code to the user's email
	utils.SendVerificationEmail(user.Email, verificationCode)

	// Registration successful
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func LoginUser(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the user's password hash from the database
	var storedPasswordHash string
	err := db.QueryRow("SELECT password FROM users WHERE email = ?", user.Email).Scan(&storedPasswordHash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate a JWT token
	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func VerifyEmail(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	email := c.Query("email")
	code := c.Query("code")

	// Check if the verification code is correct and update the verification status in the database
	result, err := db.Exec("UPDATE users SET verified = true WHERE email = ? AND verification_code = ?", email, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	// Check if any rows were affected by the update query
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or verification code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

var mutex sync.Mutex

func Recommendation(c *gin.Context) {
	// Use a mutex to block other requests while the query is executed
	mutex.Lock()
	defer mutex.Unlock()

	// Check if data exists in Redis
	redisClient := c.MustGet("redis").(*redis.Client)
	// ctx := c.Request.Context()

	data, err := redisClient.Get("recommendation_data").Result()
	if err == nil {
		var result []models.RecommendationItem
		if err := json.Unmarshal([]byte(data), &result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal data from Redis"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": result})
		return
	} else if err != redis.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data from Redis"})
		return
	}

	// Data does not exist in Redis, query MySQL
	db := c.MustGet("db").(*sql.DB)
	stmt, err := db.Prepare("SELECT item_name, price FROM recommend_item")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		return
	}
	defer rows.Close()

	var result []models.RecommendationItem
	for rows.Next() {
		var itemName string
		var price float64
		if err := rows.Scan(&itemName, &price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		item := models.RecommendationItem{ItemName: itemName, Price: price}
		result = append(result, item)
	}

	// Cache the result in Redis
	jsonData, err := json.Marshal(result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}
	if err := redisClient.Set("recommendation_data", jsonData, 10*time.Minute).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data to Redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

package main

import (
	"database/sql"
	"log"
	"net/http"

	"userServer/handlers"
	"userServer/middleware"
	"userServer/services"
	"userServer/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Initialize Redis
	redisClient := services.InitRedis()

	// Initialize MySQL
	services.InitMySQL()
	db := services.GetMySQL()

	utils.InitializeRecommendItems(db)

	r := gin.Default()

	// Pass the Redis client and database connection pool to the handlers
	r.Use(func(c *gin.Context) {
		c.Set("redis", redisClient)
		c.Set("db", db)
		c.Next()
	})

	// Route handlers
	r.POST("/register", handlers.RegisterUser)
	r.POST("/login", handlers.LoginUser)
	r.GET("/verify_email", handlers.VerifyEmail)
	r.GET("/recommendation", middleware.AuthenticateUser, handlers.Recommendation)

	// Start the server
	port := ":8080"
	log.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}

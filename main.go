package main

import (
	"database/sql"
	"log"
	"net/http"

	"userServer/handlers"
	"userServer/middleware"
	"userServer/services"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "username:password@tcp(hostname:port)/database_name")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Set the maximum number of open connections
	db.SetMaxOpenConns(10)

	// Set the maximum number of idle connections
	db.SetMaxIdleConns(5)

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Initialize recommend_item table
	if err := initializeRecommendItems(db); err != nil {
		log.Fatal("Failed to initialize recommend_item table:", err)
	}
}

func initializeRecommendItems(db *sql.DB) error {
	// Insert rows into the recommend_item table
	// Example:
	_, err := db.Exec("INSERT INTO recommend_item (item_name, price) VALUES ('Item 1', 10.99), ('Item 2', 20.99)")
	return err
}

func main() {
	// Initialize Redis
	redisClient := services.InitRedis()

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

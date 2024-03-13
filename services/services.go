package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	"userServer/utils"

	"github.com/avast/retry-go"
	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func InitRedis() *redis.Client {
	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return redisClient
}

var db *sql.DB

func InitMySQL() {
	username := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	dbName := os.Getenv("MYSQL_DATABASE")

	err := retry.Do(
		func() error {
			var err error
			db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", username, password, host, dbName))
			if err != nil {
				return err
			}

			// Set the maximum number of open connections
			db.SetMaxOpenConns(10)

			// Set the maximum number of idle connections
			db.SetMaxIdleConns(5)

			// Ping the database to verify the connection
			err = db.Ping()
			if err != nil {
				return err
			}

			// Initialize recommend_item table
			if err := utils.InitializeRecommendItems(db); err != nil {
				return err
			}

			return nil
		},
		retry.Delay(1*time.Second),
		retry.Attempts(10),
		retry.LastErrorOnly(true),
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to MySQL")
}

func GetMySQL() *sql.DB {
	return db
}

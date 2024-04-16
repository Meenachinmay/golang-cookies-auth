package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"golang-cookies/handlers"
	"golang-cookies/internal/config"
	"golang-cookies/internal/database"
	"log"
	"os"
)

func main() {
	// Initialize the Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Initialize the database
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading env file.")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in env file.")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database")
	}

	var testQuery int
	err = conn.QueryRow("SELECT 1").Scan(&testQuery)

	if err != nil {
		log.Fatal("Database connection test failed", err)
	} else {
		log.Println("Test query executed successfully. database connection verified.")
	}

	// settingUp API configuration
	apiConfig := &config.ApiConfig{
		DB:          database.New(conn),
		RedisClient: redisClient,
	}

	localApiConfig := &handlers.LocalApiConfig{
		ApiConfig: apiConfig,
	}

	// Initialize the router
	router := gin.Default()

	router.GET("/health-check", localApiConfig.HandlerReadiness)

	log.Fatal(router.Run(":8080"))
}

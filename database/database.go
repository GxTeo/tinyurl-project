package database

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"gx/learn-go/tinyurl-project/config"
	"regexp"
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// urls -> URL database structure
type Urls struct {
	Tinyurl string `gorm:"PRIMARY_KEY"`
	Longurl string
}

// InitPostgresClient -> Initialize the Postgres client
func InitPostgresClient(config *config.Config) (*gorm.DB, error) {
	var dbClient *gorm.DB
	var err error

	deadline := time.Now().Add(1 * time.Minute)
	backoff := 100 * time.Millisecond

	for time.Now().Before(deadline) {
		dbClient, err = gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			config.DBHost, config.DBPort, config.DBUser, config.DBName, config.DBPassword))

		if err == nil {
			return dbClient, nil
		}

		time.Sleep(backoff)
		backoff *= 2 // exponential backoff

		if backoff > 5*time.Second {
			backoff = 5 * time.Second // cap maximum backoff
		}
	}

	return nil, fmt.Errorf("failed to connect to postgres after 1 minute: %v", err)
}

// InitRedisClient -> Initialize the Redis client
func InitRedisClient(config *config.Config) (*redis.Client, error) {
	var client *redis.Client
	var err error

	deadline := time.Now().Add(1 * time.Minute)
	backoff := 100 * time.Millisecond

	for time.Now().Before(deadline) {
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.CACHEHost, config.CACHEPort),
			Password: config.CACHEPassword,
			DB:       config.CACHEDB(),
		})

		// Test the connection
		_, err = client.Ping().Result()
		if err == nil {
			return client, nil
		}

		client.Close() // Close the failed connection

		time.Sleep(backoff)
		backoff *= 2 // exponential backoff

		if backoff > 5*time.Second {
			backoff = 5 * time.Second // cap maximum backoff
		}
	}

	return nil, fmt.Errorf("failed to connect to redis after 1 minute: %v", err)
}

// StoreTinyURL -> Stores the urls into cache and DB
func StoreTinyURL(dbURLData Urls, longURL string, tinyURL string, dbClient *gorm.DB, redisClient *redis.Client) {
	// Store the URL in the database
	dbClient.Create(&dbURLData)
	// Store the URL in the cache
	redisClient.HSet("urls", tinyURL, longURL)
}

// GenerateHashAndInsert -> Generates a hash for the long URL and stores it in  DB
func GenerateHashAndInsert(longURL string, startIndex int, dbClient *gorm.DB, redisClient *redis.Client) string {
	byteURLData := []byte(longURL)

	// Generate a hash for the long URL using MD5
	hashedURLData := fmt.Sprintf("%x", md5.Sum(byteURLData))
	tinyURLRegex, err := regexp.Compile("[/+]")

	if err != nil {
		return "Unable to generate tiny URL"
	}
	tinyURLData := tinyURLRegex.ReplaceAllString(base64.URLEncoding.EncodeToString([]byte(hashedURLData)), "_")

	// Unable to generate tiny URL if the length of the tiny URL is less than 6
	if len(tinyURLData) < (startIndex + 6) {
		return "Unable to generate tiny URL"
	}

	tinyURL := tinyURLData[startIndex : startIndex+6] // Generate a tiny URL from the hash

	var dbURLData Urls

	dbClient.Where("tinyurl = ?", tinyURL).First(&dbURLData) // Check if the tiny URL already exists in the database

	//
	if dbURLData.Tinyurl == "" {
		fmt.Println(dbURLData, "in not found")
		go StoreTinyURL(Urls{Tinyurl: tinyURL, Longurl: longURL}, longURL, tinyURL, dbClient, redisClient)
		return tinyURL

	} else if dbURLData.Longurl == longURL && dbURLData.Tinyurl == tinyURL {
		// If the long URL and tiny URL already exist in the database, return the tiny URL
		fmt.Println(dbURLData, "in found and equal")
		return tinyURL
	} else {
		// If the tiny URL already exists in the database, generate a new tiny URL
		return GenerateHashAndInsert(longURL, startIndex+1, dbClient, redisClient)
	}

}

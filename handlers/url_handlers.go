package handlers

import (
	"gx/learn-go/tinyurl-project/database"

	"net/http"
	"strings"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

// CreateShortUrl -> Create a short URL from the long URL
func CreateShortUrl(res http.ResponseWriter, req *http.Request, dbClient *gorm.DB, redisClient *redis.Client) {
	longURL := req.URL.Query().Get("longUrl")

	if longURL == "" {
		http.Error(res, "URL parameter longUrl is missing", http.StatusBadRequest)
		return
	}

	// Ensure that the URL is valid
	if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
		longURL = "http://" + longURL
	}

	// Generate a short URL
	tinyURL := database.GenerateHashAndInsert(longURL, 0, dbClient, redisClient)

	// Return the short URL
	baseURL := "http://" + req.Host + "/"
	fullShortURL := baseURL + tinyURL
	res.Write([]byte(fullShortURL))

}

func RedirectToLongURL(res http.ResponseWriter, req *http.Request, dbClient *gorm.DB, redisClient *redis.Client) {

	// Get the tiny URL from the request
	shortCode := strings.TrimPrefix(req.URL.Path, "/")
	if shortCode == "" {
		http.Error(res, "Short URL parameter is missing", http.StatusBadRequest)
		return
	}

	// Check if the URL is present in the cache
	longURL, err := redisClient.HGet("urls", shortCode).Result()
	if err == nil && longURL != "" {
		http.Redirect(res, req, longURL, http.StatusMovedPermanently)
		return
	}

	// Else check if the URL is present in the database
	var urlData database.Urls
	dbClient.Where("tinyurl = ?", shortCode).First(&urlData)

	if urlData.Longurl == "" {
		http.Error(res, "Short URL not found", http.StatusNotFound)
		return
	}

	// Store the URL in the cache
	redisClient.HSet("urls", shortCode, urlData.Longurl)

	// Redirect to the long URL
	http.Redirect(res, req, urlData.Longurl, http.StatusMovedPermanently)
}

package router

import (
	"gx/learn-go/tinyurl-project/handlers"
	"gx/learn-go/tinyurl-project/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

func SetupRouter(dbClient *gorm.DB, redisClient *redis.Client) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORS())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to the URL shortening service")
	})
	// Using the CreateShortUrl handler
	router.POST("/shorten", func(c *gin.Context) {
		handlers.CreateShortUrl(c.Writer, c.Request, dbClient, redisClient)
	})

	// Using the  RedirectToLongURL handler
	router.GET("/:shortCode", func(c *gin.Context) {
		handlers.RedirectToLongURL(c.Writer, c.Request, dbClient, redisClient)
	})

	return router
}

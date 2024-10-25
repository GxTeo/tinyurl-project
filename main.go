package main

import (
	"gx/learn-go/tinyurl-project/config"
	"gx/learn-go/tinyurl-project/database"
	"gx/learn-go/tinyurl-project/router"
	"log"
)

func main() {

	// Load the configuration
	config, err := config.LoadConfig()

	if err != nil {
		log.Print("Error loading the configuration file ", err)
	}

	// Initalize Postgres client
	dbClient, err := database.InitPostgresClient(config)
	if err != nil {
		log.Fatalf("Error connecting to the database %v", err)

	}
	defer dbClient.Close()

	// Initialize Redis client
	redisClient, err := database.InitRedisClient(config)
	if err != nil {
		log.Fatalf("Error connecting to the cache %v", err)
	}

	// Setup the router
	router := router.SetupRouter(dbClient, redisClient)

	// Start the server
	API_PORT := ":" + config.APIPort
	router.Run(API_PORT)

}

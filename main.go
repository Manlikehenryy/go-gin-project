package main

import (
	"github.com/gin-gonic/gin"
	"github.com/manlikehenryy/go-gin-project/configs"
	"github.com/manlikehenryy/go-gin-project/database"
	"github.com/manlikehenryy/go-gin-project/helpers"
	"github.com/manlikehenryy/go-gin-project/routes"
)

func main() {
	// Initialize configuration
	configs.Init()
	helpers.Initialize()

	// Connect to the database
	database.Connect()

	// Retrieve the port from the config
	port := configs.Env.PORT
	if port == "" {
		port = "8080" // Set a default port if not specified
	}

	// Create a new Gin router
	app := gin.New()

	// Apply middleware
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	// Set up routes
	routes.Setup(app)

	// Start the server
	err := app.Run(":" + port)
	if err != nil {
		panic("Failed to start server: " + err.Error())
	}
}

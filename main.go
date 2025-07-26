package main

import (
	"Mlack/handlers"
	"Mlack/services"
	"github.com/gin-gonic/gin"
)

func main() {
	eventService := services.NewEventService() // Manages connected clients

	r := gin.Default()
	r.Static("/template", "./template") // Serves HTML/JS files
	r.LoadHTMLGlob("template/*.html")   // Allows rendering HTML

	// Routes
	r.GET("/", handlers.IndexHandler)                                            // Homepage
	r.POST("/webhook/bitbucket", handlers.BitbucketWebhookHandler(eventService)) // Bitbucket sends data here
	r.GET("/events", handlers.EventStreamHandler(eventService))                  // Browser connects here for live updates

	r.Run(":8080") // Start server on port 8080
}

package handlers

import (
	"Mlack/models"
	"Mlack/services"
	"io"

	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func BitbucketWebhookHandler(eventService *services.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verify secret token if you set one
		// token := c.GetHeader("X-Hub-Signature")
		// if !verifyToken(token) { ... }

		var payload models.BitbucketPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(400, gin.H{"error": "Invalid payload"})
			return
		}

		// Process push event
		if len(payload.Push.Changes) > 0 {
			var messages []string
			for _, change := range payload.Push.Changes {
				for _, commit := range change.Commits {
					messages = append(messages, commit.Message)
				}
			}

			event := models.CommitEvent{
				Repository: payload.Repository.Name,
				Messages:   messages,
				Committer:  payload.Actor.DisplayName,
			}

			// Broadcast to all connected clients
			eventService.Broadcast(event)
		}

		c.JSON(200, gin.H{"status": "success"})
	}
}

func EventStreamHandler(eventService *services.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set headers for SSE
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// Create new client channel
		clientChan := make(chan models.CommitEvent)
		eventService.AddClient(clientChan)

		// Close connection when client disconnects
		defer func() {
			eventService.RemoveClient(clientChan)
		}()

		// Listen to client close
		clientGone := c.Writer.CloseNotify()
		c.Stream(func(w io.Writer) bool {
			select {
			case event := <-clientChan:
				c.SSEvent("message", event)
				return true
			case <-clientGone:
				return false
			}
		})
	}
}

package controllers

import (
	"fmt"
	"io"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/kaigoh/ganjapp/models"
)

// SSE Server Side Events handler
func SSE(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)

	// Ensure entries exist in the streaming maps for this user...
	models.CreateUserStream(user.ID)

	go func() {
		for {
			select {
			case <-c.Request.Context().Done():
				models.StreamStatus[user.ID] <- true
				return
			case <-c.Done():
				models.StreamStatus[user.ID] <- true
				return
			}
		}
	}()

	c.Stream(func(w io.Writer) bool {
		for {
			select {
			case <-models.StreamStatus[user.ID]:
				c.SSEvent("end", "end")
				return false
			case msg := <-models.StreamChannel[user.ID]:
				c.Render(-1, sse.Event{
					Id:    fmt.Sprint(models.StreamMessageID[user.ID]),
					Event: msg.Type,
					Data:  msg.Data,
				})
				models.StreamMessageID[user.ID]++
				return true
			}
		}
	})
}

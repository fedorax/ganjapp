package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kaigoh/ganjapp/models"
)

// GetEvents fetches events for the logged in user
func GetEvents(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)

	limitParam := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 20
	}

	events, err := models.GetEvents(user.ID, limit)

	c.JSON(http.StatusOK, events)
}

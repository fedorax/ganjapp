package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Home handler
func Home(c *gin.Context) {
	user, _ := GetUserFromSession(c, true)
	c.HTML(http.StatusOK, "page.tmpl", gin.H{
		"user":  user,
		"title": "Dashboard",
	})
}

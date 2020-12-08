package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaigoh/ganjapp/controllers"
	"github.com/kaigoh/ganjapp/utilities"
)

// IsUserLoggedInMiddleware intercepts HTTP requests, ensuring the user is logged in
func IsUserLoggedInMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		session := sessions.Default(c)

		// Try and find the JWT token...
		jwtTokenSession := session.Get("token")
		jwtToken := ""
		if jwtTokenSession == nil {
			jwtToken = strings.TrimSpace(strings.Replace(c.Request.Header.Get("Authorization"), "Bearer ", "", 1))
		} else {
			jwtToken = jwtTokenSession.(string)
		}

		// Fetch the JWT token from the session...
		if len(jwtToken) > 0 {

			// Check the JWT
			claims, err := utilities.ParseJWT(jwtToken)
			if err != nil {
				controllers.RedirectToLoginAndAbort(c, "Session expired")
			}

			c.Set("user", claims["sub"])
			c.Next()

		} else {
			controllers.RedirectToLoginAndAbort(c, "Session expired")
		}
	}
}

// IsUserAuthenticatedMiddleware intercepts HTTP requests, ensuring that the requesting user is authenticated
func IsUserAuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		session := sessions.Default(c)

		// Try and find the JWT token...
		jwtTokenSession := session.Get("token")
		jwtToken := ""
		if jwtTokenSession == nil {
			jwtToken = strings.TrimSpace(strings.Replace(c.Request.Header.Get("Authorization"), "Bearer ", "", 1))
		} else {
			jwtToken = jwtTokenSession.(string)
		}

		// Fetch the JWT token from the session...
		if len(jwtToken) > 0 {

			// Check the JWT
			claims, err := utilities.ParseJWT(jwtToken)
			if err != nil {
				c.String(http.StatusForbidden, "403 Not Authenticated")
				c.Abort()
			}

			c.Set("user", claims["sub"])
			c.Next()

		} else {
			c.String(http.StatusForbidden, "403 Not Authenticated")
			c.Abort()
		}
	}
}

package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaigoh/ganjapp/models"
)

// Login handler
func Login(c *gin.Context) {
	message := c.DefaultQuery("message", "")
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"message": message,
		"title":   "Login",
	})
}

// GetUserFromSession trys to fetch a User struct from the HTTP session
func GetUserFromSession(c *gin.Context, full bool) (models.User, error) {
	userEmail, _ := c.Get("user")
	return models.GetUserByEmail(userEmail.(string), full)
}

// RedirectToLoginAndAbort redirects the request to the login page
func RedirectToLoginAndAbort(c *gin.Context, message string) {
	if len(message) > 0 {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/login?message="+message)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
	}
	c.Abort()
}

// Authenticate handler
func Authenticate(c *gin.Context) {
	// Initialise the session...
	session := sessions.Default(c)
	session.Clear()

	// Try and authenticate the user...
	username := c.DefaultPostForm("username", "[!!error!!]")
	password := c.DefaultPostForm("password", "[!!error!!]")

	if username == "[!!error!!]" || password == "[!!error!!]" {
		RedirectToLoginAndAbort(c, "Username and / or password not passed with request")
	}

	user, err := models.AuthenticateUserByPassword(username, password)

	// If we were unable to authenticate, redirect to the login page...
	if err != nil {
		// Log the failure...
		go models.LogSystemEvent("user-login-failed", "warning", "Failed to authenticate '"+username+"' using password authentication.", username)
		RedirectToLoginAndAbort(c, "Username and password not valid")
	} else {
		// User authenticated OK, generate a token and set it in the session...
		// Note that we issue the token so it lasts 24 hours from now
		token, err := models.IssueJWT(user, 86400)

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			c.Abort()
		}

		session.Set("token", token)
		session.Save()
		c.Set("user", username)
		// Log the failure...
		go models.LogSystemEvent("user-login-success", "info", "Authenticated '"+username+"' using password authentication.", username)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		c.Abort()
	}

}

// GetToken generates a non-expiring JWT and returns it
func GetToken(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)
	token, err := models.IssueJWT(user, -1)

	// Crap out if we can't issue the token...
	if err != nil {
		c.String(http.StatusInternalServerError, "500 Internal Server Error: Unable to issue token")
		return
	}

	// Return the token in a JSON object...
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

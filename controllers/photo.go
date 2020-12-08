package controllers

import (
	"net/http"
	"strings"

	"github.com/kaigoh/ganjapp/utilities"
	"github.com/minio/minio-go/v7"

	"github.com/kaigoh/ganjapp/models"

	"github.com/gin-gonic/gin"
)

// Photo handler
func Photo(c *gin.Context) {
	user, _ := GetUserFromSession(c, true)
	objectType := c.Param("objectType")
	objectUUID := c.Param("objectUUID")
	imageUUID := c.Param("imageUUID")
	var f *minio.Object

	switch o := strings.ToLower(objectType); o {
	case "environment":
		// Try and find a matching image...
		e, err := models.GetUserEnvironmentByUUID(user.ID, objectUUID, true)

		// Did we find the environment?
		if err != nil {
			c.AbortWithStatus(404)
		}

		i, err := e.GetImageByUUID(imageUUID)

		// Did we find the image?
		if err != nil {
			c.AbortWithStatus(404)
		}

		// Try and serve up the image...
		f, err = utilities.S3GetFile(*i.Path)

		// Did we find the image?
		if err != nil {
			c.AbortWithStatus(404)
		}

	case "tree":
		// Try and find a matching image...
		t, err := models.GetUserTreeByUUID(user.ID, objectUUID, true)

		// Did we find the environment?
		if err != nil {
			c.AbortWithStatus(404)
		}

		i, err := t.GetImageByUUID(imageUUID)

		// Did we find the image?
		if err != nil {
			c.AbortWithStatus(404)
		}

		// Try and serve up the image...
		f, err = utilities.S3GetFile(*i.Path)

		// Did we find the image?
		if err != nil {
			c.AbortWithStatus(404)
		}

	case "shroom":
		// Try and find a matching image...
		s, err := models.GetUserShroomByUUID(user.ID, objectUUID, true)

		// Did we find the environment?
		if err != nil {
			c.AbortWithStatus(404)
		}

		i, err := s.GetImageByUUID(imageUUID)

		// Did we find the image?
		if err != nil {
			c.AbortWithStatus(404)
		}

		// Try and serve up the image...
		f, err = utilities.S3GetFile(*i.Path)

		// Did we find the image?
		if err != nil {
			c.AbortWithStatus(404)
		}

	}

	// Serve the file...
	c.DataFromReader(http.StatusOK, 0, "", f, map[string]string{})

}

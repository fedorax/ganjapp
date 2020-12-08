package controllers

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kaigoh/ganjapp/models"
	"github.com/kaigoh/ganjapp/utilities"
)

// MoveShroom handler
func MoveShroom(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)
	shroomID := c.Param("id")
	environmentID := c.Param("environment")

	// Get the shroom environment...
	shroom, err := models.GetUserShroomByUUID(user.ID, shroomID, true)

	if err != nil {
		c.String(http.StatusNotFound, "404 Not Found: Shroom not found or not owned by requesting user")
		return
	}

	// Current environment ID...
	current, _ := models.GetEnvironmentByID(shroom.EnvironmentID, false)

	// Get the new environment...
	env, err := models.GetUserEnvironmentByUUID(user.ID, environmentID, true)

	if err != nil {
		c.String(http.StatusNotFound, "404 Not Found: Destination Environment not found or not owned by requesting user")
		return
	}

	shroom.EnvironmentID = env.ID
	models.DB.Save(&shroom)

	// Log everything...
	go shroom.LogEvent("shroom-moved", "info", "Shroom moved from '"+current.Name+"' to '"+env.Name+"'", shroom.UUID)
	go current.LogEvent("object-moved-out", "info", "Tree moved from '"+current.Name+"' to '"+env.Name+"'", shroom.UUID)
	go env.LogEvent("object-moved-in", "info", "Tree moved from '"+current.Name+"' to '"+env.Name+"'", shroom.UUID)

	// Publish the old environment...(we don't need to publish the destination environment, as that's already done in a gorm hook...)
	current.PublishEnvironment()

	c.String(http.StatusOK, "200 OK")

}

// UpdateShroom handler
func UpdateShroom(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)
	shroomID := c.Param("id")
	property := c.Param("property")
	value := c.DefaultPostForm("value", "[!!error!!]")

	// Did we get a value parameter with the request?
	if value == "[!!error!!]" {
		c.String(http.StatusBadRequest, "400 Bad Request: \"value\" parameter missing")
		return
	}

	shroom, err := models.GetUserShroomByUUID(user.ID, shroomID, true)

	if err != nil {
		c.String(http.StatusNotFound, "404 Not Found: Shroom not found")
		return
	}

	switch p := strings.ToLower(property); p {

	// Set the name
	case "name":
		v := strings.TrimSpace(value)
		shroom.Name = &v
		go shroom.LogEvent("shroom-updated", "info", "shroom re-named to '"+*shroom.Name+"'", *shroom.Name)

	// Set the comments
	case "comments":
		shroom.Comments = &value
		go shroom.LogEvent("shroom-updated", "info", "shroom comments set to '"+*shroom.Comments+"'", *shroom.Comments)

	// Set the breeder...
	case "breeder":
		v := strings.TrimSpace(value)
		shroom.Breeder = &v
		go shroom.LogEvent("shroom-updated", "info", "shroom breeder set to '"+*shroom.Breeder+"'", *shroom.Breeder)

	// Set the strain...
	case "strain":
		v := strings.TrimSpace(value)
		shroom.Strain = &v
		go shroom.LogEvent("shroom-updated", "info", "shroom strain set to '"+*shroom.Strain+"'", *shroom.Strain)

	// Property not found...
	default:
		c.String(http.StatusBadRequest, "400 Bad Request: '"+p+"' is not a valid shroom property")
		return
	}

	models.DB.Save(&shroom)
	c.String(http.StatusOK, "200 OK")

}

// UpdateShroomExtendedData handler
func UpdateShroomExtendedData(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)
	shroomID := c.Param("id")
	property := c.Param("property")
	value := c.DefaultPostForm("value", "[!!error!!]")

	// Did we get a value parameter with the request?
	if value == "[!!error!!]" {
		c.String(http.StatusBadRequest, "400 Bad Request: \"value\" parameter missing")
		return
	}

	shroom, err := models.GetUserShroomByUUID(user.ID, shroomID, true)

	if err != nil {
		c.String(http.StatusNotFound, "404 Not Found: Shroom not found")
		return
	}

	p := strings.TrimSpace(strings.ToLower(property))

	// Try and find and update the key...
	for _, s := range shroom.ExtendedData {
		if s.Key == p {
			s.Value = &value
			models.DB.Save(&s)

			// Log the update...
			go shroom.LogEvent("shroom-extended-status-updated", "info", "Shroom property '"+p+"' set to '"+value+"'", p+"="+value)

			c.String(http.StatusOK, "200 OK")
			return
		}
	}

	// If we've made it here, we need to create a new extended data key...
	s := models.ShroomExtendedData{ShroomID: shroom.ID, Key: p, Value: &value}
	models.DB.Create(&s)

	// Log the update...
	go shroom.LogEvent("shroom-extended-status-updated", "info", "Shroom property '"+p+"' set to '"+value+"'", p+"="+value)

	c.String(http.StatusOK, "200 OK")
	return

}

// UploadShroomImage handler
func UploadShroomImage(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)
	shroomID := c.Param("id")
	comments := c.DefaultPostForm("comments", "")
	file, _ := c.FormFile("image")
	shroom, err := models.GetUserShroomByUUID(user.ID, shroomID, true)

	if err != nil {
		c.String(http.StatusNotFound, "404 Not Found: Shroom not found")
		return
	}

	path := utilities.GetEnv("GANJAPP_TEMP_ROOT", filepath.Join(utilities.GetEnv("GANJAPP_ROOT", ""), "temp"))
	filename := filepath.Join(path, shroomID+"-"+filepath.Base(file.Filename))

	c.SaveUploadedFile(file, filename)

	clean, err := utilities.StripExif(filename)

	// Failed to strip EXIF data from the uploaded image...
	if err != nil {
		c.String(http.StatusBadRequest, "400 Bad Request: Failed to strip EXIF data from uploaded image")
		return
	}

	// Write the cleaned image back to disk...
	if err := ioutil.WriteFile(filename, clean, 0644); err != nil {
		c.String(http.StatusInternalServerError, "500 Internal Server Error: Failed to save uploaded image")
		return
	}

	// Create a new ShroomImage...
	si := models.ShroomImage{ShroomID: shroom.ID, Comments: &comments}
	models.DB.Create(&si)

	// Upload the image to S3...
	s3Name := strconv.Itoa(int(shroom.ID)) + "-" + strconv.Itoa(int(si.ID)) + "-" + filepath.Base(filename)
	s3, err := utilities.S3UploadFile(filename, s3Name)

	if err != nil || !s3 {
		c.String(http.StatusInternalServerError, "500 Internal Server Error: Failed to upload image to S3")
		utilities.FileDelete(filename)
		models.DB.Delete(&si)
		return
	}

	// Set the path to the image and save...
	si.Path = &s3Name
	models.DB.Save(&si)

	// Send the environment to the event stream...
	models.PublishEnvironmentByID(shroom.EnvironmentID)

	// Log the upload...
	go shroom.LogEvent("shroom-image-uploaded", "info", "Shroom image uploaded", *si.Path)

	c.JSON(http.StatusOK, si)

}

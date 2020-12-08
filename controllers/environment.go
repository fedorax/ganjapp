package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kaigoh/ganjapp/models"
	"github.com/kaigoh/ganjapp/utilities"
)

// CreateEnvironment handler
func CreateEnvironment(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)

	type NewEnvironment struct {
		Name string `form:"name"`
	}

	var newEnvironment NewEnvironment

	if c.ShouldBind(&newEnvironment) == nil {

		// Create the new environment...
		e := models.Environment{Name: newEnvironment.Name, UserID: user.ID}

		// Write the environment to the database...
		models.DB.Create(&e)

		// Log the event...
		// event string, severity string, message string, data string
		j, _ := json.Marshal(e)
		go e.LogEvent("environment-created", "info", "Environment created named '"+e.Name+"'", string(j))

		c.String(http.StatusOK, "200 Created Environment OK")
		return

	}

	c.String(http.StatusBadRequest, "400 Bad Request: 'name' must be passed as part of the request")

}

// GetEnvironments handler
func GetEnvironments(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)
	environments, _ := models.GetUserEnvironments(user.ID, true)
	c.JSON(http.StatusOK, environments)
}

// UpdateEnvironment handler
func UpdateEnvironment(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)
	environmentID := c.Param("id")
	property := c.Param("property")
	value := c.DefaultPostForm("value", "[!!error!!]")

	// Did we get a value parameter with the request?
	if value == "[!!error!!]" {
		c.String(http.StatusBadRequest, "400 Bad Request: \"value\" parameter missing")
		return
	}

	environment, err := models.GetUserEnvironmentByUUID(user.ID, environmentID, false)

	if err != nil {
		c.String(http.StatusNotFound, "404 Not Found: Environment not found")
		return
	}

	switch p := strings.TrimSpace(strings.ToLower(property)); p {

	// Set the name
	case "name":
		environment.Name = strings.TrimSpace(value)
		go environment.LogEvent("environment-updated", "info", "Environment re-named to '"+environment.Name+"'", environment.Name)

	// Set the comments
	case "comments":
		environment.Comments = &value
		go environment.LogEvent("environment-updated", "info", "Environment comments set to '"+environment.Name+"'", value)

	// Set the temperature...
	case "temperature":
		t, e := strconv.ParseFloat(value, 32)
		if e != nil {
			c.String(http.StatusBadRequest, "400 Bad Request: Unable to parse temperature from passed value")
			return
		}
		v := float32(t)
		environment.Status.Temperature = &v
		go environment.LogEvent("environment-status-updated", "info", "Environment temperature currently "+fmt.Sprintf("%.2f", v)+"°C", fmt.Sprintf("%.2f", v))

	// Set the humidity...
	case "humidity":
		t, e := strconv.ParseFloat(value, 32)
		if e != nil {
			c.String(http.StatusBadRequest, "400 Bad Request: Unable to parse humidity from passed value")
			return
		}
		v := float32(t)
		environment.Status.Humidity = &v
		go environment.LogEvent("environment-status-updated", "info", "Environment humidity currently "+fmt.Sprintf("%.2f", v)+"%", fmt.Sprintf("%.2f", v))

	// Set the lighting status...
	case "lighting":
		t, e := strconv.ParseBool(value)
		if e != nil {
			c.String(http.StatusBadRequest, "400 Bad Request: Unable to set lighting status from passed value")
			return
		}
		environment.Status.LightsOn = t
		// Hacky way of making a ternary operator...
		text := (map[bool]string{true: "on", false: "off"})[t]
		go environment.LogEvent("environment-status-updated", "info", "Environment lighting currently "+text, text)

	// Property not found...
	default:
		c.String(http.StatusBadRequest, "400 Bad Request: '"+p+"' is not a valid Environment property")
		return
	}

	models.DB.Save(&environment)
	c.String(http.StatusOK, "200 OK")

}

// UpdateEnvironmentExtendedData handler
func UpdateEnvironmentExtendedData(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)
	environmentID := c.Param("id")
	property := c.Param("property")
	value := c.DefaultPostForm("value", "[!!error!!]")

	// Did we get a value parameter with the request?
	if value == "[!!error!!]" {
		c.String(http.StatusBadRequest, "400 Bad Request: \"value\" parameter missing")
		return
	}

	environment, err := models.GetUserEnvironmentByUUID(user.ID, environmentID, false)

	if err != nil {
		c.String(http.StatusNotFound, "404 Not Found: Environment not found")
		return
	}

	p := strings.TrimSpace(strings.ToLower(property))

	// Try and find and update the key...
	for _, e := range environment.ExtendedData {
		if e.Key == p {
			e.Value = &value
			models.DB.Save(&e)

			// Log the update...
			go environment.LogEvent("environment-extended-status-updated", "info", "Environment property '"+p+"' set to '"+value+"'", p+"="+value)

			c.String(http.StatusOK, "200 OK")
			return
		}
	}

	// If we've made it here, we need to create a new extended data key...
	e := models.EnvironmentExtendedData{EnvironmentID: environment.ID, Key: p, Value: &value}
	models.DB.Create(&e)

	// Log the update...
	go environment.LogEvent("environment-extended-status-updated", "info", "Environment property '"+p+"' set to '"+value+"'", p+"="+value)

	c.String(http.StatusOK, "200 OK")
	return

}

// UploadEnvironmentImage handler
func UploadEnvironmentImage(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)
	environmentID := c.Param("id")
	comments := c.DefaultPostForm("comments", "")
	file, _ := c.FormFile("image")
	environment, err := models.GetUserEnvironmentByUUID(user.ID, environmentID, false)

	if err != nil {
		c.String(http.StatusNotFound, "404 Not Found: Environment not found")
		return
	}

	path := utilities.GetEnv("GANJAPP_TEMP_ROOT", filepath.Join(utilities.GetEnv("GANJAPP_ROOT", ""), "temp"))
	filename := filepath.Join(path, environmentID+"-"+filepath.Base(file.Filename))

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

	// Create a new EnvironmentImage...
	ei := models.EnvironmentImage{EnvironmentID: environment.ID, Comments: &comments}
	models.DB.Create(&ei)

	// Upload the image to S3...
	s3Name := strconv.Itoa(int(environment.ID)) + "-" + strconv.Itoa(int(ei.ID)) + "-" + filepath.Base(filename)
	s3, err := utilities.S3UploadFile(filename, s3Name)

	if err != nil || !s3 {
		c.String(http.StatusInternalServerError, "500 Internal Server Error: Failed to upload image to S3")
		utilities.FileDelete(filename)
		models.DB.Delete(&ei)
		return
	}

	// Set the path to the image and save...
	ei.Path = &s3Name
	models.DB.Save(&ei)

	// Send the environment to the event stream...
	environment.PublishEnvironment()

	// Log the upload...
	go environment.LogEvent("environment-image-uploaded", "info", "Environment image uploaded", *ei.Path)

	c.JSON(http.StatusOK, ei)

}

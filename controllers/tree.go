package controllers

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kaigoh/ganjapp/models"
	"github.com/kaigoh/ganjapp/utilities"
)

// MoveTree handler
func MoveTree(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)
	treeID := c.Param("id")
	environmentID := c.Param("environment")

	// Get the tree environment...
	tree, err := models.GetUserTreeByUUID(user.ID, treeID, true)

	if err != nil {
		c.String(http.StatusNotFound, "404 Not Found: Tree not found or not owned by requesting user")
		return
	}

	// Current environment ID...
	current, _ := models.GetEnvironmentByID(tree.EnvironmentID, false)

	// Get the new environment...
	env, err := models.GetUserEnvironmentByUUID(user.ID, environmentID, true)

	if err != nil {
		c.String(http.StatusNotFound, "404 Not Found: Destination Environment not found or not owned by requesting user")
		return
	}

	tree.EnvironmentID = env.ID
	models.DB.Save(&tree)

	// Log everything...
	go tree.LogEvent("tree-moved", "info", "Tree moved from '"+current.Name+"' to '"+env.Name+"'", tree.UUID)
	go current.LogEvent("object-moved-out", "info", "Tree moved from '"+current.Name+"' to '"+env.Name+"'", tree.UUID)
	go env.LogEvent("object-moved-in", "info", "Tree moved from '"+current.Name+"' to '"+env.Name+"'", tree.UUID)

	// Publish the old environment...(we don't need to publish the destination environment, as that's already done in a gorm hook...)
	current.PublishEnvironment()

	c.String(http.StatusOK, "200 OK")

}

// UpdateTree handler
func UpdateTree(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)
	treeID := c.Param("id")
	property := c.Param("property")
	value := c.DefaultPostForm("value", "[!!error!!]")

	// Did we get a value parameter with the request?
	if value == "[!!error!!]" {
		c.String(http.StatusBadRequest, "400 Bad Request: \"value\" parameter missing")
		return
	}

	tree, err := models.GetUserTreeByUUID(user.ID, treeID, true)

	if err != nil {
		c.String(http.StatusNotFound, "404 Not Found: Tree not found")
		return
	}

	switch p := strings.ToLower(property); p {

	// Set the name
	case "name":
		v := strings.TrimSpace(value)
		tree.Name = &v
		go tree.LogEvent("tree-updated", "info", "Tree re-named to '"+*tree.Name+"'", *tree.Name)

	// Set the comments
	case "comments":
		tree.Comments = &value
		go tree.LogEvent("tree-updated", "info", "Tree comments set to '"+*tree.Comments+"'", *tree.Comments)

	// Set the breeder...
	case "breeder":
		v := strings.TrimSpace(value)
		tree.Breeder = &v
		go tree.LogEvent("tree-updated", "info", "Tree breeder set to '"+*tree.Breeder+"'", *tree.Breeder)

	// Set the strain...
	case "strain":
		v := strings.TrimSpace(value)
		tree.Strain = &v
		go tree.LogEvent("tree-updated", "info", "Tree strain set to '"+*tree.Strain+"'", *tree.Strain)

		// Set the strain...
	case "germinationdate":
		t, err := time.Parse(time.RFC3339, value)
		if err != nil {
			c.String(http.StatusBadRequest, "400 Bad Request: Failed to parse '"+value+"' into a valid date time - Should be in YYYY-MM-DDTHH:MM:SSZTZ format")
			return
		}
		tree.GerminationDate = &t
		go tree.LogEvent("tree-updated", "info", "Tree germination date set to '"+value+"'", value)

		// Set the strain...
	case "vegetativedate":
		t, err := time.Parse(time.RFC3339, value)
		if err != nil {
			c.String(http.StatusBadRequest, "400 Bad Request: Failed to parse '"+value+"' into a valid date time - Should be in YYYY-MM-DDTHH:MM:SSZTZ format")
			return
		}
		tree.VegetativeDate = &t
		go tree.LogEvent("tree-updated", "info", "Tree vegetative date set to '"+value+"'", value)

		// Set the strain...
	case "floweringdate":
		t, err := time.Parse(time.RFC3339, value)
		if err != nil {
			c.String(http.StatusBadRequest, "400 Bad Request: Failed to parse '"+value+"' into a valid date time - Should be in YYYY-MM-DDTHH:MM:SSZTZ format")
			return
		}
		tree.FloweringDate = &t
		go tree.LogEvent("tree-updated", "info", "Tree flowering date set to '"+value+"'", value)

		// Set the strain...
	case "dryingdate":
		t, err := time.Parse(time.RFC3339, value)
		if err != nil {
			c.String(http.StatusBadRequest, "400 Bad Request: Failed to parse '"+value+"' into a valid date time - Should be in YYYY-MM-DDTHH:MM:SSZTZ format")
			return
		}
		tree.DryingDate = &t
		go tree.LogEvent("tree-updated", "info", "Tree drying date set to '"+value+"'", value)

		// Set the strain...
	case "curingdate":
		t, err := time.Parse(time.RFC3339, value)
		if err != nil {
			c.String(http.StatusBadRequest, "400 Bad Request: Failed to parse '"+value+"' into a valid date time - Should be in YYYY-MM-DDTHH:MM:SSZTZ format")
			return
		}
		tree.CuringDate = &t
		go tree.LogEvent("tree-updated", "info", "Tree curing date set to '"+value+"'", value)

	// Property not found...
	default:
		c.String(http.StatusBadRequest, "400 Bad Request: '"+p+"' is not a valid Tree property")
		return
	}

	models.DB.Save(&tree)
	c.String(http.StatusOK, "200 OK")

}

// UpdateTreeExtendedData handler
func UpdateTreeExtendedData(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)
	treeID := c.Param("id")
	property := c.Param("property")
	value := c.DefaultPostForm("value", "[!!error!!]")

	// Did we get a value parameter with the request?
	if value == "[!!error!!]" {
		c.String(http.StatusBadRequest, "400 Bad Request: \"value\" parameter missing")
		return
	}

	tree, err := models.GetUserTreeByUUID(user.ID, treeID, true)

	if err != nil {
		c.String(http.StatusNotFound, "404 Not Found: Tree not found")
		return
	}

	p := strings.TrimSpace(strings.ToLower(property))

	// Try and find and update the key...
	for _, e := range tree.ExtendedData {
		if e.Key == p {
			e.Value = &value
			models.DB.Save(&e)

			// Log the update...
			go tree.LogEvent("tree-extended-status-updated", "info", "Tree property '"+p+"' set to '"+value+"'", p+"="+value)

			c.String(http.StatusOK, "200 OK")
			return
		}
	}

	// If we've made it here, we need to create a new extended data key...
	e := models.TreeExtendedData{TreeID: tree.ID, Key: p, Value: &value}
	models.DB.Create(&e)

	// Log the update...
	go tree.LogEvent("tree-extended-status-updated", "info", "Tree property '"+p+"' set to '"+value+"'", p+"="+value)

	c.String(http.StatusOK, "200 OK")
	return

}

// UploadTreeImage handler
func UploadTreeImage(c *gin.Context) {
	user, _ := GetUserFromSession(c, false)
	treeID := c.Param("id")
	comments := c.DefaultPostForm("comments", "")
	file, _ := c.FormFile("image")
	tree, err := models.GetUserTreeByUUID(user.ID, treeID, true)

	if err != nil {
		c.String(http.StatusNotFound, "404 Not Found: Tree not found")
		return
	}

	path := utilities.GetEnv("GANJAPP_TEMP_ROOT", filepath.Join(utilities.GetEnv("GANJAPP_ROOT", ""), "temp"))
	filename := filepath.Join(path, treeID+"-"+filepath.Base(file.Filename))

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

	// Create a new TreeImage...
	ti := models.TreeImage{TreeID: tree.ID, Comments: &comments}
	models.DB.Create(&ti)

	// Upload the image to S3...
	s3Name := strconv.Itoa(int(tree.ID)) + "-" + strconv.Itoa(int(ti.ID)) + "-" + filepath.Base(filename)
	s3, err := utilities.S3UploadFile(filename, s3Name)

	if err != nil || !s3 {
		c.String(http.StatusInternalServerError, "500 Internal Server Error: Failed to upload image to S3")
		utilities.FileDelete(filename)
		models.DB.Delete(&ti)
		return
	}

	// Set the path to the image and save...
	ti.Path = &s3Name
	models.DB.Save(&ti)

	// Send the environment to the event stream...
	models.PublishEnvironmentByID(tree.EnvironmentID)

	// Log the upload...
	go tree.LogEvent("tree-image-uploaded", "info", "Tree image uploaded", *ti.Path)

	c.JSON(http.StatusOK, ti)

}

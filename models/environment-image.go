package models

import (
	"github.com/kaigoh/ganjapp/utilities"
	"gorm.io/gorm"
)

// EnvironmentImage model
type EnvironmentImage struct {
	gorm.Model
	EnvironmentID uint   `gorm:"index"`
	UUID          string `gorm:"unique;uniqueIndex"`
	Path          *string
	Comments      *string
}

// BeforeSave hook adds a UUID to the model
func (e *EnvironmentImage) BeforeSave(tx *gorm.DB) (err error) {
	if e.UUID == "" {
		e.UUID = utilities.UUID()
	}
	return
}

// AfterDelete hook removes object from S3 storage
func (e *EnvironmentImage) AfterDelete(tx *gorm.DB) (err error) {
	if *e.Path != "" {
		utilities.S3DeleteFile(*e.Path)
	}
	return
}

// GetEnvironmentImageByUUID attempts to return an image by UUID
func GetEnvironmentImageByUUID(environmentID uint, imageUUID string) (EnvironmentImage, error) {
	var image EnvironmentImage
	e := DB.Where("environment_id = ?", environmentID).Where("uuid = ?", imageUUID).First(&image)
	return image, e.Error
}

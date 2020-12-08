package models

import (
	"github.com/kaigoh/ganjapp/utilities"
	"gorm.io/gorm"
)

// ShroomImage model
type ShroomImage struct {
	gorm.Model
	ShroomID uint   `gorm:"index"`
	UUID     string `gorm:"unique;uniqueIndex"`
	Path     *string
	Comments *string
}

// BeforeSave hook adds a UUID to the model
func (s *ShroomImage) BeforeSave(tx *gorm.DB) (err error) {
	if s.UUID == "" {
		s.UUID = utilities.UUID()
	}
	return
}

// AfterDelete hook removes object from S3 storage
func (s *ShroomImage) AfterDelete(tx *gorm.DB) (err error) {
	if *s.Path != "" {
		utilities.S3DeleteFile(*s.Path)
	}
	return
}

// GetShroomImageByUUID attempts to return an image by UUID
func GetShroomImageByUUID(shroomID uint, imageUUID string) (ShroomImage, error) {
	var image ShroomImage
	e := DB.Where("shroom_id = ?", shroomID).Where("uuid = ?", imageUUID).First(&image)
	return image, e.Error
}

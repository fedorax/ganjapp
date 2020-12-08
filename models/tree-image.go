package models

import (
	"github.com/kaigoh/ganjapp/utilities"
	"gorm.io/gorm"
)

// TreeImage model
type TreeImage struct {
	gorm.Model
	TreeID   uint   `gorm:"index"`
	UUID     string `gorm:"unique;uniqueIndex"`
	Path     *string
	Comments *string
}

// BeforeSave hook adds a UUID to the model
func (t *TreeImage) BeforeSave(tx *gorm.DB) (err error) {
	if t.UUID == "" {
		t.UUID = utilities.UUID()
	}
	return
}

// AfterDelete hook removes object from S3 storage
func (t *TreeImage) AfterDelete(tx *gorm.DB) (err error) {
	if *t.Path != "" {
		utilities.S3DeleteFile(*t.Path)
	}
	return
}

// GetTreeImageByUUID attempts to return an image by UUID
func GetTreeImageByUUID(treeID uint, imageUUID string) (TreeImage, error) {
	var image TreeImage
	e := DB.Where("tree_id = ?", treeID).Where("uuid = ?", imageUUID).First(&image)
	return image, e.Error
}

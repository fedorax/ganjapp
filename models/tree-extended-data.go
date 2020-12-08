package models

import "gorm.io/gorm"

// TreeExtendedData struct for allowing custom data against tree models
type TreeExtendedData struct {
	gorm.Model
	TreeID uint `gorm:"index"`
	Key    string
	Value  *string
}

// AfterSave hook
func (t *TreeExtendedData) AfterSave(tx *gorm.DB) (err error) {
	// Send the updated environment to the event stream...
	tree, e := GetTreeByID(t.TreeID, false)
	if e != nil {
		return
	}
	PublishEnvironmentByID(tree.EnvironmentID)
	return
}

// AfterUpdate hook
func (t *TreeExtendedData) AfterUpdate(tx *gorm.DB) (err error) {
	// Send the updated environment to the event stream...
	tree, e := GetTreeByID(t.TreeID, false)
	if e != nil {
		return
	}
	PublishEnvironmentByID(tree.EnvironmentID)
	return
}

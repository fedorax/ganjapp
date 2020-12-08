package models

import "gorm.io/gorm"

// EnvironmentExtendedData struct for allowing custom data against environment models
type EnvironmentExtendedData struct {
	gorm.Model
	EnvironmentID uint `gorm:"index"`
	Key           string
	Value         *string
}

// AfterSave hook
func (e *EnvironmentExtendedData) AfterSave(tx *gorm.DB) (err error) {
	// Send the updated environment to the event stream...
	PublishEnvironmentByID(e.EnvironmentID)
	return
}

// AfterUpdate hook
func (e *EnvironmentExtendedData) AfterUpdate(tx *gorm.DB) (err error) {
	// Send the updated environment to the event stream...
	PublishEnvironmentByID(e.EnvironmentID)
	return
}

package models

import "gorm.io/gorm"

// ShroomExtendedData struct for allowing custom data against shroom models
type ShroomExtendedData struct {
	gorm.Model
	ShroomID uint `gorm:"index"`
	Key      string
	Value    *string
}

// AfterSave hook
func (s *ShroomExtendedData) AfterSave(tx *gorm.DB) (err error) {
	// Send the updated environment to the event stream...
	shroom, e := GetShroomByID(s.ShroomID, false)
	if e != nil {
		return
	}
	PublishEnvironmentByID(shroom.EnvironmentID)
	return
}

// AfterUpdate hook
func (s *ShroomExtendedData) AfterUpdate(tx *gorm.DB) (err error) {
	// Send the updated environment to the event stream...
	shroom, e := GetShroomByID(s.ShroomID, false)
	if e != nil {
		return
	}
	PublishEnvironmentByID(shroom.EnvironmentID)
	return
}

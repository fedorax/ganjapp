package models

import (
	"github.com/kaigoh/ganjapp/utilities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Shroom model
type Shroom struct {
	gorm.Model
	EnvironmentID uint   `gorm:"index"`
	UUID          string `gorm:"unique;uniqueIndex"`
	Name          *string
	Comments      *string
	Strain        *string
	Breeder       *string
	Images        []ShroomImage
	ExtendedData  []ShroomExtendedData
	Events        []Event `gorm:"-"`
}

// BeforeSave hook adds a UUID to the model
func (s *Shroom) BeforeSave(tx *gorm.DB) (err error) {
	if s.UUID == "" {
		s.UUID = utilities.UUID()
	}
	return
}

// AfterSave hook
func (s *Shroom) AfterSave(tx *gorm.DB) (err error) {
	// Send the updated environment to the event stream...
	PublishEnvironmentByID(s.EnvironmentID)
	return
}

// AfterUpdate hook
func (s *Shroom) AfterUpdate(tx *gorm.DB) (err error) {
	// Send the updated environment to the event stream...
	PublishEnvironmentByID(s.EnvironmentID)
	return
}

// AfterFind hook
func (s *Shroom) AfterFind(tx *gorm.DB) (err error) {
	events, _ := s.GetEvents(20)
	s.Events = events
	return
}

// GetEvents loads the events for this shroom
func (s *Shroom) GetEvents(limit int) ([]Event, error) {
	return GetObjectEvents("shroom", s.ID, limit)
}

// LogEvent captures an event for this shroom
func (s *Shroom) LogEvent(event string, severity string, message string, data string) error {
	e, _ := GetEnvironmentByID(s.EnvironmentID, false)
	return LogObjectEvent(e.UserID, "shroom", s.ID, event, severity, message, data)
}

// GetImageByUUID attempts to return an image belonging to this object by UUID...
func (s *Shroom) GetImageByUUID(uuid string) (ShroomImage, error) {
	return GetShroomImageByUUID(s.ID, uuid)
}

// GetShroomByID attempts to return a shroom by ID
func GetShroomByID(id uint, full bool) (Shroom, error) {
	var shroom Shroom
	if full {
		s := DB.Preload(clause.Associations).Where("id = ?", id).First(&shroom)
		return shroom, s.Error
	}
	s := DB.Preload(clause.Associations).Where("id = ?", id).First(&shroom)
	return shroom, s.Error
}

// GetUserShroomByUUID attempts to return shroom owned by a user with the passed UUID
func GetUserShroomByUUID(id uint, uuid string, full bool) (Shroom, error) {
	var shroom Shroom

	if full {
		s := DB.Joins("Environment").Preload(clause.Associations).Where("environments.user_id = ?", id).Where("shrooms.uuid = ?", uuid).First(&shroom)
		return shroom, s.Error
	}

	s := DB.Preload(clause.Associations).Joins("Environment").Where("environments.user_id = ?", id).Where("shrooms.uuid = ?", uuid).First(&shroom)
	return shroom, s.Error
}

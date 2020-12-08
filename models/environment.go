package models

import (
	"encoding/json"

	"github.com/kaigoh/ganjapp/utilities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Environment model
type Environment struct {
	gorm.Model
	UserID       uint `gorm:"index"`
	Name         string
	UUID         string `gorm:"unique;uniqueIndex"`
	Comments     *string
	Status       EnvironmentStatus `gorm:"embedded"`
	Images       []EnvironmentImage
	Trees        []Tree
	Shrooms      []Shroom
	ExtendedData []EnvironmentExtendedData
	Events       []Event `gorm:"-"`
}

// BeforeSave hook adds a UUID to the model
func (e *Environment) BeforeSave(tx *gorm.DB) (err error) {
	if e.UUID == "" {
		e.UUID = utilities.UUID()
	}
	return
}

// AfterSave hook
func (e *Environment) AfterSave(tx *gorm.DB) (err error) {
	// Send the updated environment to the event stream...
	e.PublishEnvironment()
	return
}

// AfterUpdate hook
func (e *Environment) AfterUpdate(tx *gorm.DB) (err error) {
	// Send the updated environment to the event stream...
	e.PublishEnvironment()
	return
}

// AfterFind hook
func (e *Environment) AfterFind(tx *gorm.DB) (err error) {
	events, _ := e.GetEvents(20)
	e.Events = events
	return
}

// GetEvents loads the events for this environment
func (e *Environment) GetEvents(limit int) ([]Event, error) {
	return GetObjectEvents("environment", e.ID, limit)
}

// GetImageByUUID attempts to return an image belonging to this object by UUID...
func (e *Environment) GetImageByUUID(uuid string) (EnvironmentImage, error) {
	return GetEnvironmentImageByUUID(e.ID, uuid)
}

// LogEvent captures an event for this environment
func (e *Environment) LogEvent(event string, severity string, message string, data string) error {
	return LogObjectEvent(e.UserID, "environment", e.ID, event, severity, message, data)
}

// GetUserEnvironments attempts to return environments owned by a user
func GetUserEnvironments(id uint, full bool) ([]Environment, error) {
	var env []Environment

	if full {
		e := DB.Preload("Trees.Images").Preload("Trees.ExtendedData").Preload("Shrooms.Images").Preload("Shrooms.ExtendedData").Preload(clause.Associations).Where("user_id = ?", id).Order("name").Find(&env)
		return env, e.Error
	}

	e := DB.Preload(clause.Associations).Where("user_id = ?", id).Order("name").Find(&env)
	return env, e.Error
}

// GetEnvironmentByID attempts to return an environment by its ID
func GetEnvironmentByID(id uint, full bool) (Environment, error) {
	var env Environment
	if full {
		e := DB.Preload("Trees.Images").Preload("Trees.ExtendedData").Preload("Shrooms.Images").Preload("Shrooms.ExtendedData").Preload(clause.Associations).Where("id = ?", id).First(&env)
		return env, e.Error
	}
	e := DB.Preload(clause.Associations).Where("id = ?", id).First(&env)
	return env, e.Error
}

// GetUserEnvironmentByUUID attempts to return environment owned by a user with the passed UUID
func GetUserEnvironmentByUUID(id uint, uuid string, full bool) (Environment, error) {
	var env Environment

	if full {
		e := DB.Preload("Trees.Images").Preload("Trees.ExtendedData").Preload("Shrooms.Images").Preload("Shrooms.ExtendedData").Preload(clause.Associations).Where("user_id = ?", id).Where("uuid = ?", uuid).First(&env)
		return env, e.Error
	}

	e := DB.Preload(clause.Associations).Where("user_id = ?", id).Where("uuid = ?", uuid).First(&env)
	return env, e.Error
}

// PublishEnvironment sends the passed environment to the event stream...
func (e *Environment) PublishEnvironment() bool {
	bytes, err := json.Marshal(e)
	if err == nil {
		CreateUserStream(e.UserID)
		sse := SSEEvent{Type: "environment-update", Data: string(bytes)}
		StreamChannel[e.UserID] <- sse
		return true
	}
	return false
}

// PublishEnvironmentByID attempts to fetch an environment by ID, then sends it to the event stream...
func PublishEnvironmentByID(id uint) bool {
	e, err := GetEnvironmentByID(id, true)
	if err != nil {
		return false
	}
	return e.PublishEnvironment()
}

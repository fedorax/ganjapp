package models

import (
	"time"

	"github.com/kaigoh/ganjapp/utilities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Tree model
type Tree struct {
	gorm.Model
	EnvironmentID   uint   `gorm:"index"`
	UUID            string `gorm:"unique;uniqueIndex"`
	Name            *string
	Comments        *string
	Strain          *string
	Breeder         *string
	GerminationDate *time.Time
	VegetativeDate  *time.Time
	FloweringDate   *time.Time
	DryingDate      *time.Time
	CuringDate      *time.Time
	Images          []TreeImage
	ExtendedData    []TreeExtendedData
	Events          []Event `gorm:"-"`
}

// BeforeSave hook adds a UUID to the model
func (t *Tree) BeforeSave(tx *gorm.DB) (err error) {
	if t.UUID == "" {
		t.UUID = utilities.UUID()
	}
	return
}

// AfterSave hook
func (t *Tree) AfterSave(tx *gorm.DB) (err error) {
	// Send the updated environment to the event stream...
	PublishEnvironmentByID(t.EnvironmentID)
	return
}

// AfterUpdate hook
func (t *Tree) AfterUpdate(tx *gorm.DB) (err error) {
	// Send the updated environment to the event stream...
	PublishEnvironmentByID(t.EnvironmentID)
	return
}

// AfterFind hook
func (t *Tree) AfterFind(tx *gorm.DB) (err error) {
	events, _ := t.GetEvents(20)
	t.Events = events
	return
}

// GetEvents loads the events for this tree
func (t *Tree) GetEvents(limit int) ([]Event, error) {
	return GetObjectEvents("tree", t.ID, limit)
}

// LogEvent captures an event for this tree
func (t *Tree) LogEvent(event string, severity string, message string, data string) error {
	e, _ := GetEnvironmentByID(t.EnvironmentID, false)
	return LogObjectEvent(e.UserID, "tree", t.ID, event, severity, message, data)
}

// GetImageByUUID attempts to return an image belonging to this object by UUID...
func (t *Tree) GetImageByUUID(uuid string) (TreeImage, error) {
	return GetTreeImageByUUID(t.ID, uuid)
}

// GetTreeByID attempts to return a tree by ID
func GetTreeByID(id uint, full bool) (Tree, error) {
	var tree Tree
	if full {
		t := DB.Preload(clause.Associations).Where("id = ?", id).First(&tree)
		return tree, t.Error
	}
	t := DB.Preload(clause.Associations).Where("id = ?", id).First(&tree)
	return tree, t.Error
}

// GetUserTreeByUUID attempts to return tree owned by a user with the passed UUID
func GetUserTreeByUUID(id uint, uuid string, full bool) (Tree, error) {
	var tree Tree

	if full {
		t := DB.Joins("Environment").Preload(clause.Associations).Where("environments.user_id = ?", id).Where("trees.uuid = ?", uuid).First(&tree)
		return tree, t.Error
	}

	t := DB.Preload(clause.Associations).Joins("Environment").Where("environments.user_id = ?", id).Where("trees.uuid = ?", uuid).First(&tree)
	return tree, t.Error
}

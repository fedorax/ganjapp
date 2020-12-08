package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

// Event stuct for logging
type Event struct {
	gorm.Model
	UserID     uint   `gorm:"index:event"`
	ObjectType string `gorm:"index:event"`
	ObjectID   uint   `gorm:"index:event"`
	Event      string `gorm:"index:event"`
	Severity   string
	Message    *string
	Data       *string
}

// GetEvents gets events for the passed user, sorted in most recent order
func GetEvents(userID uint, limit int) ([]Event, error) {
	var events []Event
	q := DB.Where("user_id = ?", userID).Or("user_id = 0").Order("created_at desc").Find(&events).Limit(limit)
	return events, q.Error
}

// GetSystemEvents attempts to return an array of system events
func GetSystemEvents(limit int) ([]Event, error) {
	return GetObjectEvents("system", 0, limit)
}

// GetObjectEvents attempts to return an array of Events
func GetObjectEvents(objectType string, objectID uint, limit int) ([]Event, error) {
	var events []Event
	q := DB.Where("object_type = ?", objectType).Where("object_id = ?", objectID).Order("created_at desc").Find(&events).Limit(limit)
	return events, q.Error
}

// GetSystemEventsByType fetches system events, filtered by the passed event
func GetSystemEventsByType(event string, limit int) ([]Event, error) {
	return GetObjectEventsByType("system", 0, event, limit)
}

// GetObjectEventsByType attempts to return an array of Events
func GetObjectEventsByType(objectType string, objectID uint, event string, limit int) ([]Event, error) {
	var events []Event
	q := DB.Where("object_type = ?", objectType).Where("object_id = ?", objectID).Where("event = ?", event).Order("created_at desc").Find(&events).Limit(limit)
	return events, q.Error
}

// LogSystemEvent writes a system event to the database
func LogSystemEvent(event string, severity string, message string, data string) error {
	return LogObjectEvent(0, "system", 0, event, severity, message, data)
}

// LogObjectEvent writes an event to the database
func LogObjectEvent(userID uint, objectType string, objectID uint, event string, severity string, message string, data string) error {
	e := Event{Event: event, Severity: severity, Message: &message, Data: &data, UserID: userID, ObjectType: objectType, ObjectID: objectID}
	q := DB.Create(&e)
	e.PublishEvent()
	return q.Error
}

// PublishEvent sends the event to the event stream...
func (e *Event) PublishEvent() bool {
	bytes, err := json.Marshal(e)
	if err == nil {
		CreateUserStream(e.UserID)
		sse := SSEEvent{Type: "event", Data: string(bytes)}
		StreamChannel[e.UserID] <- sse
		return true
	}
	return false
}

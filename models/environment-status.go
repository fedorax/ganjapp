package models

// EnvironmentStatus holds data about the status of an environment
type EnvironmentStatus struct {
	Temperature *float32
	Humidity    *float32
	LightsOn    bool `gorm:"default=false"`
}

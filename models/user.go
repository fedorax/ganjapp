package models

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kaigoh/ganjapp/utilities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// User model
type User struct {
	gorm.Model
	Email        string `gorm:"unique;uniqueIndex"`
	Password     string
	APIKey       string `gorm:"unique;uniqueIndex"`
	IsAdmin      bool   `gorm:"default=false;index"`
	Environments []Environment
}

// BeforeSave hook adds a UUID as an API key to the model
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.APIKey == "" {
		u.APIKey = utilities.UUID()
	}
	return
}

// GetEvents loads the events for this user
func (u *User) GetEvents(limit int) ([]Event, error) {
	return GetObjectEvents("user", u.ID, limit)
}

// LogEvent captures an event for this user
func (u *User) LogEvent(event string, severity string, message string, data string) error {
	return LogObjectEvent(u.ID, "user", u.ID, event, severity, message, data)
}

// HaveAdminAccount returns true if at least one admin account exists on the system, false otherwise
func HaveAdminAccount() bool {
	var user User
	e := DB.First(&user, "is_admin = TRUE")

	if e.Error != nil {
		return false
	}

	return true

}

// CreateInitialAdminAccount creates an initial administrator account upon first boot of the application...`
func CreateInitialAdminAccount() bool {
	te := utilities.GetEnv("GANJAPP_ADMIN_EMAIL", "admin@ganjapp.local")
	tp := utilities.GetEnv("GANJAPP_ADMIN_PASSWORD", utilities.GetRandomString(20))
	user := User{Email: te, Password: PasswordHash(tp), IsAdmin: true}
	e := DB.Create(&user)

	if e.Error != nil {
		return false
	}

	// Log the credentials we generated...
	log.Println("Warning: Initial administrator account has been created - " + te + " / " + tp)
	go LogSystemEvent("admin-account-created", "warning", "An administrator account has been created for '"+te+"'. Please note that the password will have been written to STDOUT", te)

	return true

}

// CreateUser attempts to create a new user

// GetUser attempts to return a user by ID
func GetUser(id int, full bool) (User, error) {
	var user User

	if full {
		e := DB.Preload(clause.Associations).First(&user, "id = ?", id)
		return user, e.Error
	}

	e := DB.First(&user, "id = ?", id)
	return user, e.Error
}

// GetUserByEmail attempts to return a user struct by email address
func GetUserByEmail(email string, full bool) (User, error) {
	var user User

	if full {
		e := DB.Preload(clause.Associations).First(&user, "email = ?", email)
		return user, e.Error
	}

	e := DB.First(&user, "email = ?", email)
	return user, e.Error
}

// GetUserByAPIKey attempts to return a user struct by api key
func GetUserByAPIKey(api string, full bool) (User, error) {
	var user User

	if full {
		e := DB.Preload(clause.Associations).First(&user, "api_key = ?", api)
		return user, e.Error
	}

	e := DB.First(&user, "api_key = ?", api)
	return user, e.Error
}

// AuthenticateUserByPassword attempts to return a user struct by email address and password
func AuthenticateUserByPassword(email string, password string) (User, error) {
	var user User
	e := DB.First(&user, "email = ? AND password = ?", email, PasswordHash(password))
	return user, e.Error
}

// PasswordHash calculates and returns the SHA256 sum of the passed string
func PasswordHash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

// IssueJWT generates and returns a JWT token
func IssueJWT(user User, expires int) (string, error) {

	// If we have an expiry time, set it...
	if expires >= 0 {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"aud": utilities.AppConfig.JWTAudience,
			"nbf": time.Now().Unix(),
			"exp": time.Now().Unix() + int64(expires),
			"sub": user.Email,
		})
		return token.SignedString([]byte(utilities.AppConfig.JWTKey))
	}

	// ...otherwise, return a token that never expires...
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": utilities.AppConfig.JWTAudience,
		"nbf": time.Now().Unix(),
		"sub": user.Email,
	})
	return token.SignedString([]byte(utilities.AppConfig.JWTKey))

}

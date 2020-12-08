package models

import (
	"path/filepath"
	"strings"

	"github.com/kaigoh/ganjapp/utilities"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// DB GORM instance
var DB = Connect()

// Connect to a database
func Connect() *gorm.DB {

	var db *gorm.DB

	// If the database has been overridden, try and connect...
	switch strings.ToLower(utilities.GetEnv("GANJAPP_DATABASE_TYPE", "sqlite")) {
	case "mysql":
		db = connectMysql()
	case "postgres":
		db = connectPostgres()
	case "sqlserver":
		db = connectSQLServer()
	default:
		db = connectSQLite()
	}

	// Run migrations
	Migrate(db)

	return db

}

func connectMysql() *gorm.DB {
	db, err := gorm.Open(mysql.Open(utilities.GetEnv("GANJAPP_DATABASE_DSN", "")), &gorm.Config{})
	if err != nil {
		panic("MySQL Error: Unable to connect to database, please check your configuration and restart the server")
	}
	return db
}

func connectPostgres() *gorm.DB {
	db, err := gorm.Open(postgres.Open(utilities.GetEnv("GANJAPP_DATABASE_DSN", "")), &gorm.Config{})
	if err != nil {
		panic("Postgres Error: Unable to connect to database, please check your configuration and restart the server")
	}
	return db
}

func connectSQLServer() *gorm.DB {
	db, err := gorm.Open(sqlserver.Open(utilities.GetEnv("GANJAPP_DATABASE_DSN", "")), &gorm.Config{})
	if err != nil {
		panic("SQL Server Error: Unable to connect to database, please check your configuration and restart the server")
	}
	return db
}

func connectSQLite() *gorm.DB {

	path := utilities.GetEnv("GANJAPP_DATA_ROOT", filepath.Join(utilities.GetEnv("GANJAPP_ROOT", ""), "data"))
	filename := filepath.Join(path, utilities.GetEnv("GANJAPP_DATABASE_FILENAME", "ganjapp.db"))

	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		panic("SQLite Error: Unable to connect to database, please check your configuration and restart the server")
	}
	return db
}

// Migrate all models
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Event{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&CannabisStrain{})
	db.AutoMigrate(&Environment{})
	db.AutoMigrate(&EnvironmentImage{})
	db.AutoMigrate(&EnvironmentExtendedData{})
	db.AutoMigrate(&Tree{})
	db.AutoMigrate(&TreeImage{})
	db.AutoMigrate(&TreeExtendedData{})
	db.AutoMigrate(&Shroom{})
	db.AutoMigrate(&ShroomImage{})
	db.AutoMigrate(&ShroomExtendedData{})
}

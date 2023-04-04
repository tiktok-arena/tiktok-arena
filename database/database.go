package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"tiktok-arena/configuration"
	"tiktok-arena/models"
)

var DB *gorm.DB

func ConnectDB(config *configuration.EnvConfigModel) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		config.DBHost,
		config.DBUserName,
		config.DBUserPassword,
		config.DBName,
		config.DBPort,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database!\n", err.Error())
	}
	//	Extension for postgresql uuid support
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	//  Extension for search
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"fuzzystrmatch\"")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Tournament{},
		&models.Tiktok{},
	)
	if err != nil {
		log.Fatal("Migration Failed:\n", err.Error())
	}

	tiktoksTable = DB.Table("tiktoks")
	tournamentsTable = DB.Table("tournaments")
	usersTable = DB.Table("users")

	log.Println("Successfully connected to the database")
}

var (
	tiktoksTable     *gorm.DB
	tournamentsTable *gorm.DB
	usersTable       *gorm.DB
)

// Scopes for search and pagination

func Search(searchText string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if searchText == "" {
			return db
		}
		return db.Select("*, levenshtein(name, ?) as distance", searchText).Order("distance")
	}
}

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	offset := (page - 1) * pageSize
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(pageSize)
	}
}

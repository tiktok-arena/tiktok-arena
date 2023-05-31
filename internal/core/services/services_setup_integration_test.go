package services

import (
	"gorm.io/gorm"
	"log"
	"testing"
	"tiktok-arena/configuration"
	"tiktok-arena/internal/data/database"
)

type DatabaseIntegrationTest struct {
	db *gorm.DB
}

func SetupIntegration(t *testing.T) *DatabaseIntegrationTest {
	err := configuration.LoadConfig("../../../.env.test")
	if err != nil {
		t.Fatal("Failed to load environment variables!", err.Error())
	}
	db := database.ConnectDB(&configuration.EnvConfig)
	return &DatabaseIntegrationTest{
		db: db,
	}
}

func (test *DatabaseIntegrationTest) Cleanup() {
	res := test.db.Exec("TRUNCATE TABLE users, tiktoks, tournaments CASCADE;")
	if res.Error != nil {
		log.Fatal("Error in cleanup: \n", res.Error)
	}
}

package bootstrap

import (
	"fmt"

	"github.com/koropati/population-recap/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabase(config *Config) *gorm.DB {

	dbHost := config.DBHost
	dbPort := config.DBPort
	dbUser := config.DBUser
	dbPass := config.DBPass
	dbName := config.DBName

	dbConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dbConnString), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	AutoMigrate(db)
	return db
}

func CloseDatabase(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		panic("failed to close database")
	}
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&domain.User{},
		&domain.AccessToken{},
		&domain.RefreshToken{},
		&domain.ForgotPasswordToken{},
	)
}

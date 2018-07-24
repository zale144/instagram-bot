package model

import (
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	DBInfo string
	DB     *gorm.DB
	//DE, _ = time.LoadLocation("Europe/Berlin")
)

// openDB prepares a db connection for postgres
func openDB() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}
	var err error
	DB, err = gorm.Open("postgres", DBInfo)
	if err != nil {
		return nil, err
	}
	DB.DB().SetMaxOpenConns(400)
	DB.DB().SetMaxIdleConns(0)
	DB.DB().SetConnMaxLifetime(100 * time.Second)
	return DB, nil
}

// InitDB adds tables if needed and sets the log mode
func InitDB() error {
	_, err := openDB()
	if err != nil {
		return err
	}
	// Add tables here
	DB.AutoMigrate(&ProcessedUser{}, &Job{})

	// Enable Logger, show detailed log
	DB.LogMode(true)
	return nil
}

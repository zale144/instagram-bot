package model

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DBInfo string
	DB     *gorm.DB
	//DE, _ = time.LoadLocation("Europe/Berlin")
)

func openDB() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}
	var err error
	DB, err = gorm.Open("mysql", DBInfo)
	if err != nil {
		return nil, err
	}
	DB.DB().SetMaxOpenConns(400)
	DB.DB().SetMaxIdleConns(0)
	DB.DB().SetConnMaxLifetime(100 * time.Second)
	return DB, nil
}

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

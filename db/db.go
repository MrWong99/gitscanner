package db

import (
	"github.com/MrWong99/gitscanner/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb() error {
	var err error
	db, err = gorm.Open(sqlite.Open("scanresults.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	if err = db.AutoMigrate(&utils.CheckResultConsolidated{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&utils.SingleCheck{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&utils.GlobalConfig{}); err != nil {
		return err
	}
	return nil
}

func Get() *gorm.DB {
	return db
}

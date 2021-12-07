package configrepo

import (
	"errors"

	"github.com/MrWong99/gitscanner/checks"
	"github.com/MrWong99/gitscanner/db"
	"gorm.io/gorm"
)

func ReadConfig(checkname string) (*checks.CheckConfiguration, error) {
	config := &checks.CheckConfiguration{}
	tx := db.Get().Where("check_name = ?", checkname).First(config)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, tx.Error
		}
	}
	return config, nil
}

func UpdateConfig(cfg *checks.CheckConfiguration) error {
	tx := db.Get().Save(cfg)
	return tx.Error
}

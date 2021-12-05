package configrepo

import (
	"errors"

	"github.com/MrWong99/gitscanner/db"
	"github.com/MrWong99/gitscanner/utils"
	"gorm.io/gorm"
)

func ReadConfig() (*utils.GlobalConfig, error) {
	config := &utils.GlobalConfig{}
	tx := db.Get().First(config)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, tx.Error
		}
	}
	return config, nil
}

func UpdateConfig(cfg *utils.GlobalConfig) error {
	currentCfg, err := ReadConfig()
	if err != nil {
		return err
	}
	if currentCfg == nil {
		tx := db.Get().Create(cfg)
		return tx.Error
	}
	currentCfg.BranchPattern = cfg.BranchPattern
	currentCfg.NamePattern = cfg.EmailPattern
	currentCfg.NamePattern = cfg.NamePattern
	tx := db.Get().Save(currentCfg)
	return tx.Error
}

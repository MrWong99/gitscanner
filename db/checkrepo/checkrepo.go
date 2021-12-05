package checkrepo

import (
	"time"

	"github.com/MrWong99/gitscanner/db"
	"github.com/MrWong99/gitscanner/utils"
)

func SaveChecks(checks []*utils.CheckResultConsolidated) error {
	database := db.Get()
	for _, v := range checks {
		tx := database.Create(v)
		if tx.Error != nil {
			return tx.Error
		}
	}
	return nil
}

func AcknowledgeCheck(singleCheckID uint, acknowledged bool) error {
	check := &utils.SingleCheck{}
	database := db.Get()
	tx := database.First(check, singleCheckID)
	if tx.Error != nil {
		return tx.Error
	}
	check.Acknowledged = acknowledged
	if tx = database.Save(check); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func ReadSavedChecks(checkNames []string, startDate, endDate time.Time) ([]*utils.CheckResultConsolidated, error) {
	database := db.Get()
	results := []*utils.CheckResultConsolidated{}
	if tx := database.Where("date BETWEEN ? AND ?", startDate, endDate).Find(&results); tx.Error != nil {
		return nil, tx.Error
	}
	for _, res := range results {
		allowedChecks := []utils.SingleCheck{}
		if tx := database.Where("check_result_consolidated_id = ? AND check_name IN ?", res.ID, checkNames).Find(&allowedChecks); tx.Error != nil {
			return nil, tx.Error
		}
		res.Checks = allowedChecks
	}
	return results, nil
}

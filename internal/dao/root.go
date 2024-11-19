/*
 * @Author: Liu Sainan
 * @Date: 2024-01-07 17:27:33
 */

package dao

import (
	"myadmin/internal/config"
	"myadmin/internal/db"
	"myadmin/internal/model"
	"slices"
)

func AutoCreateTable() error {
	tables, err := db.DB().Conn().Migrator().GetTables()
	if err != nil {
		return err
	}
	switch config.GlobalConfig.Server.AutoCreateTable {
	case -1:
		return nil
	case 0:
		if len(tables) != 0 {
			return nil
		}
		return db.DB().Conn().AutoMigrate(TableListForNotExist(tables)...)
	case 1:
		return db.DB().Conn().AutoMigrate(TableListForNotExist(tables)...)
	case 2:
		return db.DB().Conn().AutoMigrate(TableListForNotExist([]string{})...)
	default:
		return nil
	}
}

func TableListForNotExist(tables []string) []any {
	var tbs []any
	for tbname, table := range model.Tables {
		if !slices.Contains(tables, tbname) {
			tbs = append(tbs, table)
		}
	}
	return tbs
}

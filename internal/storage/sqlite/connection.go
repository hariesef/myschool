package sqlite

import (
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DBFile string = "./internal/storage/sqlite/main.db"

func Connect() (*gorm.DB, error) {
	//SQLite setup local file
	sqlDB, err := gorm.Open(sqliteDriver.Open(DBFile), &gorm.Config{})
	return sqlDB, err
}

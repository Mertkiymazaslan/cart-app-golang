package database

import (
	"checkoutProject/pkg/common/env"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Initialize() error {
	var err error
	db, err = gorm.Open(postgres.Open(env.DB_URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("cannot connect to the database. Error: %s", err.Error())
	}

	return nil
}

func GetInstance() *gorm.DB {
	if db == nil {
		panic("database instance is nil, please connect to the database first")
	}

	return db
}

func NewTransaction() *gorm.DB {
	return GetInstance().Begin()
}

func CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func RollbackTransaction(tx *gorm.DB) error {
	if err := tx.Rollback().Error; err != nil && !errors.Is(err, sql.ErrTxDone) {
		return err
	}

	return nil
}

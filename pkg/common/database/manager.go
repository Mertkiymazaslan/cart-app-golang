package database

import "gorm.io/gorm"

type BaseManager struct {
	DB *gorm.DB
}

func (m BaseManager) WithTx(tx *gorm.DB) BaseManager {
	if tx != nil {
		m.DB = tx
	}

	return m
}

func NewBaseManager(tx *gorm.DB) BaseManager {
	m := BaseManager{
		DB: GetInstance(),
	}

	if tx != nil {
		m.DB = tx
	}

	return m
}

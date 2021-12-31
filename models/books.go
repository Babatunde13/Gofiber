package models

import "gorm.io/gorm"

type Book struct {
	ID        uint   `gorm:"primary_key"`
	Title     string `gorm:"type:varchar(255);not null"`
	Author    string `gorm:"type:varchar(255);not null"`
	Publisher string `gorm:"type:varchar(255);not null"`
}

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Book{})
	return err
}

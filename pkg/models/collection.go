package models

import "gorm.io/gorm"

type Collection struct {
	gorm.Model
	Name  string `json:"name" binding:"required" validate:"required"`
	Books []Book `gorm:"many2many:book_collections;"`
}

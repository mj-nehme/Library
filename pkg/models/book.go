package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title       string    `json:"title" binding:"required" validate:"required" gorm:"size:255"`
	Author      string    `json:"author" binding:"required" validate:"required" gorm:"size:255"`
	Published   time.Time `json:"published" validate:"lte"`
	Edition     int       `json:"edition" validate:"gte=1"`
	Description string    `json:"description" gorm:"size:1000"`
	GenreName   string    `json:"genre_name" gorm:"size:255"`
}

func (b *Book) BeforeSave(tx *gorm.DB) error {
	b.Published = b.Published.UTC()
	return nil
}

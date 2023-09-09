package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title       string    `json:"title" binding:"required" validate:"required"`
	Author      string    `json:"author" binding:"required" validate:"required"`
	Published   time.Time `json:"published"`
	Edition     int       `json:"edition" validate:"gte=1"`
	Description string    `json:"description"`
	GenreName   string    `json:"genre_name"`
}

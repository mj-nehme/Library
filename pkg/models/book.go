package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title       string    `json:"title" binding:"required"`
	Author      string    `json:"author" binding:"required"`
	Published   time.Time `json:"published"`
	Edition     int       `json:"edition"`
	Description string    `json:"description"`
	GenreName   string    `json:"genre_name"`
}

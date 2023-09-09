package models

type Collection struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

package models

type CollectionBook struct {
	CollectionID int `json:"collection_id" binding:"required"`
	BookID       int `json:"book_id" binding:"required"`
}

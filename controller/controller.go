package controller

import (
	"github.com/RianWardanaPutra/notes-v1/model"
	"github.com/jinzhu/gorm"
)

type Controller struct {
	db *gorm.DB
}

// Create new db instance
func NewController(db *gorm.DB) *Controller {
	return &Controller{db}
}

type Returns struct {
	ID    uint         `json:"user_id"`
	Notes []model.Note `json:"notes"`
}

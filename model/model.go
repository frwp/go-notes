package model

import (
	"time"
)

type Basic struct {
	ID        uint       `json:"id" gorm:"PRIMARY_KEY"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}

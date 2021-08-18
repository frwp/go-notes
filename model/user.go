package model

type User struct {
	Basic `gorm:"EMBEDDED"`
	Name  string `json:"name" gorm:"NOT NULL"`
	Email string `json:"email" gorm:"UNIQUE;NOT NULL"`
}

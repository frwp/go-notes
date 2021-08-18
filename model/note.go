package model

type Note struct {
	Basic       `gorm:"EMBEDDED"`
	UserId      int    `json:"user_id" gorm:"NOT NULL;column:user_id"`
	Title       string `json:"title" gorm:"UNIQUE;NOT NULL"`
	Description string `json:"description"`
}

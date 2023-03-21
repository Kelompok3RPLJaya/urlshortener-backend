package entity

import "github.com/google/uuid"

type userPrivate struct {
	ID             uuid.UUID `gorm:"primary_key;not_null"`
	Password       string    `json:"password"`
	UrlShortenerID uuid.UUID `gorm:"foreignkey:ID"`
}

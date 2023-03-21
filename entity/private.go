package entity

import "github.com/google/uuid"

type userPrivate struct {
	ID             uuid.UUID `gorm:"primary_key;not_null"`
	Password       string    `json:"password"`

	UrlShortenerID   uuid.UUID `gorm:"foreignKey" json:"url_shortener_id"`
	UrlShortener     *UrlShortener     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"url_shortener,omitempty"`

	Timestamp
}

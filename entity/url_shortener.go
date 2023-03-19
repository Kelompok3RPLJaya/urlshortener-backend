package entity

import "github.com/gofrs/uuid"

type UrlShortener struct {
	ID       uuid.UUID `gorm:"primary_key;not_null" json:"id"`
	ShortUrl string    `json:"short_url"`
	LongUrl  string    `json:"long_url"`
	Views    uint64    `json:"views"`
	Private  bool      `json:"private"`
	UserID   uuid.UUID `gorm:"foreignKey" json:"user_id"`
	User     *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`

	// Timestamp
}

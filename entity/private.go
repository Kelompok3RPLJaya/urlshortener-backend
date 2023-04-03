package entity

import (
	"url-shortener-backend/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Private struct {
	ID             uuid.UUID `gorm:"primary_key;not_null"`
	Password       string    `json:"password"`

	UrlShortenerID   uuid.UUID `gorm:"foreignKey" json:"url_shortener_id"`
	UrlShortener     *UrlShortener     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"url_shortener,omitempty"`

	Timestamp
}

func (u *Private) BeforeCreate(tx *gorm.DB) error {
	var err error
	u.Password, err = helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}
	return nil
}
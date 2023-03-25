package dto

import (
	"url-shortener-backend/entity"

	"github.com/google/uuid"
)

type UrlShortenerCreateDTO struct {
	ID        uuid.UUID `gorm:"primary_key;not_null" json:"id"`
	Title  	  string    `json:"title" binding:"required"`
	ShortUrl  string    `json:"short_url" binding:"required"`
	LongUrl   string    `json:"long_url" binding:"required"`
	Views     uint64    `json:"views" form:"views"`
	IsPrivate *bool     `json:"is_private" binding:"required"`
	IsFeeds   *bool     `json:"is_feeds" binding:"required"`

	UserID uuid.UUID    `gorm:"foreignKey" json:"user_id"`
	User   *entity.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`

	Password string `json:"password"`
}

type UrlShortenerUpdateDTO struct {
	ID        uuid.UUID `gorm:"primary_key;not_null" json:"id"`
	Title  	  string    `json:"title"`
	ShortUrl  string    `json:"short_url"`
	LongUrl   string    `json:"long_url"`
	IsPrivate bool      `json:"is_private"`
	IsFeeds   bool      `json:"is_feeds"`

	Password string `json:"password"`
}

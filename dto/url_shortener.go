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

type UrlShortenerResponseDTO struct {
	ID        	uuid.UUID   `gorm:"primary_key;not_null" json:"id"`
	Title 		string 		`json:"title" form:"title"`
	LongUrl 	string 		`json:"long_url" form:"long_url"`
	ShortUrl 	string 		`json:"short_url" form:"short_url"`
	Views 		uint64  	`json:"views" form:"views"`
	IsPrivate	*bool		`json:"is_private" form:"is_private"`
	IsFeeds		*bool		`json:"is_feeds" form:"is_feeds"`
	Username	string		`json:"username"`

	UserID   	uuid.UUID 		`gorm:"foreignKey" json:"user_id"`

	CreatedAt 	string 	`json:"created_at"`
	UpdatedAt 	string 	`json:"updated_at"`
}

func BoolPointer(b bool) *bool {
    return &b
}
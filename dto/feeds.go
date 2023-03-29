package dto

import (
	"time"

	"github.com/google/uuid"
)

type Data struct {
	Before	string	`json:"before"`
	After	string	`json:"after"`
}

type FeedsResponseDTO struct {
	ID        	uuid.UUID   `gorm:"primary_key;not_null" json:"id"`
	Title		string		`json:"title"`
	Username	string		`json:"username"`
	Method		string		`json:"method"`
	UrlShortenerID   	uuid.UUID 			`gorm:"foreignKey" json:"url_shortener_id"`
	Data 		Data 		`json:"data"`

	CreatedAt 	time.Time 	`json:"created_at"`
}

package entity

import (
	"github.com/google/uuid"
)

type Feeds struct {
	ID        	uuid.UUID   `gorm:"primary_key;not_null" json:"id"`
	Data 		string 		`json:"data"`
	Method		string		`json:"method"`

	UrlShortenerID   	uuid.UUID 			`gorm:"foreignKey" json:"url_shortener_id"`
	UrlShortener     	*UrlShortener  		`gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"url_shortener,omitempty"`
	
	Timestamp
}
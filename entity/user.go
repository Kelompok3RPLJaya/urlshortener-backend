package entity

import (
	"github.com/google/uuid"
)

type User struct {
	ID        	uuid.UUID   `gorm:"primary_key;not_null" json:"id"`
	Name 		string 		`json:"name"`
	Email 		string 		`json:"email" binding:"email"`
	Password 	string  	`json:"password"`
	Role		string		`json:"role"`
	
	Timestamp
}
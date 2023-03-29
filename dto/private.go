package dto

import "github.com/google/uuid"

type PrivateUpdateDTO struct {
	ID        	uuid.UUID   `gorm:"primary_key;not_null" json:"id"`
	Password 	string 		`json:"password"`
}
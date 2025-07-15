package models

import (
	"time"
)

type User struct {
	ID        	uint      `json:"id" gorm:"primaryKey"`
	Name      	string    `json:"name" gorm:"not null" binding:"required"`
	Email     	string    `json:"email" gorm:"unique;not null" binding:"required"`
	Password  	string    `json:"-" gorm:"not null"`
	ValidFrom 	time.Time `json:"valid_from" gorm:"not null"`
	ValidUntil 	time.Time `json:"valid_until" gorm:"not null"`
}

type UserResponse struct {
	ID        	uint      `json:"id"`
	Name      	string    `json:"name"`
	Email     	string    `json:"email"`
	ValidFrom 	time.Time `json:"valid_from" gorm:"not null"`
	ValidUntil 	time.Time `json:"valid_until" gorm:"not null"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

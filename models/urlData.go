package models

import (
	"time"
)

type BrokenLink struct {
	URL        string `json:"url"`
	StatusCode int    `json:"statusCode"`
	Type       string `json:"type"` // "internal" or "external"
}

type UrlAnalysisResult struct {
    ID     			uint  	`json:"id" gorm:"primaryKey;autoIncrement"`
    Title  			string  `json:"title"`
    URL 			string  `json:"url" gorm:"not null" binding:"required"`
    HTMLVersion 	string 	`json:"htmlVersion"`
    InternalLinks 	int 	`json:"internalLinks"`
    ExternalLinks 	int 	`json:"externalLinks"`
	Status 			string 	`json:"status"`
	LoginForm 		bool 	`json:"loginForm"`
	BrokenLinks    	[]BrokenLink `json:"brokenLinks" gorm:"type:jsonb;serializer:json"`
	ValidFrom 		time.Time `json:"valid_from" gorm:"not null"`
	ValidUntil 		time.Time `json:"valid_until" gorm:"not null"`
}
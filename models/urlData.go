package models

type UrlAnalysisResult struct {
    ID     			int  	`json:"id" gorm:"primaryKey"`
    Title  			string  `json:"title" gorm:"not null" binding:"required"`
    URL 			string  `json:"url" gorm:"not null" binding:"required"`
    HTMLVersion 	string 	`json:"htmlVersion"`
    InternalLinks 	int 	`json:"internalLinks"`
    ExternalLinks 	int 	`json:"externalLinks"`
	Status 			string 	`json:"status"`
	LoginForm 		bool 	`json:"loginForm"`
}
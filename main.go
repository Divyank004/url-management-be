package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"url-management-be/config"
	"url-management-be/models"
	"time"
)

func setupRouter() *gin.Engine {

	// Initialize database
	config.ConnectDB()

	// Migrate the schema to db 
	config.DB.AutoMigrate(&models.User{})
	config.DB.AutoMigrate(&models.UrlAnalysisResult{})
	// Seed sample data into db 
	config.DB.Create(&models.User{
		ID:         1,
		Name:       "Divyank Dhadi",
		Email:      "divyank004@gmail.com",
		Password:   "password", // todo Hash the password
		ValidFrom:  time.Now(),
		ValidUntil: time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	config.DB.Create(&models.UrlAnalysisResult{
		ID:             1,
		URL:            "https://google.com",
		Title:          "Google Search",
		HTMLVersion:    "HTML5",
		InternalLinks:  25,
		ExternalLinks:  12,
		Status:         "Running",
		LoginForm:      true,
	})

	
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

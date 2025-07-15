package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"url-management-be/config"
	"url-management-be/models"
)

func setupRouter() *gin.Engine {

	// Initialize database
	config.ConnectDB()
	
	// Migrate the schema to db 
	config.DB.AutoMigrate(&models.User{})
	config.DB.AutoMigrate(&models.UrlAnalysisResult{})
	
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

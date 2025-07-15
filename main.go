package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"url-management-be/config"
	"url-management-be/models"
	"time"
	"url-management-be/utils"
	"url-management-be/middleware"
)

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token": token,
		"user": models.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			ValidFrom:  user.ValidFrom,
			ValidUntil: user.ValidUntil,
		},
	})
}

func setupRouter() *gin.Engine {

	// Initialize database
	config.ConnectDB()

	// Migrate the schema to db 
	config.DB.AutoMigrate(&models.User{})
	config.DB.AutoMigrate(&models.UrlAnalysisResult{})
	
	hashedPassword, err := utils.HashPassword("password")
	if err != nil {
		log.Fatal("Failed to seed data into database:", err)
	}
	// Seed sample data into db 
	config.DB.Create(&models.User{
		ID:         1,
		Name:       "Divyank Dhadi",
		Email:      "divyank004@gmail.com",
		Password:   hashedPassword,
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
	r.Use(middleware.CORSMiddleware())
	public := r.Group("/api")
	public.POST("/login", Login)
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

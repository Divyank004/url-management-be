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
	"url-management-be/controllers"
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
	// Drop existing tables if they exist
	config.DB.Migrator().DropTable(&models.User{})
	config.DB.Migrator().DropTable(&models.UrlAnalysisResult{})
	// Migrate the schema to db 
	config.DB.AutoMigrate(&models.User{})
	config.DB.AutoMigrate(&models.UrlAnalysisResult{})
	
	hashedPassword, err := utils.HashPassword("password")
	if err != nil {
		log.Fatal("Failed to seed data into database:", err)
	}
	// Seed sample data into db 
	config.DB.Create(&models.User{
		Name:       "Divyank Dhadi",
		Email:      "divyank004@gmail.com",
		Password:   hashedPassword,
		ValidFrom:  time.Now(),
		ValidUntil: time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	config.DB.Create(&models.UrlAnalysisResult{
		URL:            "https://google.com",
		Title:          "Google Search",
		HTMLVersion:    "HTML5",
		InternalLinks:  25,
		ExternalLinks:  12,
		Status:         "Done",
		LoginForm:      true,
		BrokenLinks: []models.BrokenLink{
			{
				URL:        "https://google.com/broken-link",
				StatusCode: 404,
				Type:       "internal",
			},
			{
				URL:        "https://external.com/missing",
				StatusCode: 500,
				Type:       "external",
			},
		},
		ValidFrom:  	time.Now(),
		ValidUntil: 	time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "API is running")
	})
	public := r.Group("/api")
	public.POST("/login", Login)
	public.Use(middleware.AuthMiddleware())
	public.GET("/urldata", controllers.GetAllUrlsAnalysisData)
	public.GET("/urldata/:id", controllers.GetSingleUrlAnalysisData)
	public.POST("/addurl", controllers.AddUrl)
	
	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

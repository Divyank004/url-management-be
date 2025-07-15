package controllers

import (
	"time"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"url-management-be/config"
	"url-management-be/models"
)

func GetAllUrlsAnalysisData(c *gin.Context) {
	var urldata []models.UrlAnalysisResult
	if err := config.DB.Order("valid_from DESC").Find(&urldata).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch urldata"})
		return
	}
	c.JSON(http.StatusOK, urldata)
}

func AddUrl(c *gin.Context) {
	var input models.UrlAnalysisResult

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	input.ValidFrom = time.Now()
	input.ValidUntil = time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed adding URL data"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func GetSingleUrlAnalysisData(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid urlAnalysisResult ID"})
		return
	}

	var urlData models.UrlAnalysisResult
	if err := config.DB.First(&urlData, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL analysis data not found"})
		return
	}

	c.JSON(http.StatusOK, urlData)
}

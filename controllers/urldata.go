package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"url-management-be/config"
	"url-management-be/models"
)

func GetUrlAnalysisData(c *gin.Context) {
	var urldata []models.UrlAnalysisResult
	if err := config.DB.Find(&urldata).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch urldata"})
		return
	}

	var urlAnalysisResults []models.UrlAnalysisResult
	for _, url := range urldata {
		urlAnalysisResults = append(urlAnalysisResults, models.UrlAnalysisResult{
			ID:             url.ID,
			URL:            url.URL,
			Title:          url.Title,
			HTMLVersion:    url.HTMLVersion,
			InternalLinks:  url.InternalLinks,
			ExternalLinks:  url.ExternalLinks,
			Status:         url.Status,
			LoginForm:      url.LoginForm,
		})
	}

	c.JSON(http.StatusOK, urlAnalysisResults)
}
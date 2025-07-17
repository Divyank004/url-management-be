package controllers

import (
	"fmt"
	"time"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"strings"

	"github.com/gin-gonic/gin"
	"url-management-be/config"
	"url-management-be/models"
	"github.com/gocolly/colly/v2"
)

const (
	StatusQueued  string = "Queued"
	StatusRunning string = "Running"
	StatusDone    string = "Done"
	StatusError  string = "Error"
)

var (
	jobs = make(map[string]*models.UrlAnalysisResult)
	mu   sync.RWMutex
)
func GetAllUrlsAnalysisData(c *gin.Context) {
	var urldata []models.UrlAnalysisResult
	if err := config.DB.Order("valid_from DESC").Find(&urldata).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch urldata"})
		return
	}
	c.JSON(http.StatusOK, urldata)
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
	// Create new job
	jobID := input.ID
	job := &models.UrlAnalysisResult{
		ID:        jobID,
		URL:       input.URL,
		Status:    StatusQueued,
	}

	mu.Lock()
	jobIDStr := strconv.FormatUint(uint64(jobID), 10)
	jobs[jobIDStr] = job
	mu.Unlock()

	// Start crawling in background
	go performCrawl(jobIDStr)

	c.JSON(http.StatusCreated, input)
}

func ReRunAnalysis(c *gin.Context) {
	id := c.Param("id")
	
	var urlData models.UrlAnalysisResult
	if err := config.DB.First(&urlData, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL analysis data not found"})
		return
	}
	// Create new job
	jobID := urlData.ID
	job := &models.UrlAnalysisResult{
		ID:        jobID,
		URL:       urlData.URL,
		Status:    StatusQueued,
	}

	mu.Lock()
	jobIDStr := strconv.FormatUint(uint64(jobID), 10)
	jobs[jobIDStr] = job
	mu.Unlock()

	// Start crawling in background
	go performCrawl(jobIDStr)

	c.JSON(http.StatusCreated, urlData)
}

func GetURLAnalysisResult(c *gin.Context) {
	jobID := c.Param("id")
	mu.RLock()
	job, exists := jobs[jobID]
	mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}


func performCrawl(jobID string) {
	mu.Lock()
	job := jobs[jobID]
	job.Status = StatusRunning
	mu.Unlock()

	// Create a new collector
	c := colly.NewCollector()

	// Set user agent
	c.UserAgent = "Web Crawler Bot 1.0"

	var (
		htmlVersion   string
		title         string
		internalLinks int
		externalLinks int
		brokenLinks   []models.BrokenLink
		baseURL       *url.URL
	)

	// Parse base URL
	var err error
	baseURL, err = url.Parse(job.URL)
	if err != nil {
		mu.Lock()
		job.Status = StatusError
		mu.Unlock()
		return
	}

	// Extract HTML version and title
	c.OnHTML("html", func(e *colly.HTMLElement) {
		htmlVersion = getHTMLVersion(e.Text)
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		title = strings.TrimSpace(e.Text)
	})

	// Count links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if href == "" {
			return
		}

		linkURL, err := url.Parse(href)
		if err != nil {
			return
		}

		// Resolve relative URLs
		linkURL = baseURL.ResolveReference(linkURL)

		// Check if internal or external
		linkType := "external"
		if linkURL.Host == baseURL.Host {
			internalLinks++
			linkType = "internal"
		} else {
			externalLinks++
		}

		// Check if link is broken (in background to avoid blocking)
		go checkBrokenLink(linkURL.String(), linkType , &brokenLinks, &mu)
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		mu.Lock()
		job.Status = StatusError
		mu.Unlock()
	})

	// Start crawling
	err = c.Visit(job.URL)
	if err != nil {
		mu.Lock()
		job.Status = StatusError
		mu.Unlock()
		return
	}

	// Wait a bit for broken link checks to complete
	time.Sleep(3 * time.Second)

	// Update job with results
	mu.Lock()
	job.HTMLVersion = htmlVersion
	job.Title = title
	job.InternalLinks = internalLinks
	job.ExternalLinks = externalLinks
	job.BrokenLinks = brokenLinks
	job.Status = StatusDone
	mu.Unlock()

	// Save results to database
	var urldata models.UrlAnalysisResult
	if err := config.DB.First(&urldata, jobID).Error; err != nil {
		fmt.Println("error: url not found", err)
		return
	}
	urldata.Title = title
	urldata.HTMLVersion = htmlVersion
	urldata.InternalLinks = internalLinks
	urldata.ExternalLinks = externalLinks
	urldata.Status = StatusDone
	urldata.BrokenLinks = brokenLinks

	if err := config.DB.Save(&urldata).Error; err != nil {
		fmt.Println("error: Failed to update url data", err)
		return
	}
}

func getHTMLVersion(content string) string {
	content = strings.ToLower(content)
	
	if strings.Contains(content, "<!doctype html>") {
		return "HTML5"
	}
	if strings.Contains(content, "html 4.01") {
		return "HTML 4.01"
	}
	if strings.Contains(content, "xhtml 1.0") {
		return "XHTML 1.0"
	}
	if strings.Contains(content, "xhtml 1.1") {
		return "XHTML 1.1"
	}
	
	return "Unknown"
}

func checkBrokenLink(linkURL string, linkType string, brokenLinks *[]models.BrokenLink, mu *sync.RWMutex) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(linkURL)
	if err != nil {
		mu.Lock()
		*brokenLinks = append(*brokenLinks, models.BrokenLink{
			URL:        linkURL,
			Type:       linkType,
			StatusCode: 0, // Network error
		})
		mu.Unlock()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		mu.Lock()
		*brokenLinks = append(*brokenLinks, models.BrokenLink{
			URL:        linkURL,
			Type:       linkType,
			StatusCode: resp.StatusCode,
		})
		mu.Unlock()
	}
}

func DeleteUrl(c *gin.Context) {
	id := c.Param("id")
	var urlData models.UrlAnalysisResult

	if err := config.DB.First(&urlData, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL analysis data not found"})
		return
	}

	if err := config.DB.Delete(&urlData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete URL analysis data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL analysis data deleted successfully"})
}
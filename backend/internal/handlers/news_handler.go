package handlers

import (
	"net/http"

	"fact-check/internal/models"
	"fact-check/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type NewsHandler struct {
	newsService   *services.NewsService
	openAIService *services.OpenAIService
	logger        *logrus.Logger
}

func NewNewsHandler(newsService *services.NewsService, openAIService *services.OpenAIService, logger *logrus.Logger) *NewsHandler {
	return &NewsHandler{
		newsService:   newsService,
		openAIService: openAIService,
		logger:        logger,
	}
}

// Submit handles news submission
func (h *NewsHandler) Submit(c *gin.Context) {
	var submission models.NewsSubmission
	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Submit news
	news, err := h.newsService.SubmitNews(userID.(string), &submission)
	if err != nil {
		h.logger.Errorf("Failed to submit news: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit news"})
		return
	}

	c.JSON(http.StatusCreated, news)
}

// Verify handles news verification using OpenAI
func (h *NewsHandler) Verify(c *gin.Context) {
	newsID := c.Param("id")
	if newsID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "News ID is required"})
		return
	}

	// Get news from database
	news, err := h.newsService.GetNewsByID(newsID)
	if err != nil {
		h.logger.Errorf("Failed to get news: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	// Verify news using OpenAI
	var link, photoURL string
	if news.Link != nil {
		link = *news.Link
	}
	if news.PhotoURL != nil {
		photoURL = *news.PhotoURL
	}
	status, explanation, err := h.openAIService.VerifyNews(news.Content, link, photoURL)
	if err != nil {
		h.logger.Errorf("Failed to verify news with OpenAI: %v", err)
		
		// Provide more specific error messages based on the error type
		errorMessage := "Failed to verify news"
		if err.Error() != "" {
			errorMessage = err.Error()
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errorMessage,
			"details": "The fact-checking service is currently unavailable. Please try again later or contact support if the issue persists.",
		})
		return
	}

	// Update news status in database
	err = h.newsService.UpdateNewsStatus(newsID, status, explanation)
	if err != nil {
		h.logger.Errorf("Failed to update news status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update news status"})
		return
	}

	verification := models.NewsVerification{
		ID:          news.ID,
		Status:      status,
		Explanation: explanation,
	}

	c.JSON(http.StatusOK, verification)
}

// GetUserNews retrieves all news submissions for a user
func (h *NewsHandler) GetUserNews(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Verify the requesting user can access this data
	requestingUserID, exists := c.Get("user_id")
	if !exists || requestingUserID.(string) != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get user's news
	newsList, err := h.newsService.GetUserNews(userID)
	if err != nil {
		h.logger.Errorf("Failed to get user news: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve news"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"news":  newsList,
		"count": len(newsList),
	})
}

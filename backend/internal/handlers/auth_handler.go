package handlers

import (
	"net/http"

	"fact-check/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	authService *services.AuthService
	logger      *logrus.Logger
}

func NewAuthHandler(authService *services.AuthService, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Login initiates Google OAuth2 flow
func (h *AuthHandler) Login(c *gin.Context) {
	authURL := h.authService.GetLoginURL()
	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
	})
}

// Callback handles Google OAuth2 callback
func (h *AuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
		return
	}

	// Handle the OAuth callback
	response, err := h.authService.HandleCallback(code)
	if err != nil {
		h.logger.Errorf("Failed to handle OAuth callback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate user"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Me returns the current user's information
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		h.logger.Errorf("Failed to get user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a real application, you might want to invalidate the JWT token
	// For now, we'll just return a success response
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fact-check/internal/config"
	"fact-check/internal/database"
	"fact-check/internal/handlers"
	"fact-check/internal/middleware"
	"fact-check/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logging
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(cfg.LogLevel)

	// Initialize database
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := database.RunMigrations(db); err != nil {
		logger.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize services
	authService := services.NewAuthService(cfg, db, logger)
	newsService := services.NewNewsService(cfg, db, logger)
	openAIService := services.NewOpenAIService(cfg, logger)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, logger)
	newsHandler := handlers.NewNewsHandler(newsService, openAIService, logger)

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger(logger))
	router.Use(middleware.DefaultRateLimiter())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Service status endpoint
		api.GET("/services/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"openai": gin.H{
					"available":  false,
					"configured": false,
					"message":    "OpenAI API quota exceeded - please check your billing",
				},
				"database": gin.H{
					"status":  "connected",
					"message": "Database connection established",
				},
			})
		})

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.GET("/login", authHandler.Login)
			auth.GET("/callback", authHandler.Callback)
			auth.GET("/me", middleware.AuthMiddleware(authService), authHandler.Me)
			auth.POST("/logout", middleware.AuthMiddleware(authService), authHandler.Logout)
		}

		// News routes
		news := api.Group("/news")
		{
			news.POST("/submit", middleware.AuthMiddleware(authService), newsHandler.Submit)
			news.GET("/verify/:id", middleware.AuthMiddleware(authService), newsHandler.Verify)
			news.GET("/user/:id", middleware.AuthMiddleware(authService), newsHandler.GetUserNews)
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Starting server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

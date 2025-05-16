package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/config"
	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/handler"
	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/middleware"
	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize services
	services, err := service.NewServices(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	// Initialize router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimit())

	// Initialize handlers
	handlers := handler.NewHandlers(services)

	// Register routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Auth.Register)
			auth.POST("/login", handlers.Auth.Login)
			auth.POST("/refresh", handlers.Auth.RefreshToken)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.Auth(services.Auth))
		{
			// Subscription routes
			subscription := protected.Group("/subscription")
			{
				subscription.GET("/plans", handlers.Subscription.GetPlans)
				subscription.POST("/subscribe", handlers.Subscription.Subscribe)
				subscription.GET("/status", handlers.Subscription.GetStatus)
				subscription.POST("/cancel", handlers.Subscription.Cancel)
			}

			// ML routes
			ml := protected.Group("/ml")
			{
				ml.POST("/analyze", handlers.ML.Analyze)
				ml.GET("/recommendations", handlers.ML.GetRecommendations)
			}
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

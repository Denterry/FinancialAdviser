package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Denterry/FinancialAdviser/Backend/sub-service/config"
	"github.com/Denterry/FinancialAdviser/Backend/sub-service/internal/controller/grpc"
	"github.com/Denterry/FinancialAdviser/Backend/sub-service/internal/repo/postgres"
	"github.com/Denterry/FinancialAdviser/Backend/sub-service/internal/usecase"
	"github.com/Denterry/FinancialAdviser/Backend/sub-service/pkg/pb/subscription/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	subscriptionRepo := postgres.NewSubscriptionRepository(db)

	// Initialize use cases
	subscriptionUseCase := usecase.NewSubscriptionUseCase(subscriptionRepo)

	// Initialize gRPC server
	server := grpc.NewServer()
	subscriptionService := grpc.NewSubscriptionService(subscriptionUseCase)
	subscription.RegisterSubscriptionServiceServer(server, subscriptionService)

	// Enable reflection for development tools
	reflection.Register(server)

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Printf("Starting gRPC server on %s:%d", cfg.ServerHost, cfg.ServerPort)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Println("Shutting down server...")
	server.GracefulStop()
	log.Println("Server stopped")
}

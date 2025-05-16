package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/config"
	controllerhttp "github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/controller/http"
	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/handler"
	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/middleware"
	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/service"
	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Run creates objects via constructors
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	l.Info("Starting %s v%s (gateway)", cfg.App.Name, cfg.App.Version)

	// GRPC clients
	svcs, err := service.NewServices(cfg)
	if err != nil {
		l.Fatal(fmt.Errorf("gateway – init services: %w", err))
	}

	// GIN router
	if cfg.HTTP.GIN.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(
		middleware.Logger(l),
		middleware.Recovery(),
		middleware.RateLimit(cfg),
	)

	handlers := handler.NewHandlers(svcs)

	// register routes
	api := router.Group("/api")
	controllerhttp.RegisterPublic(api, handlers)

	protected := api.Group("/")
	protected.Use(middleware.Auth(svcs.Auth))
	controllerhttp.RegisterProtected(protected, handlers)

	// HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	// start server in goroutine
	go func() {
		l.Info("HTTP-gateway listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal(fmt.Errorf("app - Run - gateway - ListenAndServe: %w", err))
		}
	}()

	// graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	sig := <-interrupt
	l.Info("app - Run - gateway - caught signal: %s, shutting down…", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		l.Error("app - Run - gateway - server.Shutdown: %v", err)
	}

	l.Info("app - Run - gateway stopped gracefully")
}

package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Denterry/FinancialAdviser/Backend/auth-service/config"
	grpcController "github.com/Denterry/FinancialAdviser/Backend/auth-service/internal/controller/grpc"
	"github.com/Denterry/FinancialAdviser/Backend/auth-service/internal/repo/persistent"
	"github.com/Denterry/FinancialAdviser/Backend/auth-service/internal/usecase/auth"
	grpcServer "github.com/Denterry/FinancialAdviser/Backend/auth-service/pkg/grpcserver"
	"github.com/Denterry/FinancialAdviser/Backend/auth-service/pkg/logger"
	authv1 "github.com/Denterry/FinancialAdviser/Backend/auth-service/pkg/pb/auth/v1"
	"github.com/Denterry/FinancialAdviser/Backend/auth-service/pkg/postgres"
	"google.golang.org/grpc"
)

// Run creates objects via constructors
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	l.Info("Starting %s v%s", cfg.App.Name, cfg.App.Version)

	// initialize postgres
	pg, err := postgres.New(
		cfg.PG.URL,
		postgres.MaxPoolSize(cfg.PG.PoolMax),
		postgres.MinConns(cfg.PG.MinConns),
		postgres.MaxRetries(cfg.PG.MaxRetries),
		postgres.RetryDelay(cfg.PG.RetryDelay),
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// repository and use case
	userRepository := persistent.NewUserPostgres(pg)
	authUseCase := auth.New(
		userRepository,
		cfg.JWT.Secret,
		cfg.JWT.AccessTTL,
	)

	// GRPC Server
	gs := grpcServer.New(
		grpcServer.Port("0.0.0.0:"+cfg.GRPC.Port),
		grpcServer.MaxStreams(cfg.GRPC.MaxConcurrentStreams),
		grpcServer.TLS(cfg.TLS.CertFile, cfg.TLS.KeyFile),
	)

	// register services
	gs.Serve(func(s *grpc.Server) {
		// register all services with the same server instance
		authv1.RegisterAuthServiceServer(s, grpcController.NewAuthService(authUseCase))
	})
	l.Info("gRPC server listening on " + cfg.GRPC.Port)

	// waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-gs.Notify():
		l.Error("app - Run - grpcserver.Notify: %v", err)
	}

	// graceful shutdown
	l.Info("app - Run - shutting down gRPC server")
	gs.GracefulStop(cfg.GRPC.ShutdownTimeout)
}

package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Denterry/FinancialAdviser/Backend/x-service/config"
	grpcController "github.com/Denterry/FinancialAdviser/Backend/x-service/internal/controller/grpc"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/repo/persistent"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/repo/webapi"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/usecase/admin"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/usecase/tweet"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/pkg/grpcserver"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/pkg/logger"
	adminpb "github.com/Denterry/FinancialAdviser/Backend/x-service/pkg/pb/admin/v1"
	tweetspb "github.com/Denterry/FinancialAdviser/Backend/x-service/pkg/pb/tweets/v1"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/pkg/postgres"
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
		l.Fatal("Failed to initialize postgres: %v", err)
	}
	defer pg.Close()

	// initialize repositories
	tweetRepo := persistent.NewTweetPostgres(pg)

	// scraper / parser
	fetcher, err := webapi.NewSocialFetcher(cfg.XProvider)
	if err != nil {
		l.Fatal("Failed to initialize social fetcher: %v", err)
	}

	// use cases
	tweetUseCase := tweet.New(tweetRepo, fetcher)
	adminUseCase := admin.New(tweetRepo)

	// GRPC server
	gs := grpcserver.New(
		grpcserver.Port("0.0.0.0:"+cfg.GRPC.Port),
		grpcserver.MaxStreams(cfg.GRPC.MaxConcurrentStreams),
		grpcserver.TLS(cfg.TLS.CertFile, cfg.TLS.KeyFile),
	)

	// register services
	gs.Serve(func(s *grpc.Server) {
		// register all services with the same server instance
		tweetspb.RegisterTweetServiceServer(s, grpcController.NewTweetService(tweetUseCase))
		adminpb.RegisterAdminTweetServiceServer(s, grpcController.NewAdminTweetService(adminUseCase))
	})
	l.Info("gRPC server listening on " + cfg.GRPC.Port)

	// waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err := <-gs.Notify():
		l.Error("app - Run - grpcserver.Notify: %v", err)
	}

	// graceful shutdown
	l.Info("app - Run - shutting down gRPC server")
	gs.GracefulStop(cfg.GRPC.ShutdownTimeout)
}

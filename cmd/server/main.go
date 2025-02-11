package main

import (
	"context"
	"github.com/P1ecful/jwt-grpc-test/internal/config"
	"github.com/P1ecful/jwt-grpc-test/internal/grpc/auth"
	"github.com/P1ecful/jwt-grpc-test/internal/service"
	"github.com/P1ecful/jwt-grpc-test/internal/storage/pgx"
	gen "github.com/P1ecful/pkg/gen/grpc/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger, _ := zap.NewDevelopment()
	cfg := config.LoadConfig("config/config.yaml", logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg.Storage.SetURI(logger)
	storage := pgx.NewPGX(logger, cfg.Storage.GetURI())

	if err := storage.Ping(ctx); err != nil {
		logger.Debug("failed to ping to database", zap.Error(err))
	}

	srv := service.NewAuth(logger, storage, cfg.Service)
	logger.Info("service initialized", zap.Any("Config", cfg))

	grpcServer := auth.NewGRPCServer(logger, srv)
	s := grpc.NewServer(grpc.ConnectionTimeout(cfg.GRPC.Timeout))
	gen.RegisterAuthServer(s, grpcServer)
	logger.Info("grpc server initialized", zap.Any("Config", cfg.GRPC))

	l, err := net.Listen("tcp", cfg.GRPC.Port)
	if err != nil {
		logger.Debug("failed to listen", zap.Error(err))
	}

	logger.Info("server started")
	if err := s.Serve(l); err != nil {
		logger.Debug("failed to serve", zap.Error(err))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	s.GracefulStop()
	logger.Info("Gracefully stopped")
}

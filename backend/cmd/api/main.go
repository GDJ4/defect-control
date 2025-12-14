package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"defect-tracker/internal/pkg/auth"
	"defect-tracker/internal/pkg/config"
	"defect-tracker/internal/pkg/logger"
	"defect-tracker/internal/pkg/server"
	"defect-tracker/internal/pkg/storage"
	"defect-tracker/internal/repo/postgres"
	"defect-tracker/internal/service/defect"
	"defect-tracker/internal/service/project"
	"defect-tracker/internal/service/token"
	"defect-tracker/internal/service/user"
	transporthttp "defect-tracker/internal/transport/http"
	"defect-tracker/internal/transport/http/handlers"
	"defect-tracker/internal/transport/http/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log, err := logger.New(cfg.AppEnv)
	if err != nil {
		panic(err)
	}
	defer log.Sync() //nolint:errcheck

	pool := initDatabase(log, cfg.Database.DSN)
	defer pool.Close()

	var fileStorage storage.Provider
	switch cfg.Storage.Driver {
	case "s3":
		fileStorage, err = storage.NewS3(storage.S3Config{
			Endpoint:   cfg.Storage.S3.Endpoint,
			AccessKey:  cfg.Storage.S3.AccessKey,
			SecretKey:  cfg.Storage.S3.SecretKey,
			Bucket:     cfg.Storage.S3.Bucket,
			Region:     cfg.Storage.S3.Region,
			UseSSL:     cfg.Storage.S3.UseSSL,
			PresignTTL: cfg.Storage.S3.PresignTTL,
		})
	default:
		fileStorage, err = storage.NewLocal(cfg.Storage.Path)
	}
	if err != nil {
		log.Fatal("failed to init storage", zap.Error(err))
	}

	userRepo := postgres.NewUserRepository(pool)
	userService := user.NewService(userRepo)
	tokenManager := auth.NewManager(cfg.Auth.Secret, cfg.Auth.AccessTTL)
	tokenRepo := postgres.NewTokenRepository(pool)
	tokenService := token.NewService(tokenRepo, cfg.Auth.RefreshTTL)

	defectRepo := postgres.NewDefectRepository(pool)
	defectService := defect.NewService(defectRepo)
	defectHandler := handlers.NewDefectHandler(defectService, fileStorage)

	projectRepo := postgres.NewProjectRepository(pool)
	projectService := project.NewService(projectRepo)
	projectHandler := handlers.NewProjectHandler(projectService)

	authHandler := handlers.NewAuthHandler(userService, tokenService, tokenManager)
	authMiddleware := middleware.NewAuthMiddleware(tokenManager, userService)

	router := transporthttp.NewRouter(cfg.AppName, authHandler, authMiddleware, defectHandler, projectHandler)
	httpServer := server.NewHTTPServer(cfg, router, log)

	go func() {
		if err := httpServer.Start(); err != nil {
			log.Fatal("server stopped unexpectedly", zap.Error(err))
		}
	}()

	waitForShutdown(log, httpServer)
}

func waitForShutdown(log *zap.Logger, srv *server.HTTPServer) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Info("received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
	}
}

func initDatabase(log *zap.Logger, dsn string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("failed to init postgres pool", zap.Error(err))
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("postgres not ready", zap.Error(err))
	}

	return pool
}

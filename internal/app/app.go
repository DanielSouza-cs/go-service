package app

import (
	"context"
	"go-service/internal/auth"
	"go-service/internal/config"
	"go-service/internal/httpserver"
	"go-service/internal/logger"
	"go-service/internal/student"

	"go.uber.org/zap"
)

type App struct {
	Config *config.Config
	Logger *zap.Logger
	Server *httpserver.Server
}

func New(ctx context.Context) (*App, error) {
	cfg := config.Load()

	log, err := logger.New(cfg.LogLevel, cfg.Environment)
	if err != nil {
		return nil, err
	}
	log.Info("starting go-service")

	authClient := auth.New(cfg, log)

	if err := authClient.Login(ctx); err != nil {
		log.Warn("initial authentication with Node.js API failed. will retry on first request", zap.Error(err))
	} else {
		log.Info("initial authentication successful")
	}

	studentClient := student.NewClient(authClient, cfg, log)
	studentSvc := student.NewService(studentClient, log)

	router := httpserver.NewRouter(studentSvc, log)
	srv := httpserver.NewServer(cfg.Port, router)

	return &App{
		Config: cfg,
		Logger: log,
		Server: srv,
	}, nil
}

package main

import (
	"context"
	"fmt"
	"go-service/internal/app"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	application, err := app.New(ctx)
	if err != nil {
		panic("failed to initialize app: " + err.Error())
	}

	go func() {
		application.Logger.Info("HTTP server starting", zap.String("port", application.Config.Port))
		if err := application.Server.Start(); err != nil && err != http.ErrServerClosed {
			application.Logger.Fatal("http server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	application.Logger.Info("shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if err := application.Server.Shutdown(ctxShutdown); err != nil {
		application.Logger.Error("server shutdown failed", zap.Error(err))
	}

	if err := application.Logger.Sync(); err != nil {
		fmt.Println("failed to sync logger on exit:", err)
	}
}

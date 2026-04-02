package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"omniport-api/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	setupLogger(cfg)
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "rls-service", "mode": cfg.App.Mode})
	})

	serve(cfg, "rls-service", cfg.App.PortFor("RLS"), r)
}

func serve(cfg *config.Config, service string, port string, handler http.Handler) {
	addr := fmt.Sprintf(":%s", port)
	srv := &http.Server{Addr: addr, Handler: handler}

	go func() {
		slog.Info("Server running", "service", service, "port", port, "env", cfg.App.Env, "mode", cfg.App.Mode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...", "service", service)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func setupLogger(cfg *config.Config) {
	logLevel := parseLogLevel(cfg.App.LogLevel)
	var logHandler slog.Handler
	if cfg.App.Env == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	}
	slog.SetDefault(slog.New(logHandler))
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "VERBOSE", "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

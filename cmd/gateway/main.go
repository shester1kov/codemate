package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shester1kov/codemate/internal/config"
	"github.com/shester1kov/codemate/internal/gateway/router"
	"github.com/shester1kov/codemate/internal/logger"
	"go.uber.org/zap"
)

func main() {

	// загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// инициализируем логгер
	log, err := logger.New(
		cfg.Logger.Level,
		cfg.Logger.Encoding,
		cfg.Logger.OutputPath,
	)

	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))

	}

	defer log.Sync()

	log.Info("Starting CodeMate Gateway",
		zap.String("version", "1.0.0"),
		zap.String("mode", cfg.Server.Mode),
	)

	r := router.Setup(log, cfg.Server.Mode)

	// настраиваем http сервер

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
		// таймауты
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// запускаем сервер в отдельной горутине
	go func() {
		log.Info("Server listening",
			zap.String("address", addr),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed",
				zap.Error(err),
			)
		}
	}()

	// ждем сигнал завершения
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown",
			zap.Error(err),
		)
	}

	log.Info("Server stopped gracefully")
}

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/nkolesnikov999/micro2-OK/notification/internal/app"
	"github.com/nkolesnikov999/micro2-OK/notification/internal/config"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/closer"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

const configPath = "./deploy/compose/notification/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Ошибка при работе приложения", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}

	// Затем аккуратно закрываем логгер
	if err := logger.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "logger sync error: %v\n", err)
	}
	if err := logger.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "logger close error: %v\n", err)
	}
}

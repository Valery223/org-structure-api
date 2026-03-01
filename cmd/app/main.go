package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Valery223/org-structure-api/internal/app"
	"github.com/Valery223/org-structure-api/internal/config"
	handler "github.com/Valery223/org-structure-api/internal/handler/http"
	"github.com/Valery223/org-structure-api/internal/repository/postgres"
	"github.com/Valery223/org-structure-api/internal/service"
	storage "github.com/Valery223/org-structure-api/internal/storage/postgres"
)

func main() {
	// 1. Инициализация конфига
	cfg := config.MustLoad()

	// 2. Инициализация логгера
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))

	// 3. Подключение к БД
	db, err := storage.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.String("err", err.Error()))
		os.Exit(1)
	}
	log.Info("database connection established")

	// 4. Инициализация слоев
	// Repositories
	deptRepo := postgres.NewDepartmentRepository(db)
	empRepo := postgres.NewEmployeeRepository(db)

	// Services
	deptService := service.NewDepartmentService(deptRepo, empRepo)
	empService := service.NewEmployeeService(empRepo, deptRepo)

	// Handlers
	deptHandler := handler.NewDepartmentHandler(deptService)
	empHandler := handler.NewEmployeeHandler(empService)

	// 5. Роутер
	router := app.SetupRouter(deptHandler, empHandler)

	// 6. Запуск сервера с Graceful Shutdown
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Запускаем сервер в горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", slog.String("err", err.Error()))
		}
	}()

	log.Info("server started", slog.String("address", cfg.HTTPServer.Address))

	// Ждем сигнал остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("stopping server")

	// Даем серверу 5 секунд на завершение текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.String("err", err.Error()))
		return
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

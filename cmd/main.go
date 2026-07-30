package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"EventBooker/internal/config"
	"EventBooker/internal/handler"
	"EventBooker/internal/repository"
	"EventBooker/internal/service"
	"EventBooker/internal/worker"

	pgxdriver "github.com/wb-go/wbf/dbpg/pgx-driver"
	"github.com/wb-go/wbf/logger"
)

func main() {
	cfg, err := config.New("./config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка загрузки конфигурации: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.InitLogger(
		logger.ZapEngine,
		"EventBooker",
		cfg.Env,
		logger.WithLevel(logger.InfoLevel),
		logger.WithRotation("logs/app.log", 100, 5, 30),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка инициализации логгера: %v\n", err)
		os.Exit(1)
	}

	log.Info("Запуск приложения EventBooker...")

	pg, err := pgxdriver.New(
		cfg.DSN,
		log,
		pgxdriver.MaxPoolSize(50),
		pgxdriver.MaxConnAttempts(5),
		pgxdriver.BaseRetryDelay(100*time.Millisecond),
	)
	if err != nil {
		log.Error("Не удалось подключиться к PostgreSQL", "error", err)
		os.Exit(1)
	}
	defer pg.Close()

	ctx := context.Background()
	if err := pg.Ping(ctx); err != nil {
		log.Error("PostgreSQL недоступен", "error", err)
		os.Exit(1)
	}
	log.Info("PostgreSQL успешно подключен")

	repo := repository.New(pg)
	svc := service.New(repo)
	h := handler.New(svc)

	router := h.InitRoutes(cfg.Env)

	cleanerWorker := worker.NewExpiredCleaner(svc, 30*time.Second, log)
	cleanerWorker.Start()

	srv := &http.Server{
		Addr:    cfg.ServeAddr,
		Handler: router,
	}

	go func() {
		log.Info(fmt.Sprintf("HTTP-сервер запущен на %s", cfg.ServeAddr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Критическая ошибка HTTP-сервера", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Получен сигнал на остановку, завершаем работу...")

	cleanerWorker.Stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Принудительная остановка HTTP-сервера из-за ошибки", "error", err)
	}

	log.Info("Приложение успешно и безопасно остановлено")
}

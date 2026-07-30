package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/wb-go/wbf/logger"
)

// BookerService — интерфейс сервисного слоя с методом очистки
type BookerService interface {
	ProcessExpiredBookings(ctx context.Context) (int64, error)
}

type ExpiredCleaner struct {
	service     BookerService
	interval    time.Duration
	taskTimeout time.Duration
	stopCh      chan struct{}
	wg          sync.WaitGroup
	logger      logger.Logger
}

func NewExpiredCleaner(service BookerService, interval time.Duration, logger logger.Logger) *ExpiredCleaner {
	return &ExpiredCleaner{
		service:     service,
		interval:    interval,
		taskTimeout: 10 * time.Second,
		stopCh:      make(chan struct{}),
		logger:      logger,
	}
}

// Start запускает воркер в отдельной горутине
func (w *ExpiredCleaner) Start() {
	w.wg.Add(1)
	go w.run()
	w.logger.Info("Expired bookings worker запущен", "interval", w.interval)
}

func (w *ExpiredCleaner) run() {
	defer w.wg.Done()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.process()

	for {
		select {
		case <-ticker.C:
			w.process()
		case <-w.stopCh:
			w.logger.Info("Expired bookings worker получил стоп сигнал, finishing...")
			return
		}
	}
}

func (w *ExpiredCleaner) process() {
	ctx, cancel := context.WithTimeout(context.Background(), w.taskTimeout)
	defer cancel()

	count, err := w.service.ProcessExpiredBookings(ctx)
	if err != nil {
		w.logger.Error("Failed to process expired bookings", slog.String("error", err.Error()))
		return
	}

	if count > 0 {
		w.logger.Info("Successfully processed expired bookings", slog.Int64("cancelled_count", count))
	}
}

// Stop плавно останавливает воркер и ждет завершения текущей итерации
func (w *ExpiredCleaner) Stop() {
	close(w.stopCh)
	w.wg.Wait()
	w.logger.Info("Expired bookings worker остановлен мягко")
}

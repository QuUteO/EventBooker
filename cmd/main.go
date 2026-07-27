package main

import (
	"EventBooker/internal/config"
	"fmt"
	"os"
	"time"

	"github.com/wb-go/wbf/dbpg/pgx-driver"
	"github.com/wb-go/wbf/logger"
)

func main() {
	cfg, err := config.New("./config.yaml")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
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

	pg, err := pgxdriver.New(
		cfg.DSN,
		log,
		pgxdriver.MaxPoolSize(50),
		pgxdriver.MaxConnAttempts(5),
		pgxdriver.BaseRetryDelay(100*time.Millisecond),
	)
	if err != nil {
		log.Error("Failed to connect to PostgreSQL:", err)
	}
	defer pg.Close()

}

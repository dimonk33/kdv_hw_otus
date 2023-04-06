package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig()
	if err != nil {
		fmt.Println("Ошибка инициализации конфигуратора: " + err.Error())
		os.Exit(1)
	}
	logg := logger.New(config.Logger.Level)

	var storage app.Storage
	switch config.GetStorageType() {
	case StorageInMemory:
		storage = memorystorage.New()
	case StorageDB:
		storage = sqlstorage.New(config.GetDBURL(), logg)
	default:
		logg.Error("Неподдерживаемый тип хранилища")
		os.Exit(1)
	}
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(config.GetServerAddr(), logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}
}

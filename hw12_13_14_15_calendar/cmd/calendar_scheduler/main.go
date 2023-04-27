package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/logger"
	kafkaapp "github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/queue/kafka"
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/scheduler"
	memorystorage "github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/scheduler_config.toml", "Path to configuration file")
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

	var storage scheduler.Storage
	switch config.GetStorageType() {
	case StorageInMemory:
		storage = memorystorage.New()
	case StorageDB:
		storage = sqlstorage.New(config.GetDBURL(), logg)
	default:
		logg.Error("Неподдерживаемый тип хранилища")
		os.Exit(1)
	}

	sender := kafkaapp.NewProducer(config.GetBroker(), config.GetTopic(), logg)

	notifyHour, notifyMin := config.GetNotifyTime()

	schlr := scheduler.NewScheduler(storage, sender, scheduler.NotifyTime{H: notifyHour, M: notifyMin}, logg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	schlr.Start(ctx)

	logg.Info("планировщик запущен...")

	<-ctx.Done()
}

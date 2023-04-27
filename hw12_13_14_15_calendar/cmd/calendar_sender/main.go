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
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/sender"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/sender_config.toml", "Path to configuration file")
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

	receiver := kafkaapp.NewConsumer(config.GetBroker(), config.GetTopic(), logg)

	notifier := sender.NewNotifier(os.Stdout)

	sndr := sender.NewSender(receiver, notifier, logg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	sndr.Start(ctx)

	logg.Info("планировщик запущен...")

	<-ctx.Done()
}

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	internalgrpc "github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/server/grpc"

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

	grpcServer := internalgrpc.NewServer(config.GetHttpServerAddr(), logg, calendar)
	httpServer := internalhttp.NewServer(config.GetGrpcServerAddr(), logg, grpcServer)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		if err := httpServer.Start(ctx); err != nil {
			logg.Error("failed to start httpServer: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	go func() {
		if err := grpcServer.Start(ctx); err != nil {
			logg.Error("failed to start grpcServer: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := grpcServer.Stop(ctx); err != nil {
			logg.Error("failed to stop grpcServer: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	<-ctx.Done()
}

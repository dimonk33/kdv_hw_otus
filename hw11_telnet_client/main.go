package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"time"
)

var argTimeout, host, port string

func init() {
	flag.StringVar(&argTimeout, "timeout", "10s", "connect timeout")
}

func main() {
	flag.Parse()
	for i, val := range flag.Args() {
		switch i {
		case 0:
			host = val
		case 1:
			port = val
		}
	}
	address := net.JoinHostPort(host, port)
	timeout, err := time.ParseDuration(argTimeout)
	if err != nil {
		println("неверное значение таймаута: %w", err)
		os.Exit(1)
	}

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err = client.Connect(); err != nil {
		println("%w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer exit(client, cancel)

	go func() {
	OUTER:
		for {
			select {
			case <-ctx.Done():
				break OUTER
			default:
				if err := client.Send(); err != nil {
					printExitMessage(err)
					break OUTER
				}
			}
		}
		cancel()
	}()

	go func() {
	OUTER:
		for {
			select {
			case <-ctx.Done():
				break OUTER
			default:
				if err := client.Receive(); err != nil {
					break OUTER
				}
			}
		}
		cancel()
	}()

	<-ctx.Done()
}

func exit(client TelnetClient, cancel context.CancelFunc) {
	cancel()
	if err := client.Close(); err != nil {
		println(err)
	}
}

func printExitMessage(err error) {
	if errors.Is(err, io.EOF) {
		fmt.Printf("Выход из программы ...%s", err)
	}
}

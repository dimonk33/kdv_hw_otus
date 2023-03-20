package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	fmt.Printf("%v:%v timeout=%v\n", host, port, argTimeout)
	timeout, err := time.ParseDuration(argTimeout)
	if err != nil {
		println("неверное значение таймаута: %w", err)
		os.Exit(1)
	}

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err = client.Connect(); err != nil {
		println("%w", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGQUIT)
	defer exit(client, cancel)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
	OUTER:
		for {
			select {
			case <-ctx.Done():
				break OUTER
			default:
				if err := client.Send(); err != nil {
					break OUTER
				}
			}
		}
		cancel()
	}()

	go func() {
		defer wg.Done()
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

	wg.Wait()
}

func exit(client TelnetClient, cancel context.CancelFunc) {
	cancel()
	if err := client.Close(); err != nil {
		println(err)
	}
}

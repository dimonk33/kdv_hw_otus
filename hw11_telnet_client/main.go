package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
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
	timeout, err := time.ParseDuration(argTimeout)
	if err != nil {
		println("неверное значение таймаута: %w", err)
		os.Exit(1)
	}

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err = client.Connect(); err != nil {
		println("%w", err)
	}

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	exitChan := make(chan int)

	defer func() {
		if err := client.Close(); err != nil {
			println(err)
		}
		close(signalChanel)
	}()

	go func() {
		s, ok := <-signalChanel
		if !ok {
			return
		}
		if s == syscall.SIGQUIT {
			fmt.Println(".....EOF")
		}
		exitChan <- 0
	}()

	go func() {
		processClient(client.Send, exitChan)
	}()

	go func() {
		processClient(client.Receive, exitChan)
	}()
	<-exitChan
}

func processClient(processor func() error, exitCh chan int) {
	for {
		select {
		case <-exitCh:
			return
		default:
			if err := processor(); err != nil {
				time.Sleep(100 * time.Millisecond)
				exitCh <- 2
				return
			}
		}
	}
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	ErrArgsNumber = errors.New("error arguments number")
	timeout       time.Duration
)

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout for connecting to the server")
}

func main() {
	flag.Parse()

	args := os.Args
	argsLen := len(os.Args)

	if argsLen < 3 || argsLen > 4 {
		fmt.Println(ErrArgsNumber)
		return
	}

	host := args[argsLen-2]
	port := args[argsLen-1]

	l, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Println(err)
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		if err := client.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		client.Receive()
		wg.Done()
	}()

	go func() {
		client.Send()
		wg.Done()
	}()

	wg.Wait()
}

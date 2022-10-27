package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"path"
	"sync"
	"syscall"

	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/broker"
	config "github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/config/sender"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/logger"
)

var configPath string

func init() {
	defaultConfigPath := path.Join("configs", "sender", "config.toml")
	flag.StringVar(&configPath, "config", defaultConfigPath, "Path to configuration file")
}

func main() {
	flag.Parse()

	conf, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	log := logger.New(conf.Logger)
	b := broker.New(conf.Connection)

	if err := b.Start(); err != nil {
		log.Error(err.Error())
		return
	}

	defer func() {
		if err := b.Stop(); err != nil {
			log.Error(err.Error())
			return
		}
	}()

	if err := b.QueueDeclare(conf.Queue); err != nil {
		return
	}

	msgs, err := b.Consume(conf.Consume)
	if err != nil {
		log.Error(err.Error())
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	log.Info("Waiting for messages. To exit press CTRL+C")

	go func() {
		for {
			select {
			case m := <-msgs:
				log.Info(fmt.Sprintf("Received a message: %s", m.Body))
			case <-ctx.Done():
				wg.Done()
				return
			}
		}
	}()

	wg.Wait()
}

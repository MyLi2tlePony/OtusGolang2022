package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/app"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/broker"
	config "github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/config/scheduler"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/storage/sql"
)

var (
	ErrInvalidStorageType = errors.New("invalid storage type")

	configPath  string
	storageType string
)

func init() {
	defaultConfigPath := path.Join("configs", "scheduler", "config.toml")
	flag.StringVar(&configPath, "config", defaultConfigPath, "Path to configuration file")

	flag.StringVar(&storageType, "storage", "sql", "Type of storage. Expected values: \"mem\" || \"sql\"")
}

func main() {
	flag.Parse()

	conf, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	log := logger.New(conf.Logger)

	var application *app.Calendar

	switch storageType {
	case "mem":
		s := memorystorage.New()
		application = app.New(s)
	case "sql":
		dbConf := conf.Database
		connString := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
			dbConf.Prefix, dbConf.UserName, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.DatabaseName)

		s := sqlstorage.New(connString)
		application = app.New(s)
	default:
		log.Error(ErrInvalidStorageType.Error())
		os.Exit(1)
	}

	mesBroker := broker.New(conf.Connection)

	if err := mesBroker.Start(); err != nil {
		log.Error(err.Error())
		return
	}

	if err := mesBroker.QueueDeclare(conf.Queue); err != nil {
		return
	}

	defer func() {
		if err := mesBroker.Stop(); err != nil {
			log.Error(err.Error())
			return
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		ticker := time.Tick(1 * time.Minute)
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return
			case <-ticker:
				schedule(ctx, application, log, mesBroker, conf.Publish)
			}
		}
	}()

	log.Info("scheduler is running...")
	wg.Wait()
}

func schedule(ctx context.Context, application *app.Calendar, log *logger.Logger, mesBroker broker.Broker, conf broker.PublishConfig) { //nolint:lll
	now := time.Now()
	notificationTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)
	currentEvents, err := application.SelectEventsByTime(ctx, notificationTime)
	if err != nil {
		log.Error(err.Error())
	}

	for _, event := range currentEvents {
		marshal, err := json.Marshal(event)
		if err != nil {
			return
		}

		ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
		err = mesBroker.PublishWithContext(ctxTimeout, conf, marshal)
		cancelTimeout()
		if err != nil {
			log.Error(err.Error())
		}

		log.Info(fmt.Sprintf("message send: %s", marshal))
	}

	notificationTime = time.Date(now.Year()-1, now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)
	oldEvents, err := application.SelectEventsByTime(ctx, notificationTime)
	if err != nil {
		log.Error(err.Error())
	}

	for _, event := range oldEvents {
		err := application.DeleteEvent(ctx, event.GetID())
		if err != nil {
			log.Error(err.Error())
			continue
		}

		log.Info(fmt.Sprintf("user with id: %s deleted", event.GetID()))
	}
}

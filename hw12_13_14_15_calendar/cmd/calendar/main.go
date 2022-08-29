package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/app"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/config"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/logger"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/server"
	internalgrpc "github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/storage/sql"
	"google.golang.org/grpc"
)

var (
	ErrInvalidStorageType = errors.New("invalid storage type")
	ErrInvalidServerType  = errors.New("invalid server type")

	configPath  string
	storageType string
	serverType  string
)

func init() {
	defaultConfigPath := path.Join("configs", "config.toml")
	flag.StringVar(&configPath, "config", defaultConfigPath, "Path to configuration file")

	flag.StringVar(&storageType, "storage", "sql", "Type of storage. Expected values: \"mem\" || \"sql\"")
	flag.StringVar(&serverType, "server", "grpc", "Type of server. Expected values: \"http\" || \"grpc\"")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	conf, err := config.New(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	log := logger.New(*conf.Logger)

	var application *app.Calendar

	switch storageType {
	case "mem":
		storage := memorystorage.New()
		application = app.New(storage)
	case "sql":
		storage := sqlstorage.New(*conf.Database)
		application = app.New(storage)
	default:
		log.Error(ErrInvalidStorageType.Error())
		os.Exit(1)
	}

	var serv server.Server

	switch serverType {
	case "http":
		serv = internalhttp.NewServer(log, application, conf.Server)
	case "grpc":
		serv = internalgrpc.NewServer(log, application, conf.Server)
	default:
		log.Error(ErrInvalidServerType.Error())
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		if err := serv.Stop(); err != nil {
			log.Error("failed to stop serv: " + err.Error())
		}
	}()

	log.Info("app is running...")

	err = serv.Start()
	if !errors.Is(err, grpc.ErrServerStopped) && !errors.Is(err, http.ErrServerClosed) && err != nil {
		log.Error("failed to start serv: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	log.Info("serv closed")
}

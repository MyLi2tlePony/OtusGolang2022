package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/config"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/server/dto"
)

type Server struct {
	app    Application
	logger Logger

	srv *http.Server
}

type Logger interface {
	Fatal(string)
	Error(string)
	Warn(string)
	Info(string)
	Debug(string)
	Trace(string)
}

type Application interface {
	CreateUser(ctx context.Context, dtoUser dto.User) error
	SelectUsers(ctx context.Context) ([]dto.User, error)
	DeleteUser(ctx context.Context, id string) error

	CreateEvent(ctx context.Context, dtoEvent dto.Event) error
	SelectEvents(ctx context.Context) ([]dto.Event, error)
	UpdateEvent(ctx context.Context, dtoEvent dto.Event) error
	DeleteEvent(ctx context.Context, id string) error
}

func NewServer(logger Logger, app Application, config config.ServerConfig) *Server {
	handler := newHandler(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", handler.getHello)

	middleware := newMiddleware(logger, mux)
	middleware.logging()

	return &Server{
		logger: logger,
		app:    app,
		srv: &http.Server{
			Addr:    net.JoinHostPort(config.Host, config.Port),
			Handler: middleware.Handler,
		},
	}
}

func (s *Server) Start() error {
	s.logger.Info(fmt.Sprintf("server listening: %s", s.srv.Addr))

	if err := s.srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	if err := s.srv.Close(); err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}

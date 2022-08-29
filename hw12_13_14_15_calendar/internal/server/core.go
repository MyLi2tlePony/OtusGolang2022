package server

import (
	"context"
	"time"
)

type Server interface {
	Start() error
	Stop() error
}

type Config interface {
	GetPort() string
	GetHost() string
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
	CreateUser(context.Context, User) error
	SelectUsers(context.Context) ([]User, error)
	DeleteUser(context.Context, string) error

	CreateEvent(context.Context, Event) error
	SelectEvents(context.Context) ([]Event, error)
	UpdateEvent(context.Context, Event) error
	DeleteEvent(context.Context, string) error
}

type User interface {
	GetID() string

	GetFirstName() string
	GetLastName() string

	GetEmail() string
	GetAge() int64
}

type Event interface {
	GetID() string

	GetTitle() string
	GetDescription() string

	GetBeginning() time.Time
	GetFinish() time.Time
	GetNotification() time.Time

	GetUserID() string
}

package app

import (
	"context"
	"sync"
	"time"

	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/server"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/storage"
)

type Calendar struct {
	storage Storage
	mutex   sync.RWMutex
}

type Storage interface {
	CreateUser(context.Context, storage.User) error
	SelectUsers(ctx context.Context) ([]storage.User, error)
	DeleteUser(context.Context, string) error

	CreateEvent(context.Context, storage.Event) error
	SelectEvents(context.Context) ([]storage.Event, error)
	SelectEventsByTime(context.Context, time.Time) ([]storage.Event, error)
	UpdateEvent(context.Context, storage.Event) error
	DeleteEvent(context.Context, string) error
}

func New(storage Storage) *Calendar {
	return &Calendar{
		storage: storage,
		mutex:   sync.RWMutex{},
	}
}

func (calendar *Calendar) CreateUser(ctx context.Context, user server.User) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	storageUser := storage.User{
		ID:        user.GetID(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Email:     user.GetEmail(),
		Age:       user.GetAge(),
	}

	return calendar.storage.CreateUser(ctx, storageUser)
}

func (calendar *Calendar) SelectUsers(ctx context.Context) ([]server.User, error) {
	calendar.mutex.RLock()
	defer calendar.mutex.RUnlock()

	users := make([]server.User, 0)

	storageUsers, err := calendar.storage.SelectUsers(ctx)
	if err != nil {
		return users, err
	}

	for _, storageUser := range storageUsers {
		user := storageUser
		users = append(users, &user)
	}

	return users, nil
}

func (calendar *Calendar) DeleteUser(ctx context.Context, id string) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	return calendar.storage.DeleteUser(ctx, id)
}

func (calendar *Calendar) CreateEvent(ctx context.Context, event server.Event) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	storageEvent := storage.Event{
		ID:           event.GetID(),
		Title:        event.GetTitle(),
		Description:  event.GetDescription(),
		UserID:       event.GetUserID(),
		Beginning:    event.GetBeginning(),
		Finish:       event.GetFinish(),
		Notification: event.GetNotification(),
	}

	return calendar.storage.CreateEvent(ctx, storageEvent)
}

func (calendar *Calendar) UpdateEvent(ctx context.Context, event server.Event) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	storageEvent := storage.Event{
		ID:           event.GetID(),
		Title:        event.GetTitle(),
		Description:  event.GetDescription(),
		UserID:       event.GetUserID(),
		Beginning:    event.GetBeginning(),
		Finish:       event.GetFinish(),
		Notification: event.GetNotification(),
	}

	return calendar.storage.UpdateEvent(ctx, storageEvent)
}

func (calendar *Calendar) DeleteEvent(ctx context.Context, id string) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	return calendar.storage.DeleteEvent(ctx, id)
}

func (calendar *Calendar) SelectEvents(ctx context.Context) ([]server.Event, error) {
	calendar.mutex.RLock()
	defer calendar.mutex.RUnlock()

	events := make([]server.Event, 0)

	storageEvents, err := calendar.storage.SelectEvents(ctx)
	if err != nil {
		return events, err
	}

	for _, storageEvent := range storageEvents {
		event := storageEvent
		events = append(events, &event)
	}

	return events, nil
}

func (calendar *Calendar) SelectEventsByTime(ctx context.Context, t time.Time) ([]server.Event, error) {
	calendar.mutex.RLock()
	defer calendar.mutex.RUnlock()

	events := make([]server.Event, 0)

	storageEvents, err := calendar.storage.SelectEventsByTime(ctx, t)
	if err != nil {
		return events, err
	}

	for _, storageEvent := range storageEvents {
		event := storageEvent
		events = append(events, &event)
	}

	return events, nil
}

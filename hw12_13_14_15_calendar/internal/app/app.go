package app

import (
	"context"
	"sync"

	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/server/dto"
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
	UpdateEvent(context.Context, storage.Event) error
	DeleteEvent(context.Context, string) error
}

func New(storage Storage) *Calendar {
	return &Calendar{
		storage: storage,
		mutex:   sync.RWMutex{},
	}
}

func (calendar *Calendar) CreateUser(ctx context.Context, dtoUser dto.User) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	user := storage.User{
		ID:        dtoUser.ID,
		FirstName: dtoUser.FirstName,
		LastName:  dtoUser.LastName,
		Email:     dtoUser.Email,
		Age:       dtoUser.Age,
	}

	return calendar.storage.CreateUser(ctx, user)
}

func (calendar *Calendar) SelectUsers(ctx context.Context) ([]dto.User, error) {
	calendar.mutex.RLock()
	defer calendar.mutex.RUnlock()

	users := make([]dto.User, 0)

	storageUsers, err := calendar.storage.SelectUsers(ctx)
	if err != nil {
		return users, err
	}

	for _, storageUser := range storageUsers {
		users = append(users, dto.User{
			ID:        storageUser.ID,
			FirstName: storageUser.FirstName,
			LastName:  storageUser.LastName,
			Email:     storageUser.Email,
			Age:       storageUser.Age,
		})
	}

	return users, nil
}

func (calendar *Calendar) DeleteUser(ctx context.Context, id string) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	return calendar.storage.DeleteUser(ctx, id)
}

func (calendar *Calendar) CreateEvent(ctx context.Context, dtoEvent dto.Event) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	event := storage.Event{
		ID:           dtoEvent.ID,
		Title:        dtoEvent.Title,
		Description:  dtoEvent.Description,
		UserID:       dtoEvent.UserID,
		Beginning:    dtoEvent.Beginning,
		Finish:       dtoEvent.Finish,
		Notification: dtoEvent.Notification,
	}

	return calendar.storage.CreateEvent(ctx, event)
}

func (calendar *Calendar) UpdateEvent(ctx context.Context, dtoEvent dto.Event) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	event := storage.Event{
		ID:           dtoEvent.ID,
		Title:        dtoEvent.Title,
		Description:  dtoEvent.Description,
		UserID:       dtoEvent.UserID,
		Beginning:    dtoEvent.Beginning,
		Finish:       dtoEvent.Finish,
		Notification: dtoEvent.Notification,
	}

	return calendar.storage.UpdateEvent(ctx, event)
}

func (calendar *Calendar) DeleteEvent(ctx context.Context, id string) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	return calendar.storage.DeleteEvent(ctx, id)
}

func (calendar *Calendar) SelectEvents(ctx context.Context) ([]dto.Event, error) {
	calendar.mutex.RLock()
	defer calendar.mutex.RUnlock()

	events := make([]dto.Event, 0)

	storageEvents, err := calendar.storage.SelectEvents(ctx)
	if err != nil {
		return events, err
	}

	for _, storageEvent := range storageEvents {
		events = append(events, dto.Event{
			ID:           storageEvent.ID,
			Title:        storageEvent.Title,
			Description:  storageEvent.Description,
			Beginning:    storageEvent.Beginning,
			Notification: storageEvent.Notification,
			UserID:       storageEvent.UserID,
		})
	}

	return events, nil
}

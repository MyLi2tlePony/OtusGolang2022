package memorystorage

import (
	"context"
	"errors"
	"time"

	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	events map[string]storage.Event
	users  map[string]storage.User
}

var ErrEventNotFound = errors.New("event not found")

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event, 0),
		users:  make(map[string]storage.User, 0),
	}
}

func (s *Storage) CreateUser(ctx context.Context, user storage.User) error {
	user.ID = uuid.New().String()
	s.users[user.ID] = user

	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, userID string) error {
	delete(s.users, userID)
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	event.ID = uuid.New().String()
	s.events[event.ID] = event

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	delete(s.events, eventID)
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	if _, ok := s.events[event.ID]; !ok {
		return ErrEventNotFound
	}

	s.events[event.ID] = event
	return nil
}

func (s *Storage) SelectEvents(ctx context.Context) ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	for _, event := range s.events {
		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) SelectEventsByTime(ctx context.Context, t time.Time) ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	for _, event := range s.events {
		if event.Notification.Equal(t) {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *Storage) SelectUsers(ctx context.Context) ([]storage.User, error) {
	users := make([]storage.User, 0)

	for _, user := range s.users {
		users = append(users, user)
	}

	return users, nil
}

//go:build integration
// +build integration

package sqlstorage

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/config/calendar"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	conf := calendar.DatabaseConfig{
		Prefix:       "postgresql",
		DatabaseName: "calendardb",
		Host:         "localhost",
		Port:         "5436",
		UserName:     "postgres",
		Password:     "1234512345",
	}
	mutex := sync.Mutex{}

	t.Run("user case", func(t *testing.T) {
		mutex.Lock()
		defer mutex.Unlock()

		s := New(conf)
		ctx := context.Background()

		users := []storage.User{
			{
				FirstName: "Andrey",
				LastName:  "Shabarov",
				Email:     "shabandrew@mail.ru",
				Age:       21,
			},
			{
				FirstName: "Tom",
				LastName:  "Skhot",
				Email:     "tomangry@skhot.ru",
				Age:       11,
			},
		}

		for _, user := range users {
			require.Nil(t, s.CreateUser(ctx, user))
		}

		selectedUsers, err := s.SelectUsers(ctx)
		require.Nil(t, err)

		for _, selectedUser := range selectedUsers {
			require.True(t, containsUser(users, selectedUser))
			require.Nil(t, s.DeleteUser(ctx, selectedUser.ID))
		}

		selectedUsers, err = s.SelectUsers(ctx)
		require.Nil(t, err)
		require.Len(t, selectedUsers, 0)
	})

	t.Run("event case", func(t *testing.T) {
		mutex.Lock()
		defer mutex.Unlock()

		s := New(conf)
		ctx := context.Background()

		user := storage.User{
			FirstName: "Andrey",
			LastName:  "Shabarov",
			Email:     "shabandrew@mail.ru",
			Age:       21,
		}

		require.Nil(t, s.CreateUser(ctx, user))
		selectedUsers, err := s.SelectUsers(ctx)
		require.Nil(t, err)

		user = selectedUsers[0]

		events := []storage.Event{
			{
				Title:        "Встреча с другом",
				Description:  "",
				Beginning:    time.Date(2022, time.August, 22, 18, 0, 0, 0, time.UTC),
				Finish:       time.Date(2022, time.August, 22, 20, 0, 0, 0, time.UTC),
				Notification: time.Date(2022, time.August, 22, 13, 0, 0, 0, time.UTC),
				UserID:       user.ID,
			},
			{
				Title:        "Конференция",
				Description:  "Сделать презентацию",
				Beginning:    time.Date(2022, time.August, 22, 20, 0, 0, 0, time.UTC),
				Finish:       time.Date(2022, time.August, 22, 21, 30, 0, 0, time.UTC),
				Notification: time.Date(2022, time.August, 22, 14, 0, 0, 0, time.UTC),
				UserID:       user.ID,
			},
		}

		for _, event := range events {
			require.Nil(t, s.CreateEvent(ctx, event))
		}

		selectedEvents, err := s.SelectEventsByTime(ctx, events[0].Notification)
		require.Nil(t, err)

		for _, selectedEvent := range selectedEvents {
			require.True(t, containsEvent(events, selectedEvent))
		}

		selectedEvents, err = s.SelectEvents(ctx)
		require.Nil(t, err)

		for _, selectedEvent := range selectedEvents {
			require.True(t, containsEvent(events, selectedEvent))
		}

		events = selectedEvents
		events[0].Description = "Прийти на кв к Жеке"
		require.Nil(t, s.UpdateEvent(ctx, events[0]))

		selectedEvents, err = s.SelectEvents(ctx)
		require.Nil(t, err)

		for _, selectedEvent := range selectedEvents {
			require.True(t, containsEvent(events, selectedEvent))
			require.Nil(t, s.DeleteEvent(ctx, selectedEvent.ID))
		}

		selectedEvents, err = s.SelectEvents(ctx)
		require.Nil(t, err)
		require.Len(t, selectedEvents, 0)

		require.Nil(t, s.DeleteUser(ctx, user.ID))
		selectedUsers, err = s.SelectUsers(ctx)
		require.Nil(t, err)
		require.Len(t, selectedUsers, 0)
	})
}

func containsUser(users []storage.User, u storage.User) bool {
	for _, user := range users {
		if user.FirstName == u.FirstName &&
			user.LastName == u.LastName &&
			user.Email == u.Email &&
			user.Age == u.Age {
			return true
		}
	}
	return false
}

func containsEvent(events []storage.Event, e storage.Event) bool {
	for _, event := range events {
		if event.Finish == e.Finish &&
			event.Notification == e.Notification &&
			event.Beginning == e.Beginning &&
			event.Description == e.Description &&
			event.Title == e.Title {
			return true
		}
	}
	return false
}

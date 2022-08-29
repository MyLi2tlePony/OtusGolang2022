package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/app"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/config"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/logger"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/server/http/dto"
	memorystorage "github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	logConfig := config.LoggerConfig{
		Level: "info",
	}

	log := logger.New(logConfig)
	memoryStorage := memorystorage.New()
	application := app.New(memoryStorage)

	host := "localhost"
	port := "9090"
	servConfig := config.ServerConfig{
		Host: host,
		Port: port,
	}

	serv := NewServer(log, application, &servConfig)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		_ = serv.Start()
	}()

	mutex := sync.Mutex{}
	address := "http://" + net.JoinHostPort(host, port)

	userCase(t, &mutex, address)
	eventCase(t, &mutex, address)

	err := serv.Stop()
	require.Nil(t, err)

	wg.Wait()
}

func containsUser(users []dto.User, u dto.User) bool {
	for _, user := range users {
		if user.GetFirstName() == u.GetFirstName() &&
			user.GetLastName() == u.GetLastName() &&
			user.GetEmail() == u.GetEmail() &&
			user.GetAge() == u.GetAge() {
			return true
		}
	}
	return false
}

func containsEvent(events []dto.Event, e dto.Event) bool {
	for _, event := range events {
		if event.GetFinish() == e.GetFinish() &&
			event.GetNotification() == e.GetNotification() &&
			event.GetBeginning() == e.GetBeginning() &&
			event.GetDescription() == e.GetDescription() &&
			event.GetTitle() == e.GetTitle() {
			return true
		}
	}
	return false
}

func userCase(t *testing.T, mutex *sync.Mutex, address string) {
	t.Helper()

	t.Run("user case", func(t *testing.T) {
		mutex.Lock()
		defer mutex.Unlock()

		users := []dto.User{
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

		httpClient := &http.Client{}
		ctx := context.Background()

		for _, user := range users {
			jsonUser, err := json.Marshal(user)
			require.Nil(t, err)

			request, err := http.NewRequest(http.MethodPost, address+"/create/user", bytes.NewBuffer(jsonUser))
			require.Nil(t, err)
			request = request.WithContext(ctx)

			resp, err := httpClient.Do(request)
			require.Nil(t, resp.Body.Close())
			require.Nil(t, err)
		}

		request, err := http.NewRequest(http.MethodGet, address+"/select/users", nil)
		require.Nil(t, err)
		request = request.WithContext(ctx)

		resp, err := httpClient.Do(request)
		require.Nil(t, err)

		selectedUsers := new([]dto.User)
		err = json.NewDecoder(resp.Body).Decode(&selectedUsers)
		require.Nil(t, resp.Body.Close())
		require.Nil(t, err)

		require.Equal(t, len(users), len(*selectedUsers))
		for _, selectedUser := range *selectedUsers {
			require.True(t, containsUser(users, selectedUser))

			jsonUser, err := json.Marshal(selectedUser)
			require.Nil(t, err)

			request, err := http.NewRequest(http.MethodPost, address+"/delete/user", bytes.NewBuffer(jsonUser))
			require.Nil(t, err)
			request = request.WithContext(ctx)

			resp, err := httpClient.Do(request)
			require.Nil(t, resp.Body.Close())
			require.Nil(t, err)
		}

		request, err = http.NewRequest(http.MethodGet, address+"/select/users", nil)
		require.Nil(t, err)
		request = request.WithContext(ctx)

		resp, err = httpClient.Do(request)
		require.Nil(t, err)

		selectedUsers = new([]dto.User)
		err = json.NewDecoder(resp.Body).Decode(&selectedUsers)
		require.Nil(t, resp.Body.Close())
		require.Nil(t, err)

		require.Equal(t, 0, len(*selectedUsers))
	})
}

func eventCase(t *testing.T, mutex *sync.Mutex, address string) { //nolint:funlen
	t.Helper()

	t.Run("event case", func(t *testing.T) {
		mutex.Lock()
		defer mutex.Unlock()

		users := []dto.User{
			{
				FirstName: "Andrey",
				LastName:  "Shabarov",
				Email:     "shabandrew@mail.ru",
				Age:       21,
			},
		}

		httpClient := &http.Client{}
		ctx := context.Background()

		for _, user := range users {
			jsonUser, err := json.Marshal(user)
			require.Nil(t, err)

			request, err := http.NewRequest(http.MethodGet, address+"/create/user", bytes.NewBuffer(jsonUser))
			require.Nil(t, err)
			request = request.WithContext(ctx)

			resp, err := httpClient.Do(request)
			require.Nil(t, resp.Body.Close())
			require.Nil(t, err)
		}

		request, err := http.NewRequest(http.MethodGet, address+"/select/users", nil)
		require.Nil(t, err)
		request = request.WithContext(ctx)

		resp, err := httpClient.Do(request)
		require.Nil(t, err)

		selectedUsers := new([]dto.User)
		err = json.NewDecoder(resp.Body).Decode(&selectedUsers)
		require.Nil(t, resp.Body.Close())
		require.Nil(t, err)

		require.Equal(t, len(users), len(*selectedUsers))
		for _, selectedUser := range *selectedUsers {
			require.True(t, containsUser(users, selectedUser))
		}

		user := (*selectedUsers)[0]

		events := []dto.Event{
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
				Notification: time.Date(2022, time.August, 22, 13, 0, 0, 0, time.UTC),
				UserID:       user.ID,
			},
		}

		for _, event := range events {
			jsonEvent, err := json.Marshal(event)
			require.Nil(t, err)

			request, err := http.NewRequest(http.MethodPost, address+"/create/event", bytes.NewBuffer(jsonEvent))
			require.Nil(t, err)
			request = request.WithContext(ctx)

			resp, err := httpClient.Do(request)
			require.Nil(t, resp.Body.Close())
			require.Nil(t, err)
		}

		request, err = http.NewRequest(http.MethodGet, address+"/select/events", nil)
		require.Nil(t, err)
		request = request.WithContext(ctx)

		resp, err = httpClient.Do(request)
		require.Nil(t, err)

		selectedEvents := new([]dto.Event)
		err = json.NewDecoder(resp.Body).Decode(&selectedEvents)
		require.Nil(t, resp.Body.Close())
		require.Nil(t, err)

		require.Equal(t, len(events), len(*selectedEvents))
		for _, selectedEvent := range *selectedEvents {
			require.True(t, containsEvent(events, selectedEvent))
		}

		(*selectedEvents)[0].Title = "new title"
		jsonEvent, err := json.Marshal((*selectedEvents)[0])
		require.Nil(t, err)

		request, err = http.NewRequest(http.MethodPost, address+"/update/event", bytes.NewBuffer(jsonEvent))
		require.Nil(t, err)
		request = request.WithContext(ctx)

		resp, err = httpClient.Do(request)
		require.Nil(t, resp.Body.Close())
		require.Nil(t, err)

		request, err = http.NewRequest(http.MethodGet, address+"/select/events", nil)
		require.Nil(t, err)
		request = request.WithContext(ctx)

		resp, err = httpClient.Do(request)
		require.Nil(t, err)

		updatedEvents := new([]dto.Event)
		err = json.NewDecoder(resp.Body).Decode(&updatedEvents)
		require.Nil(t, resp.Body.Close())
		require.Nil(t, err)

		require.Equal(t, len(events), len(*updatedEvents))
		for _, updatedEvent := range *updatedEvents {
			require.True(t, containsEvent(*selectedEvents, updatedEvent))
			jsonEvent, err := json.Marshal(updatedEvent)
			require.Nil(t, err)

			request, err = http.NewRequest(http.MethodPost, address+"/delete/event", bytes.NewBuffer(jsonEvent))
			require.Nil(t, err)
			request = request.WithContext(ctx)

			resp, err = httpClient.Do(request)
			require.Nil(t, resp.Body.Close())
			require.Nil(t, err)
		}

		request, err = http.NewRequest(http.MethodGet, address+"/select/events", nil)
		require.Nil(t, err)
		request = request.WithContext(ctx)

		resp, err = httpClient.Do(request)
		require.Nil(t, err)

		selectedEvents = new([]dto.Event)
		err = json.NewDecoder(resp.Body).Decode(&selectedEvents)
		require.Nil(t, resp.Body.Close())
		require.Nil(t, err)

		require.Equal(t, 0, len(*selectedEvents))

		request, err = http.NewRequest(http.MethodGet, address+"/select/users", nil)
		require.Nil(t, err)
		request = request.WithContext(ctx)

		resp, err = httpClient.Do(request)
		require.Nil(t, err)

		selectedUsers = new([]dto.User)
		err = json.NewDecoder(resp.Body).Decode(&selectedUsers)
		require.Nil(t, resp.Body.Close())
		require.Nil(t, err)

		require.Equal(t, len(users), len(*selectedUsers))
		for _, selectedUser := range *selectedUsers {
			require.True(t, containsUser(users, selectedUser))

			jsonUser, err := json.Marshal(selectedUser)
			require.Nil(t, err)

			request, err = http.NewRequest(http.MethodPost, address+"/delete/user", bytes.NewBuffer(jsonUser))
			require.Nil(t, err)
			request = request.WithContext(ctx)

			resp, err = httpClient.Do(request)
			require.Nil(t, resp.Body.Close())
			require.Nil(t, err)
		}

		request, err = http.NewRequest(http.MethodGet, address+"/select/users", nil)
		require.Nil(t, err)
		request = request.WithContext(ctx)

		resp, err = httpClient.Do(request)
		require.Nil(t, err)

		selectedUsers = new([]dto.User)
		err = json.NewDecoder(resp.Body).Decode(&selectedUsers)
		require.Nil(t, resp.Body.Close())
		require.Nil(t, err)

		require.Equal(t, 0, len(*selectedUsers))
	})
}

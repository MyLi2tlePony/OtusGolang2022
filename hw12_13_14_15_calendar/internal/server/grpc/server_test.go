package grpc

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/app"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/config"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/logger"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/server/grpc/api"
	memorystorage "github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	go func() {
		err := serv.Start()
		require.Nil(t, err)
	}()

	defer func() {
		err := serv.Stop()
		require.Nil(t, err)
	}()

	address := net.JoinHostPort(host, port)
	mutex := sync.Mutex{}

	userCase(t, &mutex, address)
	eventCase(t, &mutex, address)
}

func eventCase(t *testing.T, mutex *sync.Mutex, address string) {
	t.Helper()

	t.Run("user case", func(t *testing.T) {
		mutex.Lock()
		defer mutex.Unlock()

		ctx := context.Background()

		users := []*api.User{
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

		options := grpc.WithTransportCredentials(insecure.NewCredentials())
		clientConnection, err := grpc.Dial(address, options)
		require.Nil(t, err)

		userClient := api.NewUserServiceClient(clientConnection)

		for _, user := range users {
			_, err = userClient.CreateUser(ctx, user)
			require.Nil(t, err)
		}

		streamUsers, err := userClient.SelectUsers(ctx, &api.Void{})
		require.Nil(t, err)

		selectedUsers := make([]*api.User, 0)
		for {
			user, err := streamUsers.Recv()
			if user == nil {
				break
			}
			require.Nil(t, err)
			selectedUsers = append(selectedUsers, user)
		}

		require.Equal(t, len(users), len(selectedUsers))
		for _, selectedUser := range selectedUsers {
			require.True(t, containsUser(users, selectedUser))
			_, err = userClient.DeleteUser(ctx, selectedUser)
			require.Nil(t, err)
		}

		streamUsers, err = userClient.SelectUsers(ctx, &api.Void{})
		require.Nil(t, err)

		selectedUsers = make([]*api.User, 0)
		for {
			user, err := streamUsers.Recv()
			if user == nil {
				break
			}
			require.Nil(t, err)
			selectedUsers = append(selectedUsers, user)
		}

		require.Equal(t, 0, len(selectedUsers))
	})
}

func userCase(t *testing.T, mutex *sync.Mutex, address string) { //nolint:gocognit
	t.Helper()

	t.Run("event case", func(t *testing.T) {
		mutex.Lock()
		defer mutex.Unlock()

		ctx := context.Background()

		users := []*api.User{
			{
				FirstName: "Andrey",
				LastName:  "Shabarov",
				Email:     "shabandrew@mail.ru",
				Age:       21,
			},
		}

		options := grpc.WithTransportCredentials(insecure.NewCredentials())
		clientConnection, err := grpc.Dial(address, options)
		require.Nil(t, err)

		userClient := api.NewUserServiceClient(clientConnection)

		for _, user := range users {
			_, err = userClient.CreateUser(ctx, user)
			require.Nil(t, err)
		}

		streamUsers, err := userClient.SelectUsers(ctx, &api.Void{})
		require.Nil(t, err)

		selectedUsers := make([]*api.User, 0)
		for {
			user, err := streamUsers.Recv()
			if user == nil {
				break
			}
			require.Nil(t, err)
			selectedUsers = append(selectedUsers, user)
		}

		require.Equal(t, len(users), len(selectedUsers))
		for _, selectedUser := range selectedUsers {
			require.True(t, containsUser(users, selectedUser))
			require.Nil(t, err)
		}

		user := selectedUsers[0]

		events := []*api.Event{
			{
				Title:         "Встреча с другом",
				Description:   "",
				BeginningT:    timestamppb.New(time.Date(2022, time.August, 22, 18, 0, 0, 0, time.UTC)),
				FinishT:       timestamppb.New(time.Date(2022, time.August, 22, 20, 0, 0, 0, time.UTC)),
				NotificationT: timestamppb.New(time.Date(2022, time.August, 22, 13, 0, 0, 0, time.UTC)),
				UserID:        user.ID,
			},
			{
				Title:         "Конференция",
				Description:   "Сделать презентацию",
				BeginningT:    timestamppb.New(time.Date(2022, time.August, 22, 20, 0, 0, 0, time.UTC)),
				FinishT:       timestamppb.New(time.Date(2022, time.August, 22, 21, 30, 0, 0, time.UTC)),
				NotificationT: timestamppb.New(time.Date(2022, time.August, 22, 13, 0, 0, 0, time.UTC)),
				UserID:        user.ID,
			},
		}

		eventClient := api.NewEventServiceClient(clientConnection)

		for _, event := range events {
			_, err = eventClient.CreateEvent(ctx, event)
			require.Nil(t, err)
		}

		streamEvent, err := eventClient.SelectEvents(ctx, &api.Void{})
		require.Nil(t, err)

		selectedEvents := make([]*api.Event, 0)
		for {
			event, err := streamEvent.Recv()
			if event == nil {
				break
			}
			require.Nil(t, err)
			selectedEvents = append(selectedEvents, event)
		}

		require.Equal(t, len(events), len(selectedEvents))
		for _, selectedEvent := range selectedEvents {
			require.True(t, containsEvent(events, selectedEvent))
		}

		selectedEvents[0].Title = "new title"
		_, err = eventClient.UpdateEvent(ctx, selectedEvents[0])
		require.Nil(t, err)

		streamEvent, err = eventClient.SelectEvents(ctx, &api.Void{})
		require.Nil(t, err)

		updatedEvents := make([]*api.Event, 0)
		for {
			event, err := streamEvent.Recv()
			if event == nil {
				break
			}
			require.Nil(t, err)
			updatedEvents = append(updatedEvents, event)
		}

		require.Equal(t, len(selectedEvents), len(updatedEvents))
		for _, updatedEvent := range updatedEvents {
			require.True(t, containsEvent(selectedEvents, updatedEvent))
			_, err = eventClient.DeleteEvent(ctx, updatedEvent)
			require.Nil(t, err)
		}

		streamEvent, err = eventClient.SelectEvents(ctx, &api.Void{})
		require.Nil(t, err)

		selectedEvents = make([]*api.Event, 0)
		for {
			event, err := streamEvent.Recv()
			if event == nil {
				break
			}
			require.Nil(t, err)
			selectedEvents = append(selectedEvents, event)
		}

		require.Equal(t, 0, len(selectedEvents))

		for _, selectedUser := range selectedUsers {
			_, err = userClient.DeleteUser(ctx, selectedUser)
			require.Nil(t, err)
		}

		selectedUsers = make([]*api.User, 0)
		for {
			user, err := streamUsers.Recv()
			if user == nil {
				break
			}
			require.Nil(t, err)
			selectedUsers = append(selectedUsers, user)
		}

		require.Equal(t, 0, len(selectedUsers))
	})
}

func containsUser(users []*api.User, u *api.User) bool {
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

func containsEvent(events []*api.Event, e *api.Event) bool {
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

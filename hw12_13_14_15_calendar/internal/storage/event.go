package storage

import (
	"time"
)

type Event struct {
	ID string

	Title       string
	Description string

	Beginning    time.Time
	Finish       time.Time
	Notification time.Time

	UserID string
}

func (event *Event) GetID() string {
	return event.ID
}

func (event *Event) GetTitle() string {
	return event.Title
}

func (event *Event) GetDescription() string {
	return event.Description
}

func (event *Event) GetBeginning() time.Time {
	return event.Beginning
}

func (event *Event) GetFinish() time.Time {
	return event.Finish
}

func (event *Event) GetNotification() time.Time {
	return event.Notification
}

func (event *Event) GetUserID() string {
	return event.UserID
}

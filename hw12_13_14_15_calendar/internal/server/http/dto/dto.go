package dto

import (
	"time"
)

type Event struct {
	ID string `json:"id"`

	Title       string `json:"title"`
	Description string `json:"description"`

	Beginning    time.Time `json:"beginning"`
	Finish       time.Time `json:"finish"`
	Notification time.Time `json:"notification"`

	UserID string `json:"userId"`
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

type User struct {
	ID string `json:"id"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	Email string `json:"email"`
	Age   int64  `json:"age"`
}

func (user *User) GetID() string {
	return user.ID
}

func (user *User) GetFirstName() string {
	return user.FirstName
}

func (user *User) GetLastName() string {
	return user.LastName
}

func (user *User) GetEmail() string {
	return user.Email
}

func (user *User) GetAge() int64 {
	return user.Age
}

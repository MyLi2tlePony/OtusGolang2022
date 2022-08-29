package storage

type User struct {
	ID string

	FirstName string
	LastName  string

	Email string
	Age   int64
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

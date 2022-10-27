package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/storage"
	pgx "github.com/jackc/pgx/v4"
)

type Storage struct {
	connString string
}

func New(connString string) *Storage {
	return &Storage{
		connString: connString,
	}
}

func (s *Storage) SelectUsers(ctx context.Context) (users []storage.User, err error) {
	users = make([]storage.User, 0)
	sql := `SELECT id, firstname, lastname, email, age FROM users;`

	conn, err := pgx.Connect(ctx, s.connString)
	if err != nil {
		return users, err
	}

	defer func() {
		if closeErr := conn.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}()

	rows, err := conn.Query(ctx, sql)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var user storage.User

		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Age)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *Storage) CreateUser(ctx context.Context, user storage.User) error {
	sql := `INSERT INTO users (firstname, lastname, email, age) VALUES ($1, $2, $3, $4);`

	conn, err := pgx.Connect(ctx, s.connString)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := conn.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}()

	_, err = conn.Exec(ctx, sql, user.FirstName, user.LastName, user.Email, user.Age)
	return err
}

func (s *Storage) DeleteUser(ctx context.Context, userID string) error {
	sql := `DELETE FROM users WHERE users.id = $1;`

	conn, err := pgx.Connect(ctx, s.connString)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := conn.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}()

	_, err = conn.Exec(ctx, sql, userID)
	return err
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (err error) {
	sql := `INSERT INTO 
				events (title, description, beginning, finish, notification, userid) 
			VALUES 
				($1, $2, $3, $4, $5, $6);`

	conn, err := pgx.Connect(ctx, s.connString)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := conn.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}()

	_, err = conn.Exec(ctx, sql,
		event.Title, event.Description, event.Beginning, event.Finish, event.Notification, event.UserID)
	return err
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	sql := `DELETE FROM events WHERE events.id = $1;`

	conn, err := pgx.Connect(ctx, s.connString)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := conn.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}()

	_, err = conn.Exec(ctx, sql, eventID)
	return err
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) (err error) {
	sql := `UPDATE
				events
			SET
				title = $2, description = $3, beginning = $4, finish = $5, notification = $6, userid = $7
			WHERE
				id = $1;`

	conn, err := pgx.Connect(ctx, s.connString)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := conn.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}()

	_, err = conn.Exec(ctx, sql,
		event.ID, event.Title, event.Description, event.Beginning, event.Finish, event.Notification, event.UserID)
	return err
}

func (s *Storage) SelectEventsByTime(ctx context.Context, t time.Time) (events []storage.Event, err error) {
	events = make([]storage.Event, 0)
	sql := "SELECT id, title, description, beginning, finish, notification, userid FROM events WHERE notification = " +
		fmt.Sprintf("'%02d-%02d-%02d %02d:%02d';", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())

	conn, err := pgx.Connect(ctx, s.connString)
	if err != nil {
		return events, err
	}

	defer func() {
		if closeErr := conn.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}()

	rows, err := conn.Query(ctx, sql)
	if err != nil {
		return events, err
	}

	for rows.Next() {
		var event storage.Event

		err = rows.Scan(&event.ID, &event.Title, &event.Description,
			&event.Beginning, &event.Finish, &event.Notification, &event.UserID)
		if err != nil {
			return events, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) SelectEvents(ctx context.Context) (events []storage.Event, err error) {
	events = make([]storage.Event, 0)
	sql := `SELECT id, title, description, beginning, finish, notification, userid FROM events;`

	conn, err := pgx.Connect(ctx, s.connString)
	if err != nil {
		return events, err
	}

	defer func() {
		if closeErr := conn.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}()

	rows, err := conn.Query(ctx, sql)
	if err != nil {
		return events, err
	}

	for rows.Next() {
		var event storage.Event

		err = rows.Scan(&event.ID, &event.Title, &event.Description,
			&event.Beginning, &event.Finish, &event.Notification, &event.UserID)
		if err != nil {
			return events, err
		}

		events = append(events, event)
	}

	return events, nil
}

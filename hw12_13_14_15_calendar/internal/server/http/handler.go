package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/server"
	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/server/http/dto"
)

type handler struct {
	logger server.Logger
	app    server.Application
}

func newHandler(logger server.Logger, app server.Application) *handler {
	return &handler{
		logger: logger,
		app:    app,
	}
}

func (h *handler) createUser(w http.ResponseWriter, r *http.Request) {
	user, err := readUserFromBody(r)
	if err != nil {
		h.logger.Error(err.Error())
		w.Header().Add("status", fmt.Sprint(http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.app.CreateUser(context.Background(), user); err != nil {
		h.logger.Error(err.Error())
		w.Header().Add("status", fmt.Sprint(http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("status", fmt.Sprint(http.StatusOK))
	w.WriteHeader(http.StatusOK)
}

func (h *handler) selectUsers(w http.ResponseWriter, r *http.Request) {
	marshal, err := selectAsJSON(func(ctx context.Context) (interface{}, error) {
		events, err := h.app.SelectUsers(ctx)
		e := interface{}(events)
		return e, err
	})
	if err != nil {
		h.logger.Error(err.Error())
		w.Header().Add("status", fmt.Sprint(http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = sendData(marshal, w)
	if err != nil {
		h.logger.Error(err.Error())
		return
	}
}

func (h *handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	user, err := readUserFromBody(r)
	if err != nil {
		h.logger.Error(err.Error())
		w.Header().Add("status", fmt.Sprint(http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.app.DeleteUser(context.Background(), user.ID); err != nil {
		h.logger.Error(err.Error())
		w.Header().Add("status", fmt.Sprint(http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("status", fmt.Sprint(http.StatusOK))
	w.WriteHeader(http.StatusOK)
}

func (h *handler) createEvent(w http.ResponseWriter, r *http.Request) {
	event, err := readEventFromBody(r)
	if err != nil {
		h.logger.Error(err.Error())
		w.Header().Add("status", fmt.Sprint(http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.app.CreateEvent(context.Background(), event); err != nil {
		h.logger.Error(err.Error())
		w.Header().Add("status", fmt.Sprint(http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("status", fmt.Sprint(http.StatusOK))
	w.WriteHeader(http.StatusOK)
}

func (h *handler) selectEvents(w http.ResponseWriter, r *http.Request) {
	marshal, err := selectAsJSON(func(ctx context.Context) (interface{}, error) {
		events, err := h.app.SelectEvents(ctx)
		e := interface{}(events)
		return e, err
	})
	if err != nil {
		h.logger.Error(err.Error())
		w.Header().Add("status", fmt.Sprint(http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = sendData(marshal, w)
	if err != nil {
		h.logger.Error(err.Error())
		return
	}
}

func (h *handler) updateEvent(w http.ResponseWriter, r *http.Request) {
	event, err := readEventFromBody(r)
	if err != nil {
		h.logger.Error(err.Error())
		w.Header().Add("status", fmt.Sprint(http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.app.UpdateEvent(context.Background(), event); err != nil {
		h.logger.Error(err.Error())
		w.Header().Add("status", fmt.Sprint(http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("status", fmt.Sprint(http.StatusOK))
	w.WriteHeader(http.StatusOK)
}

func (h *handler) deleteEvent(w http.ResponseWriter, r *http.Request) {
	event, err := readEventFromBody(r)
	if err != nil {
		h.logger.Error(err.Error())
		w.Header().Add("status", fmt.Sprint(http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.app.DeleteEvent(context.Background(), event.ID); err != nil {
		h.logger.Error(err.Error())
		w.Header().Add("status", fmt.Sprint(http.StatusBadRequest))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("status", fmt.Sprint(http.StatusOK))
	w.WriteHeader(http.StatusOK)
}

func readUserFromBody(r *http.Request) (*dto.User, error) {
	user := new(dto.User)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return user, err
	}

	if err = json.Unmarshal(body, user); err != nil {
		return user, err
	}

	return user, nil
}

func readEventFromBody(r *http.Request) (*dto.Event, error) {
	event := new(dto.Event)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return event, err
	}

	if err = json.Unmarshal(body, event); err != nil {
		return event, err
	}

	return event, nil
}

func selectAsJSON(sel func(context.Context) (interface{}, error)) ([]byte, error) {
	var marshal []byte

	events, err := sel(context.Background())
	if err != nil {
		return marshal, err
	}

	marshal, err = json.Marshal(events)
	if err != nil {
		return marshal, err
	}

	return marshal, nil
}

func sendData(marshal []byte, w http.ResponseWriter) error {
	w.Header().Add("status", fmt.Sprint(http.StatusOK))
	w.WriteHeader(http.StatusOK)

	_, err := w.Write(marshal)
	return err
}

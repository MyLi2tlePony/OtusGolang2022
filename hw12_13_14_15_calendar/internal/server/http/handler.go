package internalhttp

import (
	"fmt"
	"io"
	"net/http"
)

type handler struct {
	logger Logger
}

func newHandler(logger Logger) *handler {
	return &handler{
		logger: logger,
	}
}

func (h *handler) getHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("status", fmt.Sprint(http.StatusOK))
	w.WriteHeader(http.StatusOK)

	_, err := io.WriteString(w, "hello-world\n")
	if err != nil {
		h.logger.Error(err.Error())
		return
	}
}

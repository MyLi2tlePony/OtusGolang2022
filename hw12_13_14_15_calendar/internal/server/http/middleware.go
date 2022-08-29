package internalhttp

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/server"
)

type middleware struct {
	logger  server.Logger
	Handler http.Handler
}

func newMiddleware(logger server.Logger, httpHandler http.Handler) *middleware {
	return &middleware{
		logger:  logger,
		Handler: httpHandler,
	}
}

func (m *middleware) logging() *middleware {
	curHandler := m.Handler

	m.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		curHandler.ServeHTTP(w, r)
		handleTime := time.Since(start)

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			m.logger.Error(fmt.Sprintf("error split host port, remote address:%s", r.RemoteAddr))
		}

		statusCode := w.Header().Get("status")

		logMessage := fmt.Sprintf("%s [%s] %s %s %s %s %s %s",
			ip, start, r.Method, r.URL, r.Proto, statusCode, handleTime, r.UserAgent())
		m.logger.Info(logMessage)
	})

	return m
}

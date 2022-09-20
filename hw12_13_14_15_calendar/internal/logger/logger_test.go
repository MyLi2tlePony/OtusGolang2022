package logger

import (
	"testing"

	"github.com/MyLi2tlePony/OtusGolang2022/hw12_13_14_15_calendar/internal/config/calendar"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("base case", func(t *testing.T) {
		New(&calendar.LoggerConfig{
			Level: "info",
		})
		require.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())

		New(&calendar.LoggerConfig{
			Level: "error",
		})
		require.Equal(t, zerolog.ErrorLevel, zerolog.GlobalLevel())

		New(&calendar.LoggerConfig{
			Level: "warn",
		})
		require.Equal(t, zerolog.WarnLevel, zerolog.GlobalLevel())
	})
}

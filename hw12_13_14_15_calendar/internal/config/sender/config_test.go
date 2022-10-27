package sender

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("read config case", func(t *testing.T) {
		fileName := "testConfig.*.toml"
		file, err := os.CreateTemp("", fileName)
		require.Nil(t, err)

		defer func() {
			require.Nil(t, file.Close())
			require.Nil(t, os.Remove(file.Name()))
		}()

		expectedConfig := Config{
			Logger: &LoggerConfig{
				Level: "info",
			},
			Consume: &ConsumeConfig{
				queue:     "notification",
				consumer:  "",
				autoAck:   true,
				exclusive: false,
				noLocal:   false,
				noWait:    false,
			},
			Connection: &BrokerConnectionConfig{
				login:    "guest",
				password: "guest",
				host:     "localhost",
				port:     "5672",
			},
			Queue: &BrokerQueueConfig{
				name:       "notification",
				durable:    false,
				autoDelete: false,
				exclusive:  false,
			},
		}

		configText := []byte(`
			[logger]
			level = "info"
			
			[connection]
			login = "guest"
			password = "guest"
			host = "localhost"
			port = "5672"
			
			[consume]
			queue = "notification"
			consumer = ""
			autoAck = true
			exclusive = false
			noLocal = false
			noWait = false
			
			[queue]
			name = "notification"
			durable = false
			autoDelete = false
			exclusive = false
			noWait = false
		`)

		_, err = file.Write(configText)
		require.Nil(t, err)

		config, err := NewConfig(file.Name())
		require.Nil(t, err)

		require.True(t, reflect.DeepEqual(config, expectedConfig))
	})
}

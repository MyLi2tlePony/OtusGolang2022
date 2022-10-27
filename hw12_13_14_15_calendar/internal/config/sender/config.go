package sender

import (
	"fmt"

	"github.com/spf13/viper"
)

type BrokerQueueConfig struct {
	name       string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
}

func NewBrokerQueueConfig(name string, durable, autoDelete, exclusive, noWait bool) *BrokerQueueConfig {
	return &BrokerQueueConfig{
		name:       name,
		durable:    durable,
		autoDelete: autoDelete,
		exclusive:  exclusive,
		noWait:     noWait,
	}
}

func (config *BrokerQueueConfig) Name() string {
	return config.name
}

func (config *BrokerQueueConfig) Durable() bool {
	return config.durable
}

func (config *BrokerQueueConfig) AutoDelete() bool {
	return config.autoDelete
}

func (config *BrokerQueueConfig) Exclusive() bool {
	return config.exclusive
}

func (config *BrokerQueueConfig) NoWait() bool {
	return config.noWait
}

type BrokerConnectionConfig struct {
	login    string
	password string
	host     string
	port     string
}

func NewBrokerConnectionConfig(login, password, host, port string) *BrokerConnectionConfig {
	return &BrokerConnectionConfig{
		login:    login,
		password: password,
		host:     host,
		port:     port,
	}
}

func (config *BrokerConnectionConfig) Login() string {
	return config.login
}

func (config *BrokerConnectionConfig) Password() string {
	return config.password
}

func (config *BrokerConnectionConfig) Host() string {
	return config.host
}

func (config *BrokerConnectionConfig) Port() string {
	return config.port
}

type ConsumeConfig struct {
	queue     string
	consumer  string
	autoAck   bool
	exclusive bool
	noLocal   bool
	noWait    bool
}

func NewConsumeConfig(queue, consumer string, autoAck, exclusive, noLocal, noWait bool) *ConsumeConfig {
	return &ConsumeConfig{
		queue:     queue,
		consumer:  consumer,
		autoAck:   autoAck,
		exclusive: exclusive,
		noLocal:   noLocal,
		noWait:    noWait,
	}
}

func (config *ConsumeConfig) Queue() string {
	return config.queue
}

func (config *ConsumeConfig) Consumer() string {
	return config.consumer
}

func (config *ConsumeConfig) AutoAck() bool {
	return config.autoAck
}

func (config *ConsumeConfig) Exclusive() bool {
	return config.exclusive
}

func (config *ConsumeConfig) NoLocal() bool {
	return config.noLocal
}

func (config *ConsumeConfig) NoWait() bool {
	return config.noWait
}

type LoggerConfig struct {
	Level string
}

func (config *LoggerConfig) GetLevel() string {
	return config.Level
}

type Config struct {
	Connection *BrokerConnectionConfig
	Queue      *BrokerQueueConfig
	Consume    *ConsumeConfig
	Logger     *LoggerConfig
}

func NewConfig(configPath string) (Config, error) {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("fatal error config file: %w", err)
	}

	return Config{
		Logger: &LoggerConfig{
			Level: viper.GetString("logger.level"),
		},
		Consume: &ConsumeConfig{
			queue:     viper.GetString("consume.queue"),
			consumer:  viper.GetString("consume.consumer"),
			autoAck:   viper.GetBool("consume.autoAck"),
			exclusive: viper.GetBool("consume.exclusive"),
			noLocal:   viper.GetBool("consume.noLocal"),
			noWait:    viper.GetBool("consume.noWait"),
		},
		Connection: &BrokerConnectionConfig{
			login:    viper.GetString("connection.login"),
			password: viper.GetString("connection.password"),
			host:     viper.GetString("connection.host"),
			port:     viper.GetString("connection.port"),
		},
		Queue: &BrokerQueueConfig{
			name:       viper.GetString("queue.name"),
			durable:    viper.GetBool("queue.durable"),
			autoDelete: viper.GetBool("queue.autoDelete"),
			exclusive:  viper.GetBool("queue.exclusive"),
		},
	}, nil
}

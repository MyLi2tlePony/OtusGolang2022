package scheduler

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

type PublishConfig struct {
	exchange    string
	key         string
	contentType string
	mandatory   bool
	immediate   bool
}

func NewPublishConfig(exchange, key, contentType string, mandatory, immediate bool) *PublishConfig {
	return &PublishConfig{
		exchange:    exchange,
		key:         key,
		contentType: contentType,
		mandatory:   mandatory,
		immediate:   immediate,
	}
}

func (config *PublishConfig) Exchange() string {
	return config.exchange
}

func (config *PublishConfig) Key() string {
	return config.key
}

func (config *PublishConfig) ContentType() string {
	return config.contentType
}

func (config *PublishConfig) Mandatory() bool {
	return config.mandatory
}

func (config *PublishConfig) Immediate() bool {
	return config.immediate
}

type LoggerConfig struct {
	Level string
}

func (config *LoggerConfig) GetLevel() string {
	return config.Level
}

type DatabaseConfig struct {
	Prefix       string
	DatabaseName string
	Host         string
	Port         string
	UserName     string
	Password     string
}

type Config struct {
	Connection *BrokerConnectionConfig
	Queue      *BrokerQueueConfig
	Publish    *PublishConfig
	Logger     *LoggerConfig
	Database   *DatabaseConfig
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
		Publish: &PublishConfig{
			exchange:    viper.GetString("publish.exchange"),
			key:         viper.GetString("publish.key"),
			contentType: viper.GetString("publish.contentType"),
			mandatory:   viper.GetBool("publish.mandatory"),
			immediate:   viper.GetBool("publish.immediate"),
		},
		Database: &DatabaseConfig{
			Prefix:       viper.GetString("database.Prefix"),
			DatabaseName: viper.GetString("database.DatabaseName"),
			Host:         viper.GetString("database.Host"),
			Port:         viper.GetString("database.Port"),
			UserName:     viper.GetString("database.UserName"),
			Password:     viper.GetString("database.Password"),
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

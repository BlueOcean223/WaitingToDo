package configs

import "github.com/rabbitmq/amqp091-go"

const (
	FriendRequestQueue = "friend_request"
	TeamRequestQueue   = "team_request"
)

type RabbitMQQueueConfig struct {
	Name       string `yaml:"name"`
	RoutingKey string `yaml:"routing_key"`
}

type RabbitMQConfig struct {
	Dsn      string                         `yaml:"dsn"`
	Exchange string                         `yaml:"exchange"`
	Queues   map[string]RabbitMQQueueConfig `yaml:"queues"`
}

var RabbitMQConn *amqp091.Connection

func InitRabbitMQ() error {
	conn, err := amqp091.Dial(AppConfigs.RabbitMQConfig.Dsn)
	if err != nil {
		return err
	}
	RabbitMQConn = conn
	return nil
}

package consumer

import (
	"backend/internal/configs"
	"backend/pkg/logger"
	"context"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

/*
	消费者基础类
*/

type BaseConsumer struct {
	name       string         // 消费者名称
	queueName  string         // 队列名称
	routingKey string         // 路由键
	exchange   string         // 交换机名称
	handler    MessageHandler // 消息处理函数
}

// MessageHandler 定义消息处理函数类型
type MessageHandler func(body []byte) error

func NewBaseConsumer(name, queueName, routingKey, exchange string, handler MessageHandler) *BaseConsumer {
	return &BaseConsumer{
		name:       name,
		queueName:  queueName,
		routingKey: routingKey,
		exchange:   exchange,
		handler:    handler,
	}
}

func (bc *BaseConsumer) GetName() string {
	return bc.name
}

// Start 启动消费者
func (bc *BaseConsumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := bc.consume(ctx); err != nil {
				log.Printf("%s: 消费失败，5秒后重试: %v", bc.name, err)
				time.Sleep(5 * time.Second)
			}
		}
	}
}

// consume 执行消费逻辑
func (bc *BaseConsumer) consume(ctx context.Context) error {
	// 建立连接
	conn, err := amqp091.Dial(configs.AppConfigs.RabbitMQConfig.Dsn)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 创建channel
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// 声明交换机
	err = ch.ExchangeDeclare(bc.exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// 声明队列
	queue, err := ch.QueueDeclare(bc.queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// 绑定队列
	err = ch.QueueBind(queue.Name, bc.routingKey, bc.exchange, false, nil)
	if err != nil {
		return err
	}

	// 开始消费
	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	log.Printf("%s 消费者已启动", bc.name)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}

			if err := bc.handler(msg.Body); err != nil {
				log.Printf("%s: 处理消息失败: %v", bc.name, err)
				err := msg.Nack(false, true)
				if err != nil {
					logger.Error("消息拒绝失败: ", logger.Err(err))
				}
			} else {
				err := msg.Ack(false)
				if err != nil {
					logger.Error("消息确认失败: ", logger.Err(err))
				}
			}
		}
	}
}

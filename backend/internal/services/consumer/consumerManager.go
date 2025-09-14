package consumer

import (
	"context"
	"log"
	"sync"
)

/*
	消费者管理类
*/

type ConsumerManager struct {
	consumers map[string]Consumer // 消费者集合
	ctx       context.Context     // 用于控制消费者的生命周期
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

// Consumer 消费者的基本行为
type Consumer interface {
	Start(ctx context.Context) error // 启动消费者
	GetName() string
}

func NewConsumerManager() *ConsumerManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ConsumerManager{
		consumers: make(map[string]Consumer),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// RegisterConsumer 注册消费者
func (cm *ConsumerManager) RegisterConsumer(consumer Consumer) {
	cm.consumers[consumer.GetName()] = consumer
}

// UnregisterConsumer 注销消费者
func (cm *ConsumerManager) UnregisterConsumer(name string) {
	if _, exists := cm.consumers[name]; exists {
		delete(cm.consumers, name)
		log.Printf("消费者 %s 已注销", name)
	} else {
		log.Printf("消费者 %s 不存在", name)
	}
}

// StartAll 启动所有消费者
func (cm *ConsumerManager) StartAll() {
	for name, consumer := range cm.consumers {
		cm.wg.Add(1)
		go func(name string, c Consumer) {
			defer cm.wg.Done()
			if err := c.Start(cm.ctx); err != nil {
				log.Printf("消费者 %s 启动失败: %v", name, err)
			}
		}(name, consumer)
	}
}

// Stop 停止所有消费者
func (cm *ConsumerManager) Stop() {
	cm.cancel()
	cm.wg.Wait()
	log.Println("所有消费者已停止")
}

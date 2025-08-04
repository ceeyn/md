package kafka

import (
	"sync"

	"git.code.oa.com/trpc-go/trpc-database/kafka"
)

const (
	TemplateDBProxyProducerName = "trpc.kafka.producer.db_proxy"
)

//ProducerProxy struct
type ProducerProxy struct {
	TemplateDBProxyProducer kafka.Client
}

//producerProxy 生产者单例
var producerProxyInstance ProducerProxy
var ProducerOnce sync.Once

// InitProducer 初始化生产者
func InitProducer() {
	ProducerOnce.Do(func() {
		templateDBProxyProducer := kafka.NewClientProxy(TemplateDBProxyProducerName)
		producerProxyInstance.TemplateDBProxyProducer = templateDBProxyProducer
	})
}

package rpc

import (
	"context"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	"github.com/apache/pulsar-client-go/pulsar"
	"meeting_template/config"
	"time"
)

// tdmq最长延时时间10天，预留10s的差值
const MaxTimeToDelayMq = (10 * 24 * 3600 - 120)

var TdmqProducer pulsar.Producer

// InitTdmqProducer ...
func InitTdmqProducer() {
	TdmqProducer = InitProducer(
		config.TdmqConfig.TdmqClusterUrl,
		config.TdmqConfig.TdmqRoleToken,
		config.TdmqConfig.TdmqTopic,
	)
}


// InitProducer 初始化生产者
func InitProducer(url, token, topic string) pulsar.Producer {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		// 服务接入地址
		URL: url,
		// 授权角色密钥
		Authentication:    pulsar.NewAuthenticationToken(token),
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		metrics.IncrCounter("GetTdmqClient_Failed", 1)
		log.Fatalf("Could not instantiate Pulsar client: %+v", err)
	}
	log.Infof("pulsar client created succ!",)
	// 使用客户端创建生产者
	tdmqProducer, err := client.CreateProducer(pulsar.ProducerOptions{
		// topic完整路径，格式为persistent://集群（租户）ID/命名空间/Topic名称
		Topic: topic,
	})
	if err != nil {
		metrics.IncrCounter("CreateTdmqProducerFailed", 1)
		log.Fatalf("TdmqProducerCreate failed, err:%+v, url= %v, topic= %v", err, url, topic)
	}
	log.Infof("TdmqProducer created succ!")
	return tdmqProducer
}

// SendDelayJob ... 这里用的是AfterDeliver
func SendDelayJob(ctx context.Context, msg []byte, delaySpan int64, meetingId uint64) error {
	// 发送消息
	// 最大时长只能是10天，超过10天按10天往里面放，消费方会处理；
	if delaySpan >= MaxTimeToDelayMq {
		delaySpan = MaxTimeToDelayMq
	}
	log.InfoContextf(ctx, "SendDelayJob meetingId:%+v, msg:%+v, client: %v, delaySpan: %v",
		meetingId, string(msg), TdmqProducer.Name(), delaySpan)
	tdmqMsgId, err := TdmqProducer.Send(ctx, &pulsar.ProducerMessage{
		// 消息内容
		Payload: msg,
		// 业务参数
		Properties: map[string]string{"key": "value"},
		// 延迟时间
		DeliverAfter: time.Second * time.Duration(delaySpan),
	})
	log.InfoContextf(ctx, "SendDelayJob meetingId: %v, tdmqMsgId:%+v", meetingId, tdmqMsgId)
	if err != nil {
		log.ErrorContextf(ctx, "SendDelayJobErr, err is:%+v, tdmqMsgId:%+v, msg:%+v",
			err, tdmqMsgId, string(msg))
		metrics.IncrCounter("SendDelayJobErr", 1)
		return err
	}
	metrics.IncrCounter("SendDelayJobSuccess", 1)
	return nil
}
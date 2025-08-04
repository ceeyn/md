// kafkaMsg消息处理
package service_kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.code.oa.com/meettrpc/meet_util"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	kfPb "git.code.oa.com/trpcprotocol/wemeet/wemeet_kafka_message"
	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/ants/v2"
)

// KafkaConsumer 回调对象
type KafkaConsumer struct{}

// MsgHandlePool 处理kafka消息的协程池,容量为100000个协程
var MsgHandlePool, _ = ants.NewPool(
	100000,
	ants.WithPreAlloc(true),
	ants.WithExpiryDuration(time.Hour))

var (
	gMsgHandleMap = map[uint64]func(ctx context.Context, msg *kfPb.KafkaMessage) error{
		// 修改会议
		uint64(kfPb.KAFKA_MSG_ID_KAFKA_MSG_WEBHOOK_CHANGE_MEETING_NOTIFY): HandleMeetingChangeNotify,
	}
)

// Handle 回调方法
func (KafkaConsumer) Handle(ctx context.Context, key, value []byte, topic string, partition int32, offset int64) error {
	kafkaMsg := &kfPb.KafkaMessage{}
	if err := proto.Unmarshal(value, kafkaMsg); err != nil {
		log.ErrorContextf(ctx, "kafka message Unmarshal fail ：%+v err:%+v", value, err.Error())
		return nil
	}

	//查找处理函数
	f, ok := gMsgHandleMap[kafkaMsg.GetMsgId()]
	if !ok {
		metrics.IncrCounter("HandleMap.MsgId.Err", 1)
		return nil
	}

	// 协程池 消费
	newCtx := trpc.CloneContext(ctx)
	if err := MsgHandlePool.Submit(func() {
		defer meet_util.DefPanicFun()
		log.WithContextFields(newCtx,
			"msgId", fmt.Sprintf("%v", kafkaMsg.GetMsgId()),
			"meetId", fmt.Sprintf("%v", kafkaMsg.GetMeetingId()),
			"topic", fmt.Sprintf("%v", topic),
			"partition", fmt.Sprintf("%v", partition),
			"offset", fmt.Sprintf("%v", offset),
		)
		if err := f(newCtx, kafkaMsg); err != nil {
			log.ErrorContextf(newCtx, "(KafkaConsumer) Handle kafka message consume fail, meetingId:%+v,"+
				" MsgId:%v, err:%v", kafkaMsg.GetMeetingId(), kafkaMsg.GetMsgId(), err.Error())
			return
		}
	}); err != nil {
		log.ErrorContextf(newCtx, "MsgHandlePool pool err: %v", err.Error())
		return errors.New("msg handle pool err")
	}
	return nil
}

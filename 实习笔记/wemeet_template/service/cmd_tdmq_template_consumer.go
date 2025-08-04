package service

import (
	"context"
	"encoding/json"
	"fmt"
	"git.code.oa.com/trpc-go/trpc-database/tdmq"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"google.golang.org/protobuf/proto"
	"meeting_template/config/config_rainbow"
	"meeting_template/model"
	"meeting_template/rpc"
	"meeting_template/util"
	"strconv"
	"strings"
)

// Consumer ...
type Consumer struct {
}

// Handle ...
func (Consumer) Handle(ctx context.Context, message *tdmq.Message) error {
	messageID, err := tdmq.DeserializeMessageID(message.MsgID)
	if err != nil {
		metrics.IncrCounter("tdmq.DeserializeMessageID.fail",1)
		log.ErrorContextf(ctx, "ConsumeTriggerMsg tdmq consumer msg DeserializeMessageID err: %v", err)
	}
	metrics.IncrCounter("ConsumeDelayTrigger", 1)
	log.InfoContextf(ctx,"wemeet_template ConsumeTriggerMsg Handle [key]: %+v, [payload]: %+v, [topic]: %+v, [id]: %+v",
		message.Key, string(message.PayLoad), message.Topic, messageID)

	triggerInfo := &model.TemplateMqMsg{}
	err = json.Unmarshal(message.PayLoad, triggerInfo)
	if err != nil {
		metrics.IncrCounter("TDMQUnmarshalMsgErr", 1)
		log.ErrorContextf(ctx, "TDMQ Consumer failed to unmarshal delay msg:%+v, err:%+v", string(message.PayLoad), err)
		return nil
	}
	metrics.IncrCounter("tdmq.consume.sum",1)
	log.InfoContextf(ctx, "wemeet_template consume from delay success, msg:%+v, tId:%+v", triggerInfo, messageID)
	HandleDelayQueryM3U8TransStatus(ctx, triggerInfo)
	// 返回nil才会确认消费成功
	return nil
}

// HandleDelayQueryM3U8TransStatus ...
func HandleDelayQueryM3U8TransStatus(ctx context.Context, templateMqMsg *model.TemplateMqMsg)  {
	templateId   := templateMqMsg.TemplateId
	msgMeetingId := templateMqMsg.MeetingId
    msgVideoMp4CosId := templateMqMsg.CosId
    // decoded 一下
	decodedMp4CosId := string(util.GetBase64Decoded(ctx, msgVideoMp4CosId))
    m3u8CosId := strings.Replace(decodedMp4CosId, ".mp4", ".m3u8", 1)

    isExist := rpc.JudgeCosResourceIsExist(ctx, m3u8CosId)
	if isExist {   //转码完成，更新缓存和推送
		ProcessM3U8TransDone(ctx, msgVideoMp4CosId, templateId, msgMeetingId)
	} else {
		maxRetryCnt := config_rainbow.GetMusicConfConfig().MaxRetryCnt
		if templateMqMsg.TryCount >= uint32(maxRetryCnt) {
			log.ErrorContextf(ctx,"HandleDelayQueryM3U8TransStatus TryCount over max count, templateMqMsg:%+v", templateMqMsg)
			metrics.IncrCounter("templateMqMsg.OverMaxCount",1)
			return
		}
		tempMsg := templateMqMsg
		tempMsg.TryCount += 1
		byteMsg, err := json.Marshal(tempMsg)
		if err != nil {
			metrics.IncrCounter("tdmq.MsgMashal.fail",1)
			log.ErrorContextf(ctx, "HandleDelayQueryM3U8TransStatus json Marshal TemplateMqMsg failed, tempMsg:%+v, err:%+v", tempMsg, err)
			return
		}
		// 继续延迟查询
		delaySpan := config_rainbow.GetMusicConfConfig().DelaySpan
		log.InfoContextf(ctx,"HandleDelayQueryM3U8TransStatus need retry SendDelayJob, " +
			"tdmqMsg:%+v, DelaySpan:%+v, maxRetryCnt:%+v", tempMsg, delaySpan, maxRetryCnt)
		rpc.SendDelayJob(ctx, byteMsg, delaySpan, msgMeetingId)
	}
}

// ProcessM3U8TransDone ...
func ProcessM3U8TransDone(ctx context.Context, msgVideoMp4CosId string, templateId string, msgMeetingId uint64)  {
	templateInfo := &model.TemplateInfo{}
	templateInfo, err := GetTemplateInfoSingleFlight(ctx, templateId)
	if err != nil {
		metrics.IncrCounter("ProcessM3U8TransDone.GetTemplate.Fail",1)
		log.ErrorContextf(ctx, "ProcessM3U8TransDone, get templateInfo from redis failed, templateId:%v", templateId)
		return
	}
	if templateInfo.WarmUpData == "" {
		metrics.IncrCounter("ProcessM3U8TransDone.WarmUpData.empty",1)
		log.ErrorContextf(ctx, "ProcessM3U8TransDone get warmup data empty, templateId:%+v", templateId)
		return
	}
	if templateInfo.MeetingId != fmt.Sprint(msgMeetingId) {
		metrics.IncrCounter("ProcessM3U8TransDone.MeetingId.NotSame",1)
		log.ErrorContextf(ctx,"tdmqMsg meetingId:%+v not the same with cache meetingId:%+v",
			templateInfo.MeetingId, msgMeetingId)
		return
	}
	buf := []byte(templateInfo.WarmUpData)
	warmUpData := pb.WarmUpData{}
	err = json.Unmarshal(buf, &warmUpData)
	if err != nil {
		metrics.IncrCounter("ProcessM3U8TransDone.WarmUpData.UnmarshalFail",1)
		log.ErrorContextf(ctx, "ProcessM3U8TransDone json parse warmup fail, warmUpData:%+v, error:%+v", templateInfo.WarmUpData, err)
		return
	}
	// 更新视频转码状态
	if !NeedModifyWarmUpVideoTransStatus(ctx, msgVideoMp4CosId, &warmUpData) {
		metrics.IncrCounter("DoneNeedModifyWarmUpVideoTransStatus.sum",1)
		log.InfoContextf(ctx,"ProcessM3U8TransDone No NeedModifyWarmUpVideoTransStatus, warmUpData:%+v", warmUpData)
		return
	}
	// 在这里赋值转码成功
	warmUpData.GetWarmupVideoList()[0].VideoTransStatus = proto.Uint32(util.M3U8TransDone)
	// 说明需要更新cache和推送客户端
	warmUpDataStr, err := util.GetSerializedJsonStr(ctx, warmUpData)
	if err != nil {
		metrics.IncrCounter("ProcessM3U8TransDone.WarmUpData.SerializedJsonFail",1)
		log.InfoContextf(ctx, "ProcessM3U8TransDone warmUpData SerializedJsonStr failed, warmUpData::%v, err:%v", warmUpData, err)
	} else {
		templateInfo.WarmUpData = warmUpDataStr
	}
	err = SetTemplateInfo(ctx, templateInfo)
	if err != nil {
		log.ErrorContextf(ctx, "ProcessM3U8TransDone set templateInfo fail, templateId:%v", templateId)
	}
	appId, err := strconv.ParseUint(templateInfo.AppId, 10, 32)
	if err != nil {
		metrics.IncrCounter("ProcessM3U8TransDone.AppId.ParseFail",1)
		log.ErrorContextf(ctx, "ProcessM3U8TransDone ParseUintAppIdFailed, tplInfo:%+v, err:%+v", templateInfo, err)
	}
	// 做一个推送
	UserNotifyWarmUpDataUpdate(ctx, msgMeetingId, uint32(appId), templateInfo.AppUid, templateInfo)
}

// NeedModifyWarmUpVideoTransStatus ...
func NeedModifyWarmUpVideoTransStatus(ctx context.Context, msgVideoMp4CosId string, data *pb.WarmUpData) bool {
	//  暂只支持单个视频
	if len(data.GetWarmupVideoList()) != 1 {
		log.InfoContextf(ctx, "ModifyWarmUpVideoTransStatus, warmUpData video size invalid, size: %+v",
			len(data.GetWarmupVideoList()))
		return false
	}
	cacheWarmUpVideoItem := data.GetWarmupVideoList()[0]
	if cacheWarmUpVideoItem.GetVideoTransStatus() == uint32(util.M3U8TransDone) {
		log.InfoContextf(ctx,"ModifyWarmUpVideoTransStatus cacheWarmUpVideoItem videoTransStatus has trans success, " +
			"cacheWarmUpVideoItem:%+v", cacheWarmUpVideoItem)
		return false
	}
	cacheVideoCosId := cacheWarmUpVideoItem.GetCosId()
	if cacheVideoCosId != msgVideoMp4CosId {
		metrics.IncrCounter("NeedModifyWarmUpVideoTransStatus.cosIdNotSame",1)
		log.InfoContextf(ctx, "ModifyWarmUpVideoTransStatus, cosId not same, from msg tdmq:%+v, from redis:%+v",
			msgVideoMp4CosId, cacheVideoCosId)
		return false
	}
	return true
}
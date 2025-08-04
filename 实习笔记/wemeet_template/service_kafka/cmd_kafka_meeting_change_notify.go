package service_kafka

import (
	"context"
	wlcm "meeting_template/material_control/welcome"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	cachePb "git.code.oa.com/trpcprotocol/wemeet/common_meeting_cache"
	kfPb "git.code.oa.com/trpcprotocol/wemeet/wemeet_kafka_message"
	"github.com/golang/protobuf/proto"
	"meeting_template/util"
)

// HandleMeetingChangeNotify 修改会议
func HandleMeetingChangeNotify(ctx context.Context, msg *kfPb.KafkaMessage) error {
	meetingInfo := &cachePb.MeetingInfo{}
	if err := proto.Unmarshal(msg.GetMeetingInfo(), meetingInfo); err != nil {
		log.ErrorContextf(ctx, "HandleMeetingChangeNotify meetingInfo Unmarshal err:%v, meetInfo:%+v",
			err, msg.GetMeetingInfo())
		metrics.IncrCounter("HandleMeetingChangeNotify.MeetingInfo.Unmarshal.Err", 1)
		return nil
	}

	log.InfoContextf(ctx, "HandleMeetingChangeNotify start, meetingInfo:%+v", meetingInfo)

	// 这里只处理webinar会议的
	if !util.IsWebinarMeeting(meetingInfo.GetBinaryMeetingType()) {
		return nil
	}

	// 更新入会欢迎语过期时间信息缓存过期时间
	err := setWebinarWelcomeInfoExpireTimeHandler(ctx, meetingInfo)
	if err != nil {
		log.ErrorContextf(ctx, "HandleMeetingChangeNotify error, meetingId:%+v, err:%+v", meetingInfo.GetUint64MeetingId(), err)
		return err
	}
	log.InfoContextf(ctx, "HandleMeetingChangeNotify succ, meetingId:%+v", meetingInfo.GetUint64MeetingId())

	return nil
}

// setWebinarWelcomeInfoExpireTimeHandler 设置webinar欢迎语过期时间
func setWebinarWelcomeInfoExpireTimeHandler(ctx context.Context, meetingInfo *cachePb.MeetingInfo) error {
	meetingId := meetingInfo.GetUint64MeetingId()
	orderEndTime := meetingInfo.GetUint32OrderEndTime()
	if meetingId != 0 && orderEndTime > 0 {
		timestamp := uint32(time.Now().Unix())
		expireTime := orderEndTime - timestamp
		if expireTime > 0 {
			expireTime = expireTime + wlcm.WelComeDefaultExpire
		} else {
			expireTime = wlcm.WelComeDefaultExpire
		}
		welCome := wlcm.NewWelComeImp()
		err := welCome.SetWelComeInfoExpireTime(ctx, meetingId, expireTime)
		if err != nil {
			log.ErrorContextf(ctx, "setWebinarWelcomeInfoExpireTimeHandler error, meetingId:%+v, orderEndTime:%+v, err:%+v", meetingId, orderEndTime, err)
			return err
		}
	}
	return nil
}

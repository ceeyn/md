package welcome

import (
	"context"
	"git.code.oa.com/trpc-go/trpc-database/redis"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/golang/protobuf/proto"
	"meeting_template/material_control/cache"
)

// WelComeImp ..
type WelComeImp struct {
	redisProxy redis.Client
}

// NewWelComeImp 初始化
func NewWelComeImp() cache.WelCome {
	imp := &WelComeImp{
		redisProxy: redis.NewClientProxy("trpc.meet.wemeet_webinar.wemeet_webinar_redis"),
	}
	_, err := imp.redisProxy.Do(context.Background(), "ping")
	if err != nil {
		panic(err)
	}

	return imp
}

// GetMeetingWelComeInfo 设置会议欢迎语信息
func (n *WelComeImp) GetMeetingWelComeInfo(ctx context.Context, meetingID uint64) (*pb.WelComeCache, error) {

	// 默认欢迎语
	defWelcomeInfo := &pb.WelComeCache{
		SwitchState:    proto.Uint32(WELCOME_SWITCH_ON),
		DefaultState:   proto.Uint32(WELCOME_DEFAULT_STATE_ON),
		WelcomeTitle:   proto.String(DEFAULT_WELCOME_TITLE),
		WelcomeContent: proto.String(DEFAULT_WELCOME_CONTENT),
	}
	// 会议id为0则查询默认欢迎语信息
	if meetingID == 0 {
		return defWelcomeInfo, nil
	}

	key := MakeWelcomeInfoKey(meetingID)
	rst, err := redis.Bytes(n.redisProxy.Do(ctx, "GET", key))
	if err == redis.ErrNil {
		log.InfoContextf(ctx, "GetMeetingWelComeInfo not get, key:%+v", key)
		return defWelcomeInfo, nil
	} else if err != nil {
		metrics.IncrCounter("GetMeetingWelComeInfo.Err", 1)
		log.ErrorContextf(ctx, "GetMeetingWelComeInfo redis error, key:%+v, err:%+v", key, err.Error())
		return nil, err
	}

	welcomeInfo := &pb.WelComeCache{}
	if err = proto.Unmarshal(rst, welcomeInfo); err != nil {
		metrics.IncrCounter("GetMeetingWelComeInfo.Unmarshal.Err", 1)
		log.ErrorContextf(ctx, " GetMeetingWelComeInfo unmarshal err, rst:%+v, err:%+v", rst, err)
		return nil, err
	}

	metrics.IncrCounter("GetMeetingWelComeInfo.Success", 1)
	log.InfoContextf(ctx, "GetMeetingWelComeInfo succ, meetingId:%+v, welcomeInfo: %+v", meetingID, welcomeInfo)

	return welcomeInfo, nil
}

// SetMeetingWelComeInfo 设置会议欢迎语信息
func (n *WelComeImp) SetMeetingWelComeInfo(ctx context.Context, meetingID uint64, welcomeInfo *pb.WelComeCache, expireTime uint32) error {

	key := MakeWelcomeInfoKey(meetingID)
	value, err := proto.Marshal(welcomeInfo)
	if err != nil {
		metrics.IncrCounter("SetMeetingWelComeInfo.Marshal.Err", 1)
		log.ErrorContextf(ctx, "SetMeetingWelComeInfo marshal fail, welcomeInfo:%+v err:%+v", welcomeInfo, err)
		return err
	}

	_, err = redis.Bytes(n.redisProxy.Do(ctx, "SET", key, value, "EX", expireTime))
	if err != nil {
		metrics.IncrCounter("SetMeetingWelComeInfo.Redis.Set.Err", 1)
		log.ErrorContextf(ctx, "SetMeetingWelComeInfo set redis error, key:%+v, val:%+v, err:%+v", key, value, err)
		return err
	}

	metrics.IncrCounter("SetMeetingWelComeInfo.Success", 1)
	log.InfoContextf(ctx, "SetMeetingWelComeInfo succ, meetingId:%+v, welcomeInfo: %+v", meetingID, welcomeInfo)
	return nil
}

// SetWelComeInfoExpireTime 设置会议欢迎语信息过期时间
func (n *WelComeImp) SetWelComeInfoExpireTime(ctx context.Context, meetingID uint64, expireTime uint32) error {
	key := MakeWelcomeInfoKey(meetingID)
	_, err := redis.Bytes(n.redisProxy.Do(ctx, "GET", key))
	if err == redis.ErrNil {
		log.InfoContextf(ctx, "SetWelComeInfoExpireTime not get, key:%+v", key)
		return nil
	} else if err != nil {
		metrics.IncrCounter("SetWelComeInfoExpireTime.Redis.Get.Err", 1)
		log.ErrorContextf(ctx, "SetWelComeInfoExpireTime redis error, key:%+v, err:%+v", key, err.Error())
		return err
	}

	_, err = n.redisProxy.Do(ctx, "EXPIRE", key, expireTime)
	log.DebugContextf(ctx, "SetWelComeInfoExpireTime key: %+v, expireTime: %+v, err: %+v", key, expireTime, err)
	if err != nil {
		metrics.IncrCounter("SetWelComeInfoExpireTime.Redis.EXPIRE.Err", 1)
		log.ErrorContextf(ctx, "SetWelComeInfoExpireTime error, key:%+v, err:%+v", key, err)
		return err
	}
	return nil
}

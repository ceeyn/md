package rpc

import (
	"context"
	"fmt"

	"meeting_template/util"

	"git.code.oa.com/meettrpc/meet_util"
	"git.code.oa.com/trpc-go/trpc-database/redis"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
)

const (
	// RedisMeetingOtherInfo meeting_other_info 国内redis
	RedisMeetingOtherInfo = "trpc.wemeet.redis.meeting_other_info"
	// RedisOverseaMeetingOtherInfo meeting_other_info 海外redis
	RedisOverseaMeetingOtherInfo = "trpc.wemeet.redis.meeting_other_info_oversea"
	ShowPasswordSubKey           = "show_password"
)

// RDHSetShowPassword 设置展示会议密码
func RDHSetShowPassword(ctx context.Context, isOversea bool,
	meetingID uint64, orderEndTime uint32, showPassword bool) error {
	metrics.IncrCounter("redis.hset.showpassword", 1)
	var proxy redis.Client
	// 国内和海外使用不同的redis实例
	if isOversea {
		proxy = redis.NewClientProxy(RedisOverseaMeetingOtherInfo)
	} else {
		proxy = redis.NewClientProxy(RedisMeetingOtherInfo)
	}
	key := MakeRedisMeetingOtherInfoKey(meetingID)
	_, err := proxy.Do(ctx, "HSET", key, ShowPasswordSubKey, showPassword)
	log.DebugContextf(ctx, "RDHSetShowPassword isOversea:%+v, key:%+v, subKey:%+v, showPassword: %+v, err: %+v",
		isOversea, key, ShowPasswordSubKey, showPassword, err)
	if err != nil {
		metrics.IncrCounter("redis.hset.showpassword.failed", 1)
		log.ErrorContextf(ctx, "RDHSetShowPassword error, cmd = HSET, isOversea:%+v, "+
			"key:%+v, subKey:%+v, showPassword: %v, err: %v", isOversea, key, ShowPasswordSubKey, showPassword, err)
		return err
	}
	metrics.IncrCounter("redis.hset.showpassword.succ", 1)
	log.InfoContextf(ctx, "RDHSetShowPassword isOversea:%+v, key:%+v, subKey:%+v, showPassword: %v",
		isOversea, key, ShowPasswordSubKey, showPassword)

	// 异步设置过期时间
	go func(newCtx context.Context) {
		defer meet_util.DefPanicFun()
		expireTime := util.Get32DaysExpireTimeDuration(orderEndTime)
		_, err = proxy.Do(newCtx, "EXPIRE", key, expireTime)
		log.DebugContextf(ctx, "RDHSetShowPassword isOversea:%+v, key:%+v, "+
			"orderEndTime: %+v, expireTime: %+v, err = %v", isOversea, key, orderEndTime, expireTime, err)
		if err != nil {
			metrics.IncrCounter("redis.expire.showpassword.failed", 1)
			log.Errorf("RDHSetShowPassword error, cmd = EXPIRE,  isOversea:%+v, "+
				"key:%+v, orderEndTime: %+v, expireTime: %+v, err = %v", isOversea, key, orderEndTime, expireTime, err)
		}
	}(trpc.CloneContext(ctx))
	return nil
}

// MakeRedisMeetingOtherInfoKey meeting_other_info key
func MakeRedisMeetingOtherInfoKey(meetingID uint64) string {
	return fmt.Sprintf("meeting_other_info_%+v", meetingID)
}

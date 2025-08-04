package rpc

import (
	"context"
	"errors"
	"fmt"
	"git.code.oa.com/going/attr"
	"git.code.oa.com/meettrpc/meet_util"
	"git.code.oa.com/trpc-go/trpc-database/redis"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	uuid "git.code.oa.com/wesee_ugc/go.uuid"
	"meeting_template/model"
	"meeting_template/util"
)

// 模板信息存储 redis
const (
	RedisTemplateSection  = "trpc.wemeet.wemeet_template.redis_storage"
	KeyTemplateExpireTime = 62 * 24 * 3600 // 会议模板数据过期时间，产品暂定为2个月
	WebinarInfoExpireTime = 12 * 30 * 24 * 3600
	TemplateIdPrefix      = "tpl"
)

// GenTemplateId 模板资源id生成
func GenTemplateId() string {
	// 资源id格式暂定为：tpl_  +  UnixNano(16进制格式)
	return fmt.Sprintf("%s_%v", TemplateIdPrefix, uuid.NewV4().String())
}

// GetTemplateInfoRedis 获取会议模板信息
func GetTemplateInfoRedis(ctx context.Context, templateId string) (*model.TemplateInfo, error) {
	bgTime := util.NowMs()
	util.ReportOne(util.RedisGetTemplateInfoReq) //[GetTemplateInfo]请求
	proxy := redis.NewClientProxy(RedisTemplateSection)
	rst, err := redis.StringMap(proxy.Do(ctx, "HGETALL", templateId))
	if err != nil {
		util.ReportOne(util.RedisGetTemplateInfoFail)
		log.ErrorContextf(ctx, "GetTemplateInfo redis error, cmd = HGETALL, key:%v val:%v err = %v",
			templateId, len(rst), err.Error())
		return nil, err
	}
	if len(rst) == 0 {
		return nil, redis.ErrNil
	}
	//解析；注意字段名和model.TemplateInfo保持一致，如果修改需同步
	templateInfo := &model.TemplateInfo{}
	gotWarmUpData := false
	for k, v := range rst {
		if k == "Sponsor" {
			templateInfo.Sponsor = v
		} else if k == "CoverName" {
			templateInfo.CoverName = v
		} else if k == "CoverUrl" {
			templateInfo.CoverUrl = v
		} else if k == "CoverList" {
			templateInfo.CoverList = v
		} else if k == "WarmUpData" {
			templateInfo.WarmUpData = v
			gotWarmUpData = true
		} else if k == "MeetingId" {
			templateInfo.MeetingId = v
		} else if k == "AppId" {
			templateInfo.AppId = v
		} else if k == "AppUid" {
			templateInfo.AppUid = v
		} else if k == "TemplateId" {
			templateInfo.TemplateId = v
		} else if k == "Description" {
			templateInfo.Description = v
		}
	}
	if !gotWarmUpData {
		templateInfo.WarmUpData = ""
	}
	cost := util.NowMs() - bgTime
	log.DebugContextf(ctx, "GetTemplateInfo redis ok, cmd = HGETALL, key:%v val:%+v templateInfo:%+v, cost:%vms",
		templateId, rst, templateInfo, cost)
	return templateInfo, nil
}

// SetTemplateInfoRedis 新增|更新模板信息
func SetTemplateInfoRedis(ctx context.Context, templateInfo *model.TemplateInfo) error {
	util.ReportOne(util.RedisSetTemplateInfoReq)
	bgTime := util.NowMs()
	key := templateInfo.TemplateId
	proxy := redis.NewClientProxy(RedisTemplateSection)
	// 注意字段名和model.TemplateInfo保持一致，如果修改需同步
	rst, err := redis.String(proxy.Do(ctx, "HMSET", key,
		"Sponsor", templateInfo.Sponsor,
		"CoverName", templateInfo.CoverName,
		"CoverUrl", templateInfo.CoverUrl,
		"Description", templateInfo.Description,
		"CoverList", templateInfo.CoverList,
		"WarmUpData", templateInfo.WarmUpData,
		"MeetingId", templateInfo.MeetingId,
		"AppId", templateInfo.AppId,
		"AppUid", templateInfo.AppUid,
		"TemplateId", templateInfo.TemplateId,
	))

	cost := util.NowMs() - bgTime
	//处理结果
	if err != nil {
		util.ReportOne(util.RedisSetTemplateInfoFail)
		log.ErrorContextf(ctx, "SetTemplateInfo redis error, cmd = HMSET, key:%v TemPlateInfo:%+v,"+
			" err = %v, cost:%vms", key, templateInfo, err.Error(), cost)
		return err
	}
	log.InfoContextf(ctx, "SetTemplateInfo redis ok, cmd = HMSET, key:%v TemPlateInfo:%+v, rst:%v, cost:%vms",
		key, templateInfo, rst, cost)
	return nil
}

// SetExpireRedis 设置过期时间
func SetExpireRedis(ctx context.Context, key string, expireTime int) error {
	util.ReportOne(util.RedisSetExpireReq) //[SetExpire]请求
	proxy := redis.NewClientProxy(RedisTemplateSection)
	rst, err := redis.Int(proxy.Do(ctx, "Expire", key, expireTime))
	if err != nil {
		util.ReportOne(util.RedisSetExpireFail) //[SetExpire]失败
		log.ErrorContextf(ctx,
			"SetExpire redis error, cmd = Expire, key:%v KeyTemplateExpireTime:%v rst:%v err = %v",
			key, expireTime, rst, err.Error())
		return err
	}
	if rst == 1 {
		log.DebugContextf(ctx, "SetExpire  redis ok key:%v ,KeyTemplateExpireTime:%v rst:%v",
			key, expireTime, rst)
		return nil
	}
	util.ReportOne(util.RedisSetExpireFail) //[SetExpire]失败
	log.ErrorContextf(ctx, "SetExpire  redis fail key:%v, KeyTemplateExpireTime:%v rst:%v ",
		key, expireTime, rst)
	return errors.New(fmt.Sprintf("key is not exist"))
}

// 获取待检测的模板id列表
func GetNeedCheckIds(ctx context.Context, needCheckIds *[]string) error {
	proxy := redis.NewClientProxy(RedisTemplateSection)
	// 暂定一次获取10000个模板id用于检测
	rst, err := redis.Strings(proxy.Do(ctx, "ZREVRANGE", "need_check_tpl_infos", 0, 10000))
	if err != nil {
		log.ErrorContextf(ctx, "GetNeedCheckIds redis error, cmd = ZREVRANGE, need_check_tpl_inofos 0, 1000, "+
			"rst:%v err = %v", rst, err)
		return err
	}
	for _, v := range rst {
		*needCheckIds = append(*needCheckIds, v)
	}
	return nil
}

// 从待检测有序集合删除模板id
func RemoveNeedCheck(ctx context.Context, templateId string) error {
	proxy := redis.NewClientProxy(RedisTemplateSection)
	rst, err := proxy.Do(ctx, "ZREM", "need_check_tpl_infos", templateId)
	if err != nil {
		log.ErrorContextf(ctx, "RemoveNeedCheck redis error, cmd = ZREM nneed_check_tpl_infos %v, rst:%v err = %v",
			templateId, rst, err)
		return err
	}
	return nil
}

// RDUpdateInviteSwitch..
func RDUpdateInviteSwitch(ctx context.Context, key string, val uint32) error {
	proxy := redis.NewClientProxy(RedisTemplateSection)
	_, err := redis.String(proxy.Do(ctx, "SET", key, val))
	if err != nil {
		attr.AttrAPI(35765628, 1) //[redis]设置暖场邀请开关失败
		log.ErrorContextf(ctx, "RDUpdateInviteSwitch redis failed, cmd = SET, key:%v, val:%+v, err:%+v", key, val, err)
		return err
	}
	attr.AttrAPI(35765629, 1) //[redis]设置暖场邀请开关成功
	log.InfoContextf(ctx, "RDUpdateInviteSwitch redis ok, cmd = SET, key:%v", key)
	return nil
}

// RDGetInviteSwitch
func RDGetInviteSwitch(ctx context.Context, key string) (int, error) {
	proxy := redis.NewClientProxy(RedisTemplateSection)
	rst, err := redis.Int(proxy.Do(ctx, "GET", key))
	if err == redis.ErrNil { // 此场会议默认值开关
		attr.AttrAPI(35765631, 1)
		log.InfoContextf(ctx, "RDGetInviteSwitch key not exist. key:%+v", key)
		return 1, nil //默认打开邀请开关
	}
	if err != nil {
		attr.AttrAPI(35765632, 1) //[redis]RDGetInviteSwitch失败
		log.ErrorContextf(ctx, "RDGetInviteSwitch redis failed, cmd = GET, key:%+v, val:%+v, err:%+v", key, err)
		return 1, err
	}
	attr.AttrAPI(35765633, 1) //[redis]RDGetInviteSwitch成功
	log.InfoContextf(ctx, "RDGetInviteSwitch redis ok, cmd = GET, key:%+v", key, rst)
	return rst, nil
}

// RDIncrForGetId ...
func RDIncrForGetId(ctx context.Context, key string) (int64, error) {
	attr.AttrAPI(35914645, 1) //[RDIncrForGetId]请求
	proxy := redis.NewClientProxy(RedisTemplateSection)
	rst, err := redis.Int64(proxy.Do(ctx, "INCR", key))
	if err != nil {
		attr.AttrAPI(35914646, 1) //[RDIncrForGetId]失败
		log.ErrorContextf(ctx, "RedisIncrSeq redis error, cmd = INCR, key:%v val:%v err = %v",
			key, rst, err.Error())
		return 0, err
	}
	attr.AttrAPI(35914647, 1) //[RDIncrForGetId]成功
	log.InfoContextf(ctx, "RedisIncrSeq redis ok, cmd = INCR, key:%v  val:%v  rst:%v ", key, rst, rst)
	return rst, err
}

// RDHSetWebinarInfo ...
func RDHSetWebinarInfo(ctx context.Context, key string, subKey uint64, value string) error {
	attr.AttrAPI(35914648, 1) //[RDHSetWebinarInfo]请求
	proxy := redis.NewClientProxy(RedisTemplateSection)
	_, err := proxy.Do(ctx, "HSET", key, subKey, value)
	if err != nil {
		attr.AttrAPI(35914649, 1) //[RDHSetWebinarInfo]请求失败
		log.ErrorContextf(ctx, "RDHSetWebinarInfo Fail, key:%v, err:%+v", key, err)
		return err
	}
	attr.AttrAPI(35914650, 1) //[RDHSetWebinarInfo]请求成功
	log.InfoContextf(ctx, "RDHSetWebinarInfo redis ok, cmd = HSET, key:%v, subKey:%v, val:%v",
		key, subKey, value)
	//异步设置过期时间
	go func(newCtx context.Context) {
		defer meet_util.DefPanicFun()
		_, err = proxy.Do(newCtx, "EXPIRE", key, WebinarInfoExpireTime)
		if err != nil {
			attr.AttrAPI(0, 1) //设置缓存过期时间失败
			log.Errorf("RDHSetWebinarInfo error, cmd = EXPIRE, key:[%v],  err = %v", key, err)
		}
	}(trpc.CloneContext(ctx))
	return nil
}

// RDHLenWebinarInfo ...
func RDHLenWebinarInfo(ctx context.Context, key string) (int64, error) {
	attr.AttrAPI(35914651, 1) //[RDHLenWebinarInfo]请求
	proxy := redis.NewClientProxy(RedisTemplateSection)
	rst, err := redis.Int64(proxy.Do(ctx, "HLEN", key))
	if err != nil {
		attr.AttrAPI(35914652, 1) //[RDHLenWebinarInfo]请求失败
		log.ErrorContextf(ctx, "RDHLenWebinarInfo Fail, key:%v, err:%+v", key, err)
		return 0, err
	}
	attr.AttrAPI(35914653, 1) //[RDHSetWebinarInfo]请求成功
	log.InfoContextf(ctx, "RDHLenWebinarInfo redis ok, cmd = HLEN, key:%v, rst:%v",
		key, rst)
	return rst, nil
}

// RDHDelWebinarInfo 支持批量删除
func RDHDelWebinarInfo(ctx context.Context, key string, fields []string) error {
	attr.AttrAPI(35914654, 1) //[RDHDelWebinarInfo]请求
	proxy := redis.NewClientProxy(RedisTemplateSection)
	args := make([]interface{}, 0, len(fields)+1)
	args = append(args, key)
	for _, filed := range fields {
		args = append(args, filed)
	}
	_, err := proxy.Do(ctx, "HDEL", args...)
	if err != nil {
		attr.AttrAPI(35914655, 1) //[RDHDelWebinarInfo]请求失败
		log.ErrorContextf(ctx, "RDHDelWebinarInfo Fail, key:%v, fields:%+v, err:%+v", key, fields, err)
		return err
	}
	attr.AttrAPI(35914656, 1) //[RDHDelWebinarInfo]请求成功
	log.InfoContextf(ctx, "RDHDelWebinarInfo redis ok, cmd = HDEL, key:%v, fields:%+v",
		key, fields)
	return nil
}

// RDHMSETWebinarInfo ...
func RDHMSETWebinarInfo(ctx context.Context, key string, webinarInfoMap map[uint64]string) error {
	attr.AttrAPI(35914657, 1) //[RDHMSETWebinarInfo]请求
	proxy := redis.NewClientProxy(RedisTemplateSection)
	if len(webinarInfoMap) == 0 {
		metrics.IncrCounter("WebinarInfoMapEmpty", 1)
		log.InfoContextf(ctx, "RDHMSETWebinarInfo webinarInfoMap is empty, key:%+v", key)
		return nil
	}
	//设置缓存信息
	_, err := redis.Bytes(proxy.Do(ctx, "HMSET", redis.Args{}.Add(key).AddFlat(webinarInfoMap)...))
	if err != nil {
		attr.AttrAPI(35914658, 1) //[redis]RDHMSETWebinarInfo失败
		log.ErrorContextf(ctx, "RDHMSETWebinarInfo redis error, cmd = HMSET, key:%v ,err:%v", key, err.Error())
		return errors.New(fmt.Sprintf("RDHMSETWebinarInfo falied, key:%v, err:%v", key, err.Error()))
	}
	attr.AttrAPI(35914659, 1) //[redis]RDHMSETWebinarInfo成功
	log.InfoContextf(ctx, "RDHMSETWebinarInfo succ. key:%+v, webinarInfoMap:%+v", key, webinarInfoMap)
	return nil
}

// RDSetIncrValue ...
func RDSetIncrValue(ctx context.Context, key string, val uint32) error {
	attr.AttrAPI(35914660, 1) //[RDSetIncrValue]请求
	proxy := redis.NewClientProxy(RedisTemplateSection)
	_, err := proxy.Do(ctx, "SET", key, val)
	if err != nil {
		attr.AttrAPI(35914661, 1) //[RDSetIncrValue]请求失败
		log.ErrorContextf(ctx, "RDSetIncrValue Fail, key:%v, val:%+v, err:%+v", key, val, err)
		return err
	}
	attr.AttrAPI(35914662, 1) //[RDSetIncrValue]请求成功
	log.InfoContextf(ctx, "RDSetIncrValue redis ok, cmd = SET, key:%v, val:%+v", key, val)
	//异步设置过期时间
	go func(newCtx context.Context) {
		defer meet_util.DefPanicFun()
		_, err = proxy.Do(newCtx, "EXPIRE", key, WebinarInfoExpireTime)
		if err != nil {
			attr.AttrAPI(0, 1) //设置缓存过期时间失败
			log.Errorf("RDSetIncrValue error, cmd = EXPIRE, key:[%v],  err = %v", key, err)
		}
	}(trpc.CloneContext(ctx))
	return nil
}

// RDHValsWebinarInfo ...
// NOCA:RedisHashSlowCmd(EnsureSafe@yucachen)
func RDHValsWebinarInfo(ctx context.Context, key string) ([]string, error) {
	attr.AttrAPI(35914663, 1) //[RDHValsWebinarInfo]请求
	proxy := redis.NewClientProxy(RedisTemplateSection)
	reply, err := redis.Strings(proxy.Do(ctx, "HVALS", key)) //NOCA:RedisHashSlowCmd(EnsureSafe@yucachen)
	if err != nil {
		attr.AttrAPI(35914664, 1)
		log.ErrorContextf(ctx, "RDHValsWebinarInfo fail with err:%v, key:%v ",
			err, key) //NOCA:RedisHashSlowCmd(EnsureSafe@yucachen)
		return nil, err
	}
	attr.AttrAPI(35914665, 1)
	log.InfoContextf(ctx, "RDHValsWebinarInfo ok, key:%v, reply:%v ",
		key, reply) //NOCA:RedisHashSlowCmd(EnsureSafe@yucachen)
	return reply, nil
}

// RDHGetParticipantInfo ...
func RDHGetParticipantInfo(ctx context.Context, key string, subKey string) (string, error) {
	attr.AttrAPI(35914666, 1) //[RDHGetParticipantInfo]请求
	proxy := redis.NewClientProxy(RedisTemplateSection)
	reply, err := redis.String(proxy.Do(ctx, "HGET", key, subKey))
	if err != nil {
		attr.AttrAPI(35914667, 1) //[redis]RDHGetParticipantInfo
		log.ErrorContextf(ctx, "RDHGetParticipantInfo fail with err:%v, cmd = HGET, key:%v ", err, key)
		return "", err
	}
	attr.AttrAPI(35914668, 1) //[redis]RDHGetParticipantInfo
	log.InfoContextf(ctx, "RDHGetParticipantInfo ok, cmd = HGET, key:%v, reply:%v ", key, reply)
	return reply, nil
}

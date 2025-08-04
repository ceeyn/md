package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"meeting_template/dao"
	"meeting_template/kafka"
	"meeting_template/model"
	"meeting_template/rpc"
	"meeting_template/util"

	"git.code.oa.com/meettrpc/meet_util"
	"git.code.oa.com/trpc-go/trpc-database/redis"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

var (
	getTemplateInfoGroup singleflight.Group
)

// 基础服务类，主要提供一些业务共用方法

// 校验用户是否是会议创建者或超管
func checkCreatorOrSuperAdmin(ctx context.Context, appid, appUid, meetingId string) bool {

	meetingIdUint, err := strconv.ParseUint(meetingId, 10, 64)
	if err != nil {
		return false
	}
	meetingInfo, _, _, _ := rpc.GetMeetingInfo(ctx, meetingIdUint)
	if meetingInfo == nil || meetingInfo.GetUint64MeetingId() == 0 {
		return false
	}

	if appUid == meetingInfo.GetStrCreatorAppUid() {
		return true
	}

	appIdStr := strconv.Itoa(int(meetingInfo.GetUint32CreatorSdkappid()))
	log.DebugContextf(ctx, "print meetingAppId: %+v, curAppId: %+v", appIdStr, appid)
	if appid != appIdStr {
		return false
	}

	isSuperAdmin, err := rpc.CheckIsSuperAdmin(ctx, appid, appUid)
	if err != nil {
		return false
	}
	return isSuperAdmin
}

// GetTemplateInfoSingleFlight 获取会议模板信息SingleFlight
func GetTemplateInfoSingleFlight(ctx context.Context, templateId string) (*model.TemplateInfo, error) {
	v, err, shared := getTemplateInfoGroup.Do(templateId, func() (interface{}, error) {
		util.ReportOne(util.GetTemplateInfoSingleFlightCall) //[GetTemplateInfo]真正执行查询
		go func() {
			meet_util.DefPanicFun()
			forget := 500
			time.Sleep(time.Duration(forget) * time.Millisecond)
			getTemplateInfoGroup.Forget(templateId)
		}()
		return GetTemplateInfo(ctx, templateId)
	})
	if shared {
		util.ReportOne(util.GetTemplateInfoSingleFlightCallShared) //[GetTemplateInfo]复用查询结果
	}
	templateInfo := v.(*model.TemplateInfo)
	return templateInfo, err
}

// GetTemplateInfo 获取会议模板信息
func GetTemplateInfo(ctx context.Context, templateId string) (*model.TemplateInfo, error) {
	if !CheckTemplateId(templateId) {
		util.ReportOne(util.GetTemplateInfoTemplateIdIllegal) //[GetTemplateInfo]templateId不合法，查询失败
		return nil, errors.New("invalid templateId: " + templateId)
	}
	templateInfo, err := rpc.GetTemplateInfoRedis(ctx, templateId)

	if err != nil {
		if err != redis.ErrNil {
			// 1、查询redis报错，返回错误
			util.ReportOne(util.GetTemplateInfoFromRedisFail) //[GetTemplateInfo]从redis中查询失败
			return nil, err
		}
		// 2、redis中没有查到，再查db
		templateInfo, err = dao.GetTemplateInfoByTemplateId(ctx, templateId, nil)
		if err != nil {
			// 2.1 db查询错误
			if err != gorm.ErrRecordNotFound {
				util.ReportOne(util.GetTemplateInfoFromMySQLFail) //[GetTemplateInfo]从MySQL中失败
				return nil, err
			}
			// 2.2 db中没查到，将空TemplateInfo写入redis，防止缓存击穿
			util.ReportOne(util.GetTemplateInfoFromMySQLNotFound) //[GetTemplateInfo]从MySQL中没查到
			templateInfo = &model.TemplateInfo{TemplateId: templateId}
		} else {
			// 2.3 在db中查询到
			util.ReportOne(util.GetTemplateInfoFromMySQLSucc) //[GetTemplateInfo]从MySQL中查询成功
		}
		// 异步templateInfo写入redis，设置过期时间
		go func(ctx context.Context) {
			defer meet_util.DefPanicFun()
			rpc.SetTemplateInfoRedis(ctx, templateInfo)
			rpc.SetExpireRedis(ctx, templateId, rpc.KeyTemplateExpireTime)
		}(trpc.CloneContext(ctx))

		return templateInfo, nil
	} else {
		// 3、从redis中查到
		// todo 这里先和原来保持一致每次查询都会设置一次过期时间。如果后期优化redis访问量，
		// 可以将上次过期时间字段也作为templateInfo字段存储在redis中，这样可以根据上次过期时间来判断本次是否需要再续期
		go func(ctx context.Context) {
			defer meet_util.DefPanicFun()
			expireTime, err := GetTemplateInfoExpireTime(ctx, templateInfo.MeetingId, true)
			if err != nil {
				log.ErrorContextf(ctx, "GetTemplateInfoExpireTime error. "+
					"meetingId: %+v, expireTime: %+v, err: %+v", templateInfo.MeetingId, expireTime, err)
			}
			rpc.SetExpireRedis(ctx, templateId, expireTime)
		}(trpc.CloneContext(ctx))
		util.ReportOne(util.GetTemplateInfoFromRedisSucc) //[GetTemplateInfo]从redis中查询成功

		return templateInfo, nil
	}
}

// SetTemplateInfo 新增/更新模板信息
func SetTemplateInfo(ctx context.Context, templateInfo *model.TemplateInfo) error {
	isUpdate := false
	if templateInfo.TemplateId != "" {
		if !CheckTemplateId(templateInfo.TemplateId) {
			util.ReportOne(util.UpdateTemplateTemplateIdIllegal) //[SetTemplateInfo]templateId不合法
			return errors.New("invalid templateId: " + templateInfo.TemplateId)
		}
		util.ReportOne(util.UpdateTemplateInfoReq) //[SetTemplateInfo update]请求
		isUpdate = true
	} else {
		util.ReportOne(util.CreateTemplateInfoReq) //[SetTemplateInfo create]请求
		templateInfo.TemplateId = rpc.GenTemplateId()
	}

	var err error
	// 先写DBProxy(kafka生产)，再写redis，最大限度保证redis和MySQL一致性
	if isUpdate {
		err = kafka.UpdateTemplateInfoDbProxy(ctx, *templateInfo)
	} else {
		err = kafka.InsertTemplateInfoDbProxy(ctx, *templateInfo)
	}
	if err != nil {
		return err
	}
	// 异步templateInfo写入redis，设置过期时间
	go func(ctx context.Context) {
		defer meet_util.DefPanicFun()
		// 不检查redis是否写入成功（即使redis写入失败，查询时也会从MySQL中查询到再插入redis）
		rpc.SetTemplateInfoRedis(ctx, templateInfo)
		expireTime, err := GetTemplateInfoExpireTime(ctx, templateInfo.MeetingId, isUpdate)
		if err != nil {
			log.ErrorContextf(ctx, "GetTemplateInfoExpireTime error. "+
				"meetingId: %+v, expireTime: %+v, err: %+v", templateInfo.MeetingId, expireTime, err)
		}
		rpc.SetExpireRedis(ctx, templateInfo.TemplateId, expireTime)
	}(trpc.CloneContext(ctx))

	return nil
}

// GetTemplateInfoExpireTime 获取模版信息过期时间
func GetTemplateInfoExpireTime(ctx context.Context, meetingId string, isAfterMeetingEnd bool) (int, error) {
	log.DebugContextf(ctx, "GetTemplateInfoExpireTime meetingId: %+v, isAfterMeetingEnd: %+v",
		meetingId, isAfterMeetingEnd)
	if !isAfterMeetingEnd {
		// 不要求会议结束时间，默认2个月
		return rpc.KeyTemplateExpireTime, nil
	}
	// 要求会议结束时间，会议结束时间+2个月
	meetingIdUint, err := strconv.ParseUint(meetingId, 10, 64)
	if err != nil {
		return rpc.KeyTemplateExpireTime, errors.WithMessage(err, "parse meetingId error")
	}
	// 获取会议预定结束时间
	meetingInfo, _, _, err := rpc.GetMeetingInfo(ctx, meetingIdUint)
	if err != nil {
		return rpc.KeyTemplateExpireTime, errors.WithMessage(err, "rpc.GetMeetingInfo error")
	}
	orderEndTime := meetingInfo.GetUint32OrderEndTime()
	now := util.NowS()
	if orderEndTime < now {
		// 如果会议结束时间小于当前时间，返回默认2个月
		return rpc.KeyTemplateExpireTime, nil
	}
	expireTime := int(orderEndTime-now) + rpc.KeyTemplateExpireTime
	log.DebugContextf(ctx, "GetTemplateInfoExpireTime orderEndTime: %+v, now: %+v, expireTime: %+v",
		orderEndTime, now, expireTime)
	return expireTime, nil
}

// CheckTemplateId 校验templateId
func CheckTemplateId(templateId string) bool {
	if !strings.HasPrefix(templateId, rpc.TemplateIdPrefix) {
		return false
	}
	return true
}

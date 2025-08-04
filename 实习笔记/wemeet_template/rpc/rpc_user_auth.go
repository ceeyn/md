package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	authPb "git.code.oa.com/trpcprotocol/tencent_meeting/common_service_role_auth"
	"github.com/golang/protobuf/proto"
	"meeting_template/util"
)

// CheckUserHasPower ... 判断用户是否有某个权限
func CheckUserHasPower(ctx context.Context, appId uint32, appUid string, powerPoint string) bool {

	req := &authPb.DescribeUserAuthReq{
		CorpId: proto.String(fmt.Sprint(appId)),
		Uid:    proto.String(appUid),
	}

	proxy := authPb.NewRoleAuthOidbClientProxy()

	rsp, err := proxy.DescribeUserAuth(ctx, req)
	if err != nil {
		metrics.IncrCounter("proxy.DescribeUserAuth_fail",1)
		log.ErrorContextf(ctx, "CheckUserHasPower rpc DescribeUserAuth failed. appId:%+v, appUid:%+v, err:%+v",
			appId, appUid, err)
		return true
	}
	metrics.IncrCounter("proxy.DescribeUserAuth_succ",1)
	log.InfoContextf(ctx,"CheckUserHasPower rpc DescribeUserAuth succ, req:%+v, rsp:%+v", req, rsp)

	if rsp.RspCode.GetCode() != "0" {
		metrics.IncrCounter("DescribeUserAuth_RspCode_Not_0",1)
		log.ErrorContextf(ctx, "CheckUserHasPower rpc DescribeUserAuthReq code not 0, req:%+v, rsp:%+v", req, rsp)
		return false
	}
	if len(rsp.GetData().GetAccessList()) == 0 {
		log.ErrorContextf(ctx, "CheckUserHasPower rpc DescribeUserAuthReq AccessList length is 0," +
			" req:%+v, rsp:%+v", req, rsp)
		return false
	}

	for _, v := range rsp.GetData().GetAccessList() {
		if v.GetAuthCode() == powerPoint {
			log.InfoContextf(ctx, "CheckUserHasPower Person have auth. powerPoint:%+v, appId:%+v, " +
				"appUid:%+v, authCode:%+v", powerPoint, appId, appUid, v.GetAuthCode())
			return true
		}
	}
	log.InfoContextf(ctx, "CheckUserHasPower Person dont have auth. powerPoint:%+v, appId:%+v, appUid:%+v",
		powerPoint, appId, appUid)
	return false
}


// WebinarAbility ...
type WebinarAbility struct {
	WebinarAbilityGradientSwitch  bool   `json:"webinar_ability_gradient_switch"`
}

// UserSetting ...
type UserSetting struct {
	WebinarMeeting  WebinarAbility `json:"webinar-meeting"`
}

//GetUserSettings ...
func GetUserSettings(ctx context.Context, appUid string, sdkAppId uint32) bool {
	req := &authPb.GetSettingReq{
		CorpId: proto.String(fmt.Sprint(sdkAppId)),
		Filters: []*authPb.KVs{
			{
				Key:    proto.String("webinar-meeting"),
				Values: []string{"webinar_ability_gradient_switch"},   //是否读取能力项配置
			},
		},
		Uid: proto.String(appUid),
		IsFilterDefault: proto.Bool(true),
	}
	proxy := authPb.NewRoleAuthOidbClientProxy()
	rsp, err := proxy.GetSettings(ctx, req)
	if err != nil || rsp.GetRspCode().GetCode() != "0" {
		metrics.IncrCounter("proxy.GetSettings_fail",1)
		log.ErrorContextf(ctx, "[GetUserSettings]  GetSettings fail, sdkAppId:%+v, appUid:%+v, err:%+v",
			sdkAppId, appUid, err)
		return false    // 降级
	}
	metrics.IncrCounter("proxy.GetSettings_succ",1)
	log.InfoContextf(ctx,"GetUserSettings rpc GetSettings succ. req:%+v, rsp:%+v", req, rsp)
	strSetting := rsp.GetSettings()
	userSettings := &UserSetting{}
	err = json.Unmarshal([]byte(strSetting), userSettings)
	if err != nil {
		metrics.IncrCounter("json.Unmarshal_userSettings_fail",1)
		log.ErrorContextf(ctx,"json.Unmarshal strSetting failed, err:%+v", err)
		return false
	}
	log.InfoContextf(ctx,"GetUserSettings get userSettings:%+v, sdkAppId:%+v, appUid:%+v", userSettings,
		sdkAppId, appUid)
	return userSettings.WebinarMeeting.WebinarAbilityGradientSwitch    // 为true需要去校验权限点
}

// CheckMeetingPermission ... 判断用户是有会议控制权限
func CheckMeetingPermission(ctx context.Context, appId uint32, appUid string, creatorAppId uint32) bool {
	if appId != creatorAppId {
		log.ErrorContextf(ctx, "CheckMeetingPermission appId[%v] != creatorAppId[%v]", appId, creatorAppId)
		return false
	}

	return CheckUserHasPower(ctx, appId, appUid, util.RoleAuthMeetingModification)
}
package rpc

import (
	"context"
	"encoding/base64"
	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	gatePb "git.code.oa.com/trpcprotocol/wemeet/wemeet_safe_gate"
	uuid "git.code.oa.com/wesee_ugc/go.uuid"
)

const WEB = 5

// GetSafe ... 获取安全接口
func GetSafe(ctx context.Context, action string, req *pb.UpdateTemplateReq, url string, scenes string) (int32, string) {
	fromUser := &gatePb.UserInfo{
		AppId: req.GetUint32AppId(),
		AppUid: req.GetStrAppUid(),
		InstanceId: uint32(WEB),    // web
		Ip: req.GetStrUserIp(),   // 用户ip
	}
	meetingInfo := &gatePb.MeetingInfo{
		BinartMeetingType: 1,     // webinar
		MeetingId: req.GetMeetingId(),
		Attrs: map[string]string{
			"webinar_url" : url,
		},
	}
	serInfo := &gatePb.SerInfo{
		AppFrom: "wemeet_template",
		Scenes: scenes,
		TraceId: uuid.NewV4().String(),    // uuid库生成
	}
	safeReq := &gatePb.GetSafeReq{
		Action: action,    //
		FromUser: fromUser,
		Meet: meetingInfo,
		Ser: serInfo,
	}
	proxy := gatePb.NewSafeTrpcClientProxy()
	safeRsp, err := proxy.GetSafe(ctx, safeReq)
	if err != nil {
		attr.AttrAPI(35907532,1)
		log.ErrorContextf(ctx,"rpc GetSafe failed. req:%+v, rsp:%+v, err:%+v", safeReq, safeRsp, err)
		return 0, ""
	}
	attr.AttrAPI(35907533,1)
	log.InfoContextf(ctx, "check User block rpc GetSafe succ, req:%+v, rsp:%+v", safeReq, safeRsp)
	if safeRsp.GetErrorCode() == 0 && safeRsp.GetScore() >= 900 {
		tips, _ := base64.StdEncoding.DecodeString(safeRsp.GetTips())
		attr.AttrAPI(35907534,1)
		log.InfoContextf(ctx,"rpc GetSafe need block, rsp:%+v", safeRsp)
		return 9000, string(tips)
	}
	return 0, ""
}
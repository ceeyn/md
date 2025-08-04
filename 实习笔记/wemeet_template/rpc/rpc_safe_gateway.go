package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"meeting_template/util"

	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	sf "git.code.oa.com/trpcprotocol/wemeet/wemeet_safe_gate"
	safePb "git.code.oa.com/trpcprotocol/wemeet/wemeet_safe_gateway"
	"google.golang.org/protobuf/proto"
)

const (
	GateWayAppFrom = "wemeet_template"
	GateWayScenes  = "sc_webinar_video"
)

// SafetyAuditReq 校验的参数
type SafetyAuditReq struct {
	Uuid              string `json:"uuid"`
	AppId             uint32 `json:"app_id"`
	AppUid            string `json:"app_uid"`
	MeetingId         uint64 `json:"meeting_id"`
	Url               string `json:"url"`
	Scenes            string `json:"scenes"`
	Action            string `json:"action"`
	CallbackProtoType string `json:"callback_proto_type"` // 异步回调必填，回调协议类型固定值：http
	CallbackTarget    string `json:"callback_target"`     //异步回调必填，回调地址。例如：polaris://trpc.wecalendar.calendar.api 或者 ip://30.42.53.101:9000
	CallbackUri       string `json:"callback_uri"`        //异步回调必填，回调URI。例如：/safe/v1/get_text_callback
	CallbackNameSpace string `json:"callback_name_space"` //异步回调必填，
	CallbackEnvName   string `json:"callback_env_name"`   //异步回调必填，
}

// WarmVideoSafetyAudit 视频送检, 送检成功返回true, 失败返回false  get_webinar_warm_video_check
func WarmVideoSafetyAudit(ctx context.Context, warmReq *SafetyAuditReq) {
	util.ReportOne(util.GetVidoReq)
	warmReqBytes, err := json.Marshal(warmReq)
	if err != nil {
		metrics.IncrCounter("Rpc.GetSafe.Err", 1)
		log.ErrorContext(ctx, "rpc GetImage json marsha fail, err:%v", err)
		//作为附带的str，json失败了，也没有关系
	}
	warmReqString := string(warmReqBytes)

	req := &sf.GetSafeReq{
		Action: warmReq.Action,
		FromUser: &sf.UserInfo{
			AppId:  warmReq.AppId,
			AppUid: warmReq.AppUid,
			Attrs: map[string]string{
				"video_url": warmReq.Url,
				"video_act": util.VideoAuditUploadAction,
			},
		},
		Meet: &sf.MeetingInfo{
			MeetingId:         warmReq.MeetingId,
			BinartMeetingType: 1, //1-webinar会议
		},
		Ser: &sf.SerInfo{
			AppFrom: GateWayAppFrom,
			Scenes:  warmReq.Scenes,
			TraceId: warmReq.Uuid,
			EchoStr: warmReqString,
			Callback: &sf.CallbackInfo{
				ProtoType: warmReq.CallbackProtoType,
				Targe:     warmReq.CallbackTarget,
				Uri:       warmReq.CallbackUri,
				NameSpace: warmReq.CallbackNameSpace,
				EnvName:   warmReq.CallbackEnvName,
			},
		},
	}

	// 发送请求
	metrics.IncrCounter("Rpc.GetSafe", 1)
	proxy := sf.NewSafeTrpcClientProxy()
	rsp, err := proxy.GetSafe(ctx, req)
	if err != nil {
		metrics.IncrCounter("Rpc.GetSafe.Err", 1)
		util.ReportOne(util.GetVidoReqFail) //[GetVidoReq请求失败]
		log.ErrorContext(ctx, "WarmVideoSafetyAudit fail  req:%v, err:%v", req, err)
		return
	}
	log.InfoContextf(ctx, "WarmVideoSafetyAudit ok, \nreq:%v\nrsp:%v", req, rsp)
	return
}

// VideoSafetyAudit 视频送检, 送检成功返回true, 失败返回false
func VideoSafetyAudit(ctx context.Context, uuid string, appId uint32,
	appUid string, meetingId uint64, url string) {
	util.ReportOne(util.GetVidoReq)
	proxy := safePb.NewVideoTrpcClientProxy()

	req := &safePb.GetVideoReq{
		Url: proto.String(url),
		Act: proto.String(util.VideoAuditUploadAction),
	}

	req.User = &safePb.TUserInfo{}
	req.User.AppId = proto.Uint32(appId)
	req.User.AppUid = proto.String(appUid)

	req.Meet = &safePb.TMeetingInfo{}
	req.Meet.MeetingId = proto.Uint64(meetingId)

	req.Ser = &safePb.TSerInfo{}
	req.Ser.AppFrom = proto.String(GateWayAppFrom)
	req.Ser.Scenes = proto.String(GateWayScenes)
	req.Ser.TraceId = proto.String(uuid)
	rsp, err := proxy.GetVideo(ctx, req)

	if err != nil {
		util.ReportOne(util.GetVidoReqFail) //[GetVidoReq请求失败]
		log.ErrorContext(ctx, "VideoSafetyAudit fail  req:%v, err:%v", req, err)
		return
	}
	log.InfoContextf(ctx, "VideoSafetyAudit ok, \nreq:%v\nrsp:%v", req, rsp)
	return
}

// CheckHasSensitiveWords ...
func CheckHasSensitiveWords(ctx context.Context, meetingId uint64, appId uint32, appUid string,
	word string, oldScenes string) (bool, error) {
	log.DebugContextf(ctx, "template-CheckHasSensitiveWords check begin. scenes:%+v,word :%+v, meetingId:%+v",
		oldScenes, word, meetingId)
	if word == "" {
		return false, nil
	}

	//根据旧场景，确定新的action和scenes
	action, scenes := getActionAndSensesByOldScenes(oldScenes)
	log.DebugContextf(ctx, "template-CheckHasSensitiveWords check scenes. scenes:%+v,after check :%+v, action:%+v",
		oldScenes, scenes, action)
	if action != "" && scenes != "" {
		metrics.IncrCounter("CheckHasSensitiveWords.New", 1)
		return CheckHasSensitiveWordsNew(ctx, meetingId, appId, appUid, word, scenes, action)
	}
	metrics.IncrCounter("CheckHasSensitiveWords.Old", 1)
	req := &safePb.GetTextReq{
		Text: proto.String(word),
		User: &safePb.TUserInfo{
			AppId:  proto.Uint32(appId),
			AppUid: proto.String(appUid),
		},
		Ser: &safePb.TSerInfo{
			AppFrom: proto.String(fmt.Sprintf("%v", "wemeet_template")),
			Scenes:  proto.String(oldScenes),
			TraceId: proto.String(fmt.Sprint(rand.Uint32())),
		},
	}
	if meetingId != 0 {
		req.Meet = &safePb.TMeetingInfo{
			MeetingId:  proto.Uint64(meetingId),
			BinaryType: proto.Uint32(1), // webinar
		}
	}

	// 发送请求
	attr.AttrAPI(35914642, 1) // [CheckHasSensitiveWords]查询敏感词信息请求
	proxy := safePb.NewTextOidbClientProxy()
	rsp, err := proxy.GetText(ctx, req)

	// 处理返回内容
	if err != nil {
		attr.AttrAPI(35914643, 1) // [CheckHasSensitiveWords]查询敏感词信息失败
		log.ErrorContextf(ctx, "(CheckHasSensitiveWords) rpc failed. req:%+v, meetingId:%+v, err:%+v",
			req, meetingId, err)
		return false, err
	}
	log.InfoContextf(ctx, "(CheckHasSensitiveWords) rpc succ. req:%+v, rsp:%+v", req, rsp)

	if rsp.GetResultCode() == 1 { // 有敏感词
		attr.AttrAPI(35914644, 1)
		log.InfoContextf(ctx, "meetingId:%+v, have sensitive words: %+v", meetingId, rsp.GetHitWords())
		return true, nil
	}

	return false, nil
}

// ImgSafetyAuditV2 图片送审，将会有回调处理结果,更换调用信安的接口，后续稳定以后，删除ImgSafetyAudit方法
func ImgSafetyAuditV2(ctx context.Context, safeReq *SafetyAuditReq) {
	safeReqBytes, err := json.Marshal(safeReq)
	if err != nil {
		//作为附带的str，json失败了，也没有关系
		log.ErrorContext(ctx, "rpc GetImage json marsha fail, err:%v", err)
	}
	safeReqString := string(safeReqBytes)

	req := &sf.GetSafeReq{
		Action: safeReq.Action,
		FromUser: &sf.UserInfo{
			AppId:  safeReq.AppId,
			AppUid: safeReq.AppUid,
			Attrs: map[string]string{
				"image_url": safeReq.Url,
			},
		},
		Meet: &sf.MeetingInfo{
			MeetingId:         safeReq.MeetingId,
			BinartMeetingType: 1, //1-webinar会议
		},
		Ser: &sf.SerInfo{
			AppFrom: GateWayAppFrom,
			Scenes:  safeReq.Scenes,
			TraceId: safeReq.Uuid,
			EchoStr: safeReqString,
			Callback: &sf.CallbackInfo{
				ProtoType: safeReq.CallbackProtoType,
				Targe:     safeReq.CallbackTarget,
				Uri:       safeReq.CallbackUri,
				NameSpace: safeReq.CallbackNameSpace,
				EnvName:   safeReq.CallbackEnvName,
			},
		},
	}

	// 发送请求
	metrics.IncrCounter("Rpc.GetSafe", 1)
	proxy := sf.NewSafeTrpcClientProxy()
	rsp, err := proxy.GetSafe(ctx, req)
	if err != nil {
		metrics.IncrCounter("Rpc.GetSafe.Err", 1)
		log.ErrorContext(ctx, "ImgSafetyAuditV2 fail req:%v, err:%v", req, err)
		return
	}
	log.InfoContextf(ctx, "ImgSafetyAuditV2 ok,req:%v,rsp:%v", req, rsp)
	return
}

// ImgSafetyAudit 图片送审，将会有回调处理结果
func ImgSafetyAudit(ctx context.Context, uuid string, appId uint32, appUid string, meetingId uint64, url string,
	scence string)  {
	proxy := safePb.NewImageTrpcClientProxy()
	user := &safePb.TUserInfo{
		AppId: proto.Uint32(appId),
		AppUid: proto.String(appUid),
	}
	meet := &safePb.TMeetingInfo{
		MeetingId: proto.Uint64(meetingId),
		BinaryType: proto.Uint32(1),       // webinar
	}
	imgInfo := &safePb.TImageInfo{
		Url: proto.String(url),        // url
	}
	ser := &safePb.TSerInfo {
		AppFrom: proto.String(GateWayAppFrom),
		Scenes: proto.String(scence),
		TraceId: proto.String(uuid),
	}
	req := &safePb.GetImageReq{
		User: user,
		Meet: meet,
		Ser: ser,
		Image: imgInfo,
	}
	rsp, err := proxy.GetImage(ctx, req)
	if err != nil {
		attr.AttrAPI(35907551,1)
		log.ErrorContextf(ctx, "rpc ImgSafetyAudit fail  req:%+v, err:%+v", req, err)
		return
	}
	attr.AttrAPI(35907552,1)
	log.InfoContextf(ctx, "rpc ImgSafetyAudit ok, req:%+v, rsp:%+v", req, rsp)
	return
}
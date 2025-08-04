package rpc

import (
	"context"
	"encoding/json"

	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	sf "git.code.oa.com/trpcprotocol/wemeet/wemeet_safe_gate"
	safePb "git.code.oa.com/trpcprotocol/wemeet/wemeet_safe_gateway"
	"github.com/golang/protobuf/proto"
)

const GateWayAppFrom = "wemeet_template"

// ImgSafetyAuditReq 图片安全的Req
type ImgSafetyAuditReq struct {
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

// ImgSafetyAuditV2 图片送审，将会有回调处理结果,更换调用信安的接口，后续稳定以后，删除ImgSafetyAudit方法
func ImgSafetyAuditV2(ctx context.Context, imgReq *ImgSafetyAuditReq) error {
	imgReqBytes, err := json.Marshal(imgReq)
	if err != nil {
		//作为附带的str，json失败了，也不影响调用逻辑
		log.ErrorContext(ctx, "rpc GetImage json marsha fail, err:%v", err)
	}
	imgReqString := string(imgReqBytes)
	req := &sf.GetSafeReq{
		Action: imgReq.Action,
		FromUser: &sf.UserInfo{
			AppId:  imgReq.AppId,
			AppUid: imgReq.AppUid,
			Attrs: map[string]string{
				"image_url": imgReq.Url,
			},
		},
		Meet: &sf.MeetingInfo{
			MeetingId:         imgReq.MeetingId,
			BinartMeetingType: 1, //1-webinar会议
		},

		Ser: &sf.SerInfo{
			AppFrom: GateWayAppFrom,
			Scenes:  imgReq.Scenes,
			TraceId: imgReq.Uuid,
			EchoStr: imgReqString,
			Callback: &sf.CallbackInfo{
				ProtoType: imgReq.CallbackProtoType,
				Targe:     imgReq.CallbackTarget,
				Uri:       imgReq.CallbackUri,
				NameSpace: imgReq.CallbackNameSpace,
				EnvName:   imgReq.CallbackEnvName,
			},
		},
	}

	// 发送请求
	metrics.IncrCounter("Rpc.GetSafe", 1)
	proxy := sf.NewSafeTrpcClientProxy()
	rsp, err := proxy.GetSafe(ctx, req)
	if err != nil {
		metrics.IncrCounter("Rpc.GetSafe.Err", 1)
		log.ErrorContext(ctx, "rpc GetImage ImgSafetyAudit fail  req:%v, err:%v", req, err)
		return err
	}
	attr.AttrAPI(35907552, 1)
	log.InfoContextf(ctx, "rpc.ImgSafetyAuditV2 ok, req:%+v, rsp:%v", req, rsp)
	return nil
}

// ImgSafetyAudit 图片送审，将会有回调处理结果
func ImgSafetyAudit(ctx context.Context, uuid string, appId uint32, appUid string, meetingId uint64, url string,
	scenes string) {
	proxy := safePb.NewImageTrpcClientProxy()
	user := &safePb.TUserInfo{
		AppId:  proto.Uint32(appId),
		AppUid: proto.String(appUid),
	}
	meet := &safePb.TMeetingInfo{
		MeetingId:  proto.Uint64(meetingId),
		BinaryType: proto.Uint32(1), // webinar
	}
	imgInfo := &safePb.TImageInfo{
		Url: proto.String(url), // url
	}
	ser := &safePb.TSerInfo{
		AppFrom: proto.String(GateWayAppFrom),
		Scenes:  proto.String(scenes),
		TraceId: proto.String(uuid),
	}
	req := &safePb.GetImageReq{
		User:  user,
		Meet:  meet,
		Ser:   ser,
		Image: imgInfo,
	}
	rsp, err := proxy.GetImage(ctx, req)
	if err != nil {
		attr.AttrAPI(35907551, 1)
		log.ErrorContextf(ctx, "rpc GetImage ImgSafetyAudit fail  req:%+v, err:%+v", req, err)
		return
	}
	attr.AttrAPI(35907552, 1)
	log.InfoContextf(ctx, "rpc.ImgSafetyAudit ok, req:%+v, rsp:%+v", req, rsp)
	return
}

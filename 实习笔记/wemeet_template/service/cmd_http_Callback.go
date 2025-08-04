package service

import (
	"context"
	"encoding/json"
	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	safePb "git.code.oa.com/trpcprotocol/wemeet/wemeet_safe_gateway"
	"github.com/golang/protobuf/proto"
	ctlrpc "meeting_template/material_control/rpc"
	"meeting_template/util"
	"time"
)

// VirtualBackgroundSafeCallback 虚拟背景图片审核回调
func (s *WemeetMeetingTemplateHttpServiceImpl) VirtualBackgroundSafeCallback(ctx context.Context, req *pb.SafeCallbackReq, rsp *pb.SafeCallbackRsp) (err error) {
	log.DebugContextf(ctx, "VirtualBackgroundSafeCallback: req:%+v", req)
	metrics.IncrCounter("VirtualBackgroundSafeCallback Total", 1)
	start := time.Now()

	handleReq, err := convert2GetImageCallbackReq(ctx, req)
	if err != nil {
		metrics.IncrCounter("VirtualBackgroundSafeCallback fail", 1)
		log.ErrorContextf(ctx, "VirtualBackgroundSafeCallback invalid param, err:%+v", err)
		rsp.ErrorCode = proto.Int32(util.InvalidParam)
		rsp.ErrorMsg = proto.String(err.Error())
		return nil
	}
	log.DebugContextf(ctx, "VirtualBackgroundSafeCallback,req:%+v", handleReq)
	rst, err := HandleVirtualBackgroundSafeCallback(ctx, handleReq)
	if err != nil {
		metrics.IncrCounter("VirtualBackgroundSafeCallback fail", 1)
		rsp.ErrorCode = proto.Int32(rst)
		rsp.ErrorMsg = proto.String(err.Error())
	} else {
		metrics.IncrCounter("VirtualBackgroundSafeCallback success", 1)
		rsp.ErrorCode = proto.Int32(0)
		rsp.ErrorMsg = proto.String("ok")
	}
	log.InfoContextf(ctx, "http_VirtualBackgroundSafeCallback, cost:%v, rst:%v, err:%+v, req:%+v, rsp:%+v",
		time.Since(start), rst, err, req, rsp)
	return nil
}

// ImgCensorCallback 暖场图片审核回调，用于获取暖场图片审核结果
func (s *WemeetMeetingTemplateHttpServiceImpl) ImgCensorCallback(ctx context.Context, req *pb.SafeCallbackReq, rsp *pb.SafeCallbackRsp) (err error) {
	log.DebugContextf(ctx, "ImgCensorCallback,req:%+v", req)
	handleReq, err := convert2GetImageCallbackReq(ctx, req)
	if err != nil {
		log.ErrorContextf(ctx, "ImgCensorCallback invalid param, err:%+v", err)
		rsp.ErrorCode = proto.Int32(util.InvalidParam)
		rsp.ErrorMsg = proto.String(err.Error())
		return nil
	}

	log.DebugContextf(ctx, "ImgCensorCallback,req:%+v", handleReq)

	rst, err := HandleImgCensorResult(ctx, handleReq)
	rsp.ErrorCode = proto.Int32(rst)
	if err != nil {
		attr.AttrAPI(35907579, 1)
		rsp.ErrorMsg = proto.String(err.Error())
		log.ErrorContextf(ctx, "http_ImgCensorCallback fail req:%+v, rsp:%+v", req, rsp)
	} else {
		attr.AttrAPI(35907580, 1)
		rsp.ErrorMsg = proto.String("ok")
		log.InfoContextf(ctx, "http_ImgCensorCallback ok  req:%+v, rsp:%+v", req, rsp)
	}
	return nil
}

// AvatarSafeCallback webinar活动页嘉宾头像图片审核
func (s *WemeetMeetingTemplateHttpServiceImpl) AvatarSafeCallback(ctx context.Context, req *pb.SafeCallbackReq, rsp *pb.SafeCallbackRsp) (err error) {
	log.DebugContextf(ctx, "AvatarSafeCallback,req:%+v", req)
	handleReq, err := convert2GetImageCallbackReq(ctx, req)
	if err != nil {
		log.ErrorContextf(ctx, "AvatarSafeCallback invalid param, err:%+v", err)
		rsp.ErrorCode = proto.Int32(util.InvalidParam)
		rsp.ErrorMsg = proto.String(err.Error())
		return nil
	}
	log.DebugContextf(ctx, "AvatarSafeCallback,req:%+v", handleReq)
	rst, err := HandleAvatarSafeCallback(ctx, handleReq)
	rsp.ErrorCode = proto.Int32(rst)
	if err != nil {
		attr.AttrAPI(35914631, 1)
		rsp.ErrorMsg = proto.String(err.Error())
		log.ErrorContextf(ctx, "http_AvatarSafeCallback fail req:%+v, rsp:%+v", req, rsp)
	} else {
		attr.AttrAPI(35914632, 1)
		rsp.ErrorMsg = proto.String("ok")
		log.InfoContextf(ctx, "http_AvatarSafeCallback ok  req:%+v, rsp:%+v", req, rsp)
	}
	return nil
}

// VideoCensorCallback 暖场视频审核
func (s *WemeetMeetingTemplateHttpServiceImpl) VideoCensorCallback(ctx context.Context, req *pb.SafeCallbackReq, rsp *pb.SafeCallbackRsp) (err error) {
	log.DebugContextf(ctx, "VideoCensorCallback,req:%+v", req)
	bgTime := util.NowMs()
	util.ReportOne(util.VideoCensorCallback) //[VideoCensorCallback]请求
	handleReq, err := convert2GetVideoCallbackReq(req)
	if err != nil {
		util.ReportOne(util.VideoCensorCallbackFail) //[VideoCensorCallback]请求失败
		rsp.ErrorCode = proto.Int32(1)
		rsp.ErrorMsg = proto.String(err.Error())
		log.ErrorContextf(ctx, "VideoCensorCallback invalid param, err:%+v", err)
		return nil
	}
	log.DebugContextf(ctx, "VideoCensorCallback,req:%+v", handleReq)

	rst, err := HandleVideoCensorResult(ctx, handleReq)
	rsp.ErrorCode = proto.Int32(rst)
	cost := util.NowMs() - bgTime
	if err != nil {
		util.ReportOne(util.VideoCensorCallbackFail) //[VideoCensorCallback]请求失败
		rsp.ErrorMsg = proto.String(err.Error())
		log.ErrorContextf(ctx, "http_VideoCentsorCallback fail req:%v, rsp:%v, cost:%vms", req, rsp, cost)
	} else {
		rsp.ErrorMsg = proto.String("ok")
		log.InfoContextf(ctx, "http_VideoCentsorCallback ok  req:%v, rsp:%v, cost:%vms", req, rsp, cost)
	}
	return nil
}

func convert2GetImageCallbackReq(ctx context.Context, req *pb.SafeCallbackReq) (*safePb.GetImageCallbackReq, error) {
	// 初始化user, meet
	user, meet, err := iniUserAndMeet(req)
	if err != nil {
		return nil, err
	}

	//是否是敏感图片的判断
	resultCode := iniResultCode(req)

	oidbReq := &safePb.GetImageCallbackReq{
		ResultCode: proto.Uint32(resultCode),
		Ser: &safePb.TSerInfo{
			AppFrom: proto.String(req.GetSer().GetAppFrom()),
			Scenes:  proto.String(req.GetSer().GetScenes()),
			TraceId: proto.String(req.GetSer().GetTraceId()),
		},
		User: user,
		Meet: meet,
	}
	return oidbReq, nil
}

func convert2GetVideoCallbackReq(req *pb.SafeCallbackReq) (*safePb.GetVideoCallbackReq, error) {
	// 初始化user, meet
	user, meet, err := iniUserAndMeet(req)
	if err != nil {
		return nil, err
	}

	//是否是敏感视频的判断
	resultCode := iniResultCode(req)

	handleReq := &safePb.GetVideoCallbackReq{
		ResultCode: proto.Uint32(resultCode),
		Ser: &safePb.TSerInfo{
			AppFrom: proto.String(req.GetSer().GetAppFrom()),
			Scenes:  proto.String(req.GetSer().GetScenes()),
			TraceId: proto.String(req.GetSer().GetTraceId()),
		},
		User: user,
		Meet: meet,
	}
	return handleReq, nil
}

// iniResultCode 确定是否是安全的图片，或者视频
func iniResultCode(req *pb.SafeCallbackReq) uint32 {
	if req.GetScore() >= 900 {
		return 1
	} else {
		return 0
	}
}

// iniUserAndMeet 初始化user, meet,
func iniUserAndMeet(req *pb.SafeCallbackReq) (*safePb.TUserInfo, *safePb.TMeetingInfo, error) {
	user := &safePb.TUserInfo{}
	meet := &safePb.TMeetingInfo{}
	if req.GetSer() != nil && req.GetSer().GetEchoStr() != "" {
		imgReq := &ctlrpc.ImgSafetyAuditReq{}
		err := json.Unmarshal([]byte(req.GetSer().GetEchoStr()), imgReq)
		if err != nil {
			return user, meet, err
		}
		user.AppId = proto.Uint32(imgReq.AppId)
		user.AppUid = proto.String(imgReq.AppUid)
		meet.MeetingId = proto.Uint64(imgReq.MeetingId)
		meet.BinaryType = proto.Uint32(util.WebinarMeetingType)
	}
	return user, meet, nil
}

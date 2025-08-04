package service

import (
	"context"
	"time"

	"meeting_template/util"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/golang/protobuf/proto"
	wlcm "meeting_template/material_control/welcome"
)

// GetWelcomeInfo ...
func (s *WemeetMeetingTemplateOidbServiceImpl) GetWelcomeInfo(ctx context.Context, req *pb.GetWelcomeInfoReq,
	rsp *pb.GetWelcomeInfoRsp) error {

	bgTime := util.NowMs()

	rst, err := HandleGetWelcomeInfo(ctx, req, rsp)
	rsp.ErrorCode = proto.Int32(rst)
	cost := util.NowMs() - bgTime
	if err != nil {
		metrics.IncrCounter("GetWelcomeInfo.Failed", 1) //[GetWelcomeInfo]请求失败
		rsp.ErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "GetWelcomeInfo fail req:%v, rsp:%+v, cost:%+vms", req, rsp, cost)
	} else {
		rsp.ErrorMessage = proto.String("ok")
		log.InfoContextf(ctx, "GetWelcomeInfo ok  req:%+v, rsp:%+v, cost:%+vms", req, rsp, cost)
	}
	metrics.IncrCounter("GetWelcomeInfo.Success", 1)
	return nil
}

// HandleGetWelcomeInfo ...
func HandleGetWelcomeInfo(ctx context.Context, req *pb.GetWelcomeInfoReq,
	rsp *pb.GetWelcomeInfoRsp) (ret int32, err error) {

	timestamp := uint32(time.Now().Unix())
	rsp.Timestamp = proto.Uint32(timestamp)
	rsp.WelcomeData = new(pb.WelComeInfo)

	welCome := wlcm.NewWelComeImp()
	welcomeData, err := welCome.GetMeetingWelComeInfo(ctx, req.GetMeetingId())
	log.InfoContextf(ctx, "GetMeetingWelComeInfo, req:%+v, rsp:%+v", req, welcomeData)
	if err != nil {
		metrics.IncrCounter("GetMeetingWelComeInfo.Err", 1)
		rsp.ErrorMessage = proto.String(err.Error())
		return -1, err // redis读取失败，返回失败
	}

	rsp.ErrorCode = proto.Int32(0)
	rsp.SwitchState = proto.Uint32(welcomeData.GetSwitchState())
	// 如果开关打开则返回欢迎语信息
	if wlcm.IsSwitchOpen(welcomeData.GetSwitchState()) {
		if welcomeData.GetDefaultState() == wlcm.WELCOME_DEFAULT_STATE_ON {
			// 返回默认欢迎语
			rsp.WelcomeData.WelcomeTitle = proto.String(wlcm.DEFAULT_WELCOME_TITLE)
			rsp.WelcomeData.WelcomeContent = proto.String(wlcm.DEFAULT_WELCOME_CONTENT)
		} else {
			// 返回自定义欢迎语
			rsp.WelcomeData.WelcomeTitle = proto.String(welcomeData.GetWelcomeTitle())
			rsp.WelcomeData.WelcomeContent = proto.String(welcomeData.GetWelcomeContent())
		}
	}
	return ret, err
}

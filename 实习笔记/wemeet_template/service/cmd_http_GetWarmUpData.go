package service

import (
	"context"
	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/golang/protobuf/proto"
	"meeting_template/util"
)

// GetWarmUpData ...
func (s *WemeetMeetingTemplateHttpServiceImpl) GetWarmUpData(ctx context.Context, req *pb.GetWarmUpDataReq,
	rsp *pb.GetWarmUpDataRsp) error {
	bgTime := util.NowMs()
	util.ReportOne(util.GetWarmUpDataReq) //[GetWarmUpData]请求

	rst, err := HandleGetWarmUpData(ctx, req, rsp)
	rsp.ErrorCode = proto.Int32(rst)
	cost := util.NowMs() - bgTime
	if err != nil {
		util.ReportOne(util.GetWarmUpDataFail) //[GetWarmUpData]请求失败
		rsp.ErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "http GetWarmUpData fail req:%v, rsp:%v, cost:%vms", req, rsp, cost)
	} else {
		rsp.ErrorMessage = proto.String("ok")
		log.InfoContextf(ctx, "http GetWarmUpData ok  req:%v, rsp:%v, cost:%vms", req, rsp, cost)
	}
	return nil
}

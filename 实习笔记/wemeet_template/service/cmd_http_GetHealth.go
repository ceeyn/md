package service

import (
	"context"

	"meeting_template/util"

	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/golang/protobuf/proto"
)

// WemeetMeetingTemplateHttpServiceImpl ...
type WemeetMeetingTemplateHttpServiceImpl struct{}

// GetHealth ...
func (s *WemeetMeetingTemplateHttpServiceImpl) GetHealth(ctx context.Context, req *pb.GetHealthReq,
	rsp *pb.GetHealthRsp) error {
	util.ReportOne(util.GetHealthReq) //[GetHealth]请求

	rsp.Status = proto.String("ok")
	log.InfoContextf(ctx, " Handle GetHealth ok  req:%v, rsp:%v", req, rsp)

	return nil
}

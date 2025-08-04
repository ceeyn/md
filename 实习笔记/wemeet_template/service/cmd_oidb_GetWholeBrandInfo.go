package service

import (
	"context"
	"meeting_template/rpc"
	"meeting_template/util"

	"git.code.oa.com/trpc-go/trpc-go/log"
	cachePb "git.code.oa.com/trpcprotocol/wemeet/common_meeting_cache"
	errpb "git.code.oa.com/trpcprotocol/wemeet/common_xcast_meeting_error_code"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/golang/protobuf/proto"
)

// GetTemplateInfo ...
func (s *WemeetMeetingTemplateOidbServiceImpl) GetWholeBrandInfo(ctx context.Context,
	req *pb.GetWholeBrandInfoReq, rsp *pb.GetWholeBrandInfoRsp) error {
	bgTime := util.NowMs()
	util.ReportOne(util.GetWholeBrandInfoReqAttr) //[GetWholeBrandInfo]请求

	rst, err := HandleGetWholeBrandInfo(ctx, req, rsp)
	rsp.Int32ErrorCode = proto.Int32(rst)
	cost := util.NowMs() - bgTime
	if err != nil {
		util.ReportOne(util.GetWholeBrandInfoFailAttr) //[GetWholeBrandInfo]失败
		rsp.StrErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "GetWholeBrandInfo fail req:%v, rsp:%v, cost:%vms",
			req, rsp, cost)
	} else {
		util.ReportOne(util.GetWholeBrandInfoSuccAttr) //[GetWholeBrandInfo]成功
		rsp.StrErrorMessage = proto.String("ok")
		log.InfoContextf(ctx, "GetWholeBrandInfo ok  req:%v, rsp:%v, cost:%vms",
			req, rsp, cost)
	}
	return nil
}

// HandleGetWholeBrandInfo
func HandleGetWholeBrandInfo(ctx context.Context, req *pb.GetWholeBrandInfoReq,
	rsp *pb.GetWholeBrandInfoRsp) (ret int32, err error) {

	//获取meeting info，拿到template id
	meetingInfo, _, _, err := rpc.GetMeetingInfo(ctx, req.GetUint64MeetingId())
	if err != nil {
		util.ReportOne(util.HandleGetWholeBrandInfoGetMeetingInfoFailAttr) //[HandleGetWholeBrandInfo]获取会议信息失败
		log.ErrorContextf(ctx, "get meeting info from grocery fail, meetingId:[%v], err: %v",
			req.GetUint64MeetingId(), err)
		return int32(errpb.ERROR_CODE_WEMEET_TEMPLATE_NOT_GET_MEETING_INFO), err
	}

	//获取品牌信息
	templateInfo, ret, err := DoGetTemplateInfo(ctx, meetingInfo.GetStrMeetingTemplateId(), "")
	if err != nil {
		util.ReportOne(util.HandleGetWholeBrandInfoGetTemplateFailAttr) //[HandleGetWholeBrandInfo]获取模板信息失败
		log.ErrorContextf(ctx, "DoGetTemplateInfo fail, meetingId:[%v], template_id:[%v], err: %v",
			req.GetUint64MeetingId(), meetingInfo.GetStrMeetingTemplateId(), err)
		return int32(errpb.ERROR_CODE_WEMEET_TEMPLATE_NOT_GET_TEMPLATE_INFO), err
	}

	BuildGetWholeBrandInfoRsp(meetingInfo, templateInfo, rsp)

	return 0, nil
}

// BuildGetWholeBrandInfoRsp
func BuildGetWholeBrandInfoRsp(meetingInfo *cachePb.MeetingInfo, templateInfo *pb.TemplateInfo,
	rsp *pb.GetWholeBrandInfoRsp) {
	rsp.Int32ErrorCode = proto.Int32(0)
	rsp.MsgMeetingInfo = meetingInfo
	rsp.MsgTemplateInfo = templateInfo
}

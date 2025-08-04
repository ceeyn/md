package service

import (
	"context"
	"errors"
	"strconv"

	"meeting_template/rpc"
	"meeting_template/util"

	"git.code.oa.com/trpc-go/trpc-go/log"
	errpb "git.code.oa.com/trpcprotocol/wemeet/common_xcast_meeting_error_code"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/golang/protobuf/proto"
)

// UpdateMeetingShareItem ...
func (s *WemeetMeetingTemplateOidbServiceImpl) UpdateMeetingShareItem(ctx context.Context,
	req *pb.UpdateMeetingShareItemReq, rsp *pb.UpdateMeetingShareItemRsp) error {
	bgTime := util.NowMs()

	util.ReportOne(util.UpdateMeetingShareItemReq) //[UpdateMeetingShareItem]请求
	if req.GetFromUserInfo() == nil || req.GetMeetingId() == 0 {
		return errors.New("meetingId or user info is empty")
	}
	isCreatorOrSuperAdmin := checkCreatorOrSuperAdmin(ctx, strconv.Itoa(int(req.GetFromUserInfo().GetAppId())),
		req.GetFromUserInfo().GetAppUserId(), strconv.FormatUint(req.GetMeetingId(), 10))
	if !isCreatorOrSuperAdmin {
		return errors.New("permission denied，current login user is not creator or admin")
	}
	rst, err := HandleUpdateMeetingShareItem(ctx, req)
	rsp.ErrorCode = proto.Int32(rst)
	rsp.MeetingId = proto.Uint64(req.GetMeetingId())
	rsp.ShareItemId = proto.String(req.GetShareItemId())
	rsp.ShowPassword = proto.Bool(req.GetShowPassword())
	cost := util.NowMs() - bgTime
	if err != nil {
		util.ReportOne(util.UpdateMeetingShareItemFail) //[UpdateMeetingShareItem]请求失败
		rsp.ErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "UpdateMeetingShareItem fail req:%v, rsp:%v, cost:%vms", req, rsp, cost)
	} else {
		rsp.ErrorMessage = proto.String("ok")
		log.InfoContextf(ctx, "UpdateMeetingShareItem ok  req:%v, rsp:%v, cost:%vms", req, rsp, cost)
	}
	return nil
}

// HandleUpdateMeetingShareItem ...
func HandleUpdateMeetingShareItem(ctx context.Context,
	req *pb.UpdateMeetingShareItemReq) (ret int32, err error) {

	// 参数校验
	rst, err := IsValidUpdateShareItemReqBody(req)
	if err != nil {
		util.ReportOne(util.UpdateMeetingShareItemInValidReq) //[UpdateMeetingShareItem]参数出错
		log.ErrorContextf(ctx, "isValidReqBody failed,  err:[%v]", err)
		return rst, err
	}

	meetingId := req.GetMeetingId()
	// 获取会议信息
	meetingInfo, isOversea, cas, err := rpc.GetMeetingInfo(ctx, meetingId)
	if err != nil {
		util.ReportOne(util.UpdateMeetingShareItemGetMeetingInfoFail) //[UpdateMeetingShareItem]获取会议信息失败
		log.ErrorContextf(ctx, "get meeting info from grocery fail,meetingId:%v, err: %v", meetingId, err)
		return int32(errpb.ERROR_CODE_COMM_ERROR_CODE_COMM_MEETING_NOT_EXIST), err
	}

	log.InfoContextf(ctx, "GetMeetingInfo succ,meetingInfo:%v, cas: %v", meetingInfo, cas)

	// //创建者判断
	// if meetingInfo.GetStrCreatorAppUid() != req.FromUserInfo.GetAppUserId() {
	// 	util.ReportOne(util.UpdateMeetingShareItemAuthFail) //[UpdateMeetingShareItem]鉴权失败
	// 	log.ErrorContextf(ctx, "Not Creator Update,meetingId:[%v], CreatorAppUid:[%v]",
	// 		meetingId, meetingInfo.GetStrCreatorAppUid())
	// 	return int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_INVALID_PARA), err
	// }

	// 设置StrShareItemId
	shareItemId := req.GetShareItemId()
	meetingInfo.StrShareItemId = proto.String(shareItemId)
	err = rpc.SetMeetingInfo(ctx, meetingInfo, isOversea, cas)
	if err != nil {
		//重试一次
		meetingInfo, isOversea, cas, err = rpc.GetMeetingInfo(ctx, meetingId)
		if err != nil {
			util.ReportOne(util.UpdateMeetingShareItemGetMeetingInfoFail) //[UpdateMeetingShareItem]获取会议信息失败
			log.ErrorContextf(ctx, "get meeting info from grocery fail,meetingId:%v, err: %v", meetingId, err)
			return int32(errpb.ERROR_CODE_COMM_ERROR_CODE_COMM_MEETING_NOT_EXIST), err
		} else {
			meetingInfo.StrShareItemId = proto.String(shareItemId)
			err = rpc.SetMeetingInfo(ctx, meetingInfo, isOversea, cas)
			if err != nil {
				util.ReportOne(util.UpdateMeetingShareItemSetMeetingInfoFail) //[UpdateMeetingShareItem]写会议信息失败
				log.ErrorContextf(ctx, "SetMeetingInfo failed, meetingId:[%v], shareItemId:[%v]",
					meetingId, shareItemId)
				return int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_SYSTEM_INNER_ERROR), err
			}
		}
	}
	// 设置showPassword
	showPassword := req.GetShowPassword()
	err = rpc.RDHSetShowPassword(ctx, isOversea, req.GetMeetingId(), meetingInfo.GetUint32OrderEndTime(), showPassword)
	if err != nil {
		return int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_SYSTEM_INNER_ERROR), err
	}

	return 0, nil
}

// IsValidUpdateShareItemReqBody ...
func IsValidUpdateShareItemReqBody(req *pb.UpdateMeetingShareItemReq) (code int32, err error) {
	if req.FromUserInfo == nil || req.FromUserInfo.AppUserId == nil {
		code = int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_INVALID_PARA)
		return code, errors.New("Invalid Parameter ")
	}
	if req.GetMeetingId() == 0 || req.FromUserInfo.GetAppUserId() == "" || req.GetShareItemId() == "" {
		code = int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_INVALID_PARA)
		return code, errors.New("Parameter Empty ")
	}
	return 0, nil
}

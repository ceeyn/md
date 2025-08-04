package service

import (
	"context"
	"errors"
	common "git.code.oa.com/trpcprotocol/wemeet/common_xcast_im_conversion_common"
	errpb "git.code.oa.com/trpcprotocol/wemeet/common_xcast_meeting_error_code"
	"strconv"
	"time"

	"meeting_template/rpc"
	"meeting_template/util"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/golang/protobuf/proto"
	wlcm "meeting_template/material_control/welcome"
)

// UpdateWelcomeInfo ...
func (s *WemeetMeetingTemplateOidbServiceImpl) UpdateWelcomeInfo(ctx context.Context,
	req *pb.UpdateWelcomeInfoReq, rsp *pb.UpdateWelcomeInfoRsp) error {

	bgTime := util.NowMs()

	if req.GetUint32AppId() == 0 || req.GetStrAppUid() == "" || req.GetMeetingId() == 0 {
		return errors.New("meetingId or user info is empty")
	}

	// 权鉴校验
	isCreatorOrSuperAdmin := checkCreatorOrSuperAdmin(ctx, strconv.Itoa(int(req.GetUint32AppId())), req.GetStrAppUid(), strconv.FormatUint(req.GetMeetingId(), 10))
	if !isCreatorOrSuperAdmin {
		return errors.New("permission denied，current login user is not creator or admin")
	}

	// 企业校验（只允许企业版和教育版）
	corpInfo, err := rpc.QueryCorpInfo(ctx, req.GetUint32AppId())
	if err != nil {
		metrics.IncrCounter("QueryCorpInfo.Err", 1) //[QueryCorpInfo]获取会议信息失败
		log.ErrorContextf(ctx, "Get CorpInfo Fail, req:%+v, err: %+v", req, err)
		return err
	}
	if corpInfo.GetTagInfo().GetTagType() != wlcm.ENTERPRISE_EDITION && corpInfo.GetTagInfo().GetTagType() != wlcm.EDUCATE_EDITION {
		log.ErrorContextf(ctx, "Corp Not Allow Update WelcomeInfo, req:%+v, tag_type: %+v", req, corpInfo.GetTagInfo().GetTagType())
		return errors.New("this operation is not possible with this version")
	}

	rst, err := HandleUpdateWelcomeInfo(ctx, req, rsp)
	rsp.ErrorCode = proto.Int32(rst)
	cost := util.NowMs() - bgTime
	if err != nil {
		metrics.IncrCounter("UpdateWelcomeInfo.Failed", 1) //[UpdateWelcomeInfo]请求失败
		rsp.ErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "UpdateWelcomeInfo fail, req:%+v, rsp:%+v, cost:%+vms",
			req, rsp, cost)
	} else {
		rsp.ErrorMessage = proto.String("ok")
		log.InfoContextf(ctx, "UpdateWelcomeInfo ok, req:%+v, rsp:%+v, cost:%+vms", req, rsp, cost)
	}
	metrics.IncrCounter("UpdateWelcomeInfo.Success", 1)
	return nil
}

// HandleUpdateWelcomeInfo ...
func HandleUpdateWelcomeInfo(ctx context.Context, req *pb.UpdateWelcomeInfoReq,
	rsp *pb.UpdateWelcomeInfoRsp) (ret int32, err error) {

	timestamp := uint32(time.Now().Unix())
	rsp.Timestamp = proto.Uint32(timestamp)
	rsp.SafeResult = &pb.WelComeSafeResult{
		TitleSafeCode:   proto.Int32(wlcm.WELCOME_SAFE_PASS),
		ContentSafeCode: proto.Int32(wlcm.WELCOME_SAFE_PASS),
	}

	//获取meeting_info
	meetingInfo, _, _, err := rpc.GetMeetingInfo(ctx, req.GetMeetingId())
	if err != nil {
		util.ReportOne(util.HandleGetWholeBrandInfoGetMeetingInfoFailAttr) //[HandleGetWholeBrandInfo]获取会议信息失败
		log.ErrorContextf(ctx, "get meeting info from grocery fail, meetingId:%+v, err: %+v", req.GetMeetingId(), err)
		return int32(errpb.ERROR_CODE_WEMEET_TEMPLATE_NOT_GET_MEETING_INFO), err
	}

	// 设置过期时间为会议预定结束时间
	expireTime := meetingInfo.GetUint32OrderEndTime() - timestamp
	if expireTime > 0 {
		expireTime = expireTime + wlcm.WelComeDefaultExpire
	} else {
		expireTime = wlcm.WelComeDefaultExpire
	}

	// 判断会议状态状态，不允许修改
	switch meetingInfo.GetUint32State() {
	case uint32(common.MEETING_STATE_MEETING_STATE_STARTED):
		return -1, errors.New(wlcm.MEETING_STATE_VALID_ERR_TIPS)
	case uint32(common.MEETING_STATE_MEETING_STATE_CANCELLED), uint32(common.MEETING_STATE_MEETING_STATE_ENDED), uint32(common.MEETING_STATE_MEETING_STATE_RECYCLED):
		return -1, errors.New(wlcm.MEETING_STATE_INVALID_ERR_TIPS)
	default:
		break
	}

	needSaveCache := false
	welCome := wlcm.NewWelComeImp()
	welcomeData := &pb.WelComeCache{
		AppUid:     proto.String(req.GetStrAppUid()),
		UpdateTime: proto.Uint32(timestamp),
	}
	// 如果关闭开关或者设置默认语，直接更新redis，无需送审
	if !wlcm.IsSwitchOpen(req.GetSwitchState()) || (wlcm.IsSwitchOpen(req.GetSwitchState()) && wlcm.IsDefaultWelcome(req.GetWelcomeData().GetWelcomeTitle(), req.GetWelcomeData().GetWelcomeContent())) {
		// 如果开关关闭
		if !wlcm.IsSwitchOpen(req.GetSwitchState()) {
			welcomeData.SwitchState = proto.Uint32(wlcm.WELCOME_SWITCH_OFF)
		} else {
			welcomeData.SwitchState = proto.Uint32(wlcm.WELCOME_SWITCH_ON)
			welcomeData.DefaultState = proto.Uint32(wlcm.WELCOME_DEFAULT_STATE_ON)
		}
		needSaveCache = true
	} else {
		containsUrl := false
		if util.ContainsURL(req.GetWelcomeData().GetWelcomeTitle()) {
			rsp.SafeResult.TitleSafeCode = proto.Int32(wlcm.WELCOME_SAFE_FAIL)
			rsp.SafeResult.TitleSafeDesc = proto.String(wlcm.UPDATE_FAIL_DEFAULT_TIPS)
			containsUrl = true
		}
		if util.ContainsURL(req.GetWelcomeData().GetWelcomeContent()) {
			rsp.SafeResult.ContentSafeCode = proto.Int32(wlcm.WELCOME_SAFE_FAIL)
			rsp.SafeResult.ContentSafeDesc = proto.String(wlcm.UPDATE_FAIL_DEFAULT_TIPS)
			containsUrl = true
		}

		if containsUrl {
			log.InfoContextf(ctx, "title/content contains url, meetingId:%+v, title:%+v, content:%+v", req.GetMeetingId(), req.GetWelcomeData().GetWelcomeTitle(), req.GetWelcomeData().GetWelcomeContent())
			return 0, nil
		}

		welcomeData.SwitchState = proto.Uint32(wlcm.WELCOME_SWITCH_ON)
		welcomeData.DefaultState = proto.Uint32(wlcm.WELCOME_DEFAULT_STATE_OFF)
		welcomeData.WelcomeTitle = proto.String(req.GetWelcomeData().GetWelcomeTitle())
		welcomeData.WelcomeContent = proto.String(req.GetWelcomeData().GetWelcomeContent())

		checkWords := make(map[string]string)
		checkWords[wlcm.WELCOME_TITLE_NAME] = req.GetWelcomeData().GetWelcomeTitle()
		checkWords[wlcm.WELCOME_CONTENT_NAME] = req.GetWelcomeData().GetWelcomeContent()
		checkResult, err := rpc.BatchCheckHasSensitiveWords(ctx, req.GetMeetingId(), req.GetUint32AppId(), req.GetStrAppUid(), req.GetStrUserIp(), checkWords)
		if err != nil {
			metrics.IncrCounter("BatchCheckHasSensitiveWords.Err", 1)
			// 调用失败，默认送审成功
			needSaveCache = true
		} else {
			if tips, ok := checkResult[wlcm.WELCOME_TITLE_NAME]; ok {
				rsp.SafeResult.TitleSafeCode = proto.Int32(wlcm.WELCOME_SAFE_FAIL)
				rsp.SafeResult.TitleSafeDesc = proto.String(tips)
			}
			if tips, ok := checkResult[wlcm.WELCOME_CONTENT_NAME]; ok {
				rsp.SafeResult.ContentSafeCode = proto.Int32(wlcm.WELCOME_SAFE_FAIL)
				rsp.SafeResult.ContentSafeDesc = proto.String(tips)
			}
			// 如果title和content没有同时审批通过，则直接返回，信息不存储
			if rsp.SafeResult.GetTitleSafeCode() == wlcm.WELCOME_SAFE_PASS && rsp.SafeResult.GetContentSafeCode() == wlcm.WELCOME_SAFE_PASS {
				needSaveCache = true
			}
		}
	}

	log.InfoContextf(ctx, "UpdateWelcomeInfo sensitive, need_save_cache:%+v, safe_result:%+v", needSaveCache, rsp.SafeResult)
	if needSaveCache {
		err = welCome.SetMeetingWelComeInfo(ctx, req.GetMeetingId(), welcomeData, expireTime)
		if err != nil {
			metrics.IncrCounter("SetMeetingWelComeInfo.Err", 1)
			rsp.ErrorMessage = proto.String(err.Error())
			return -3, err // redis读取失败，返回失败
		}
	}

	return 0, nil
}

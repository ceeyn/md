package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"meeting_template/config/config_rainbow"
	"meeting_template/es"
	"meeting_template/model"
	"meeting_template/rpc"
	"meeting_template/util"

	"git.code.oa.com/going/attr"
	"git.code.oa.com/meettrpc/meet_util"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	cachePb "git.code.oa.com/trpcprotocol/wemeet/common_meeting_cache"
	errpb "git.code.oa.com/trpcprotocol/wemeet/common_xcast_meeting_error_code"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	uuid "git.code.oa.com/wesee_ugc/go.uuid"
	"github.com/golang/protobuf/proto"
	t_extractor "meeting_template/util/t-helper/t-html/t-extractor"
)

const (
	IMGAction   = "get_webinar_image"
	IMGScene    = "sc_webinar_image"
	VideoAction = "get_webinar_video"
	VideoScene  = "sc_webinar_video"

	WarmUpShowVideoAction      = "get_webinar_warm_video_check"
	WarmUpShowVideoScene       = "sc_webinar_video"
	WarmUpShowVideoCallBackUri = "/wemeet-template/VideoCensorCallback"

	GateWayImgScenes            = "sc_webinar_warm_img"        // 暖场图片送审场景
	GateWayImgScenesAction      = "get_webinar_warm_img_check" // 暖场图片送审场景 Action
	GateWayImgScenesCallBackUri = "/wemeet-template/ImgCensorCallback"

	WebinarTemplateSponsor = "sc_webinar_template_sponsor" // webinar会议-主办方场景 送审场景
	WebinarTemplateDesc    = "sc_webinar_template_desc"    // webinar会议-会议介绍 送审场景

	CallbackTarget    = "polaris://trpc.wemeet.wemeet_template.http" //回调地址的target
	CallbackProtoType = "http"                                       //回调地址的配型
)

// UpdateTemplate ...
func (s *WemeetMeetingTemplateOidbServiceImpl) UpdateTemplate(ctx context.Context,
	req *pb.UpdateTemplateReq, rsp *pb.UpdateTemplateRsp) error {
	bgTime := util.NowMs()
	util.ReportOne(util.UpdateTemplateReq) //[GetTemplateInfo]请求

	if req.GetUint32AppId() == 0 || req.GetStrAppUid() == "" || req.GetMeetingId() == 0 {
		return errors.New("meetingId or user info is empty")
	}
	isCreatorOrSuperAdmin := checkCreatorOrSuperAdmin(ctx, strconv.Itoa(int(req.GetUint32AppId())),
		req.GetStrAppUid(), strconv.FormatUint(req.GetMeetingId(), 10))
	if !isCreatorOrSuperAdmin {
		return errors.New("permission denied，current login user is not creator or admin")
	}
	rst, err := HandleUpdateTemplate(ctx, req, rsp)
	rsp.ErrorCode = proto.Int32(rst)
	rsp.StrNonceId = proto.String(req.GetStrNonceId())
	if rsp.GetStrNonceId() == "" {
		rsp.StrNonceId = proto.String(uuid.NewV4().String())
	}
	cost := util.NowMs() - bgTime
	if err != nil {
		util.ReportOne(util.UpdateTemplateFail) //[UpdateTemplate]请求失败
		rsp.ErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "UpdateTemplate fail req:%v, rsp:%v, cost:%vms",
			req, rsp, cost)
	} else {
		rsp.ErrorMessage = proto.String("ok")
		log.InfoContextf(ctx, "UpdateTemplate ok  req:%v, rsp:%v, cost:%vms", req, rsp, cost)
	}
	return nil
}

// HandleUpdateTemplate
// NOCA:golint/fnsize(设计如此)
func HandleUpdateTemplate(ctx context.Context, req *pb.UpdateTemplateReq,
	rsp *pb.UpdateTemplateRsp) (ret int32, err error) {
	meetingId := req.GetMeetingId()
	appId := req.GetUint32AppId()
	if req.GetMeetingId() != 0 {
		_, err = HandleUpdateInviteSwitch(ctx, req)
		if err != nil {
			attr.AttrAPI(35765634, 1)
			log.InfoContextf(ctx, "do HandleUpdateInviteSwitch logic failed, req:%+v, err:%+v", req, err)
		}
	}
	ret, err = HandleUserSafeCheck(ctx, req) // 先进行用户风控检查
	if ret != 0 {
		attr.AttrAPI(35907557, 1)
		log.InfoContextf(ctx, "HandleUpdateTemplate check user safe need block. req:%+v", req)
		return ret, err
	}
	if !CheckUserHasWarmUpPower(ctx, req) { //接入暖场相关权限点
		metrics.IncrCounter("HandleUpdateTemplate_NO_Warm_Power", 1)
		log.ErrorContextf(ctx, "HandleUpdateTemplate check user dont have warmUp power. req:%+v", req)
		return -6, errors.New("dont have warmUp power")
	}
	templateInfo := &model.TemplateInfo{}
	FillUpdataTemplateDirectInfo(req, meetingId, appId, templateInfo)
	//输入参数校验
	ret, err = IsUpdateRequestValid(ctx, req, templateInfo)
	if ret != 0 {
		return ret, err
	}
	warmUpData := req.GetWarmUpData()
	err = util.CheckWarmUpDataFormat(ctx, warmUpData)
	if err != nil {
		util.ReportOne(util.CoverUrlCosIdInvalid)
		log.ErrorContextf(ctx, "WarmUpData cosId invalid, templateId:%v", req.GetTemplateId())
		return -4, err
	}
	oldTemplateInfo := &model.TemplateInfo{}
	oldTemplateInfo, err = GetTemplateInfoSingleFlight(ctx, req.GetTemplateId())
	//先上报，不报错返回，考虑预订相同配置的会议
	if oldTemplateInfo.MeetingId != "" && oldTemplateInfo.MeetingId != strconv.FormatUint(meetingId, 10) {
		metrics.IncrCounter("warn_template_HandleUpdateTemplate_meeting_not_match", 1)
		util.ReportOne(36354620) //[HandleUpdateTemplate]模板会议不匹配
		log.ErrorContextf(ctx, "Template meeting NOT match. oldTemplateInfo.MeetingId:%v, "+
			"meetingId:%v", oldTemplateInfo.MeetingId, meetingId)
	}

	needCheckImg, needDelayQueryM3U8Status := false, false
	WarmUpInfoPreprocess(ctx, req, warmUpData, oldTemplateInfo, &needCheckImg, &needDelayQueryM3U8Status) // warmUpData预处理，指针里面会修改值
	warmUpDataStr, err := util.GetSerializedJsonStr(ctx, warmUpData)
	if err != nil {
		util.ReportOne(util.JsonMarshalFail)
		log.InfoContextf(ctx, "get warmUpDataStr failed, warmUpData::%v, err:%v", warmUpData, err)
	}
	needPushWarmUpUpdate := false
	if warmUpData != nil && 0 < len(warmUpData.WarmupImgList) && 0 < len(warmUpData.WarmupVideoList) {
		bgTime := util.NowMs()
		getRedisCost := util.NowMs() - bgTime
		if err != nil {
			rsp.ErrorMessage = proto.String(err.Error())
			log.ErrorContextf(ctx, "judge warmUpData, redis GetTemplateInfo failed,req:%v,", req)
		} else {
			needPushWarmUpUpdate = (warmUpDataStr != oldTemplateInfo.WarmUpData)
			newStatusWarmUpStr := ProcessWarmUpDataAudit(ctx, req, warmUpData, oldTemplateInfo, &needPushWarmUpUpdate)
			if newStatusWarmUpStr != "" {
				warmUpDataStr = newStatusWarmUpStr
			}
		}
		log.InfoContextf(ctx, "judge warmUpData,Id:[%v],getRedisCost:%vms", req.GetTemplateId(), getRedisCost)
	}
	judgeOtherCaseNeedPush(ctx, req, oldTemplateInfo, &needPushWarmUpUpdate)
	//这里会议介绍内容写es
	if CanDoEsLogic() {
		go func(newCtx context.Context) {
			defer meet_util.DefPanicFun()
			HandleMeetingDescriptionESLogic(newCtx, req, oldTemplateInfo)
		}(trpc.CloneContext(ctx))
	}
	templateInfo.WarmUpData = warmUpDataStr
	templateInfo.TemplateId = req.GetTemplateId()
	err = SetTemplateInfo(ctx, templateInfo)
	if err != nil {
		rsp.ErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "UpdateTemplate fail req:%v, rsp:%v", req, rsp)
		return -2, err
	}
	if needPushWarmUpUpdate {
		defer UserNotifyWarmUpDataUpdate(ctx, meetingId, appId, req.GetStrAppUid(), templateInfo)
	}
	if needCheckImg && len(req.GetWarmUpData().GetWarmupImgList()) == 1 { // 图片送审
		AuditImage(ctx, req)
	}
	if needDelayQueryM3U8Status {
		metrics.IncrCounter("UpdateTemplate.SendDelayJob.Sum", 1)
		log.InfoContextf(ctx, "HandleUpdateTemplate need delay query m3u8 trans status, req:%+v", req)
		go func(newCtx context.Context) {
			defer meet_util.DefPanicFun()
			HandleWarmUpVideoDelayJob(newCtx, req)
		}(trpc.CloneContext(ctx))
	}
	timestamp := uint32(time.Now().Unix())
	nonce := int32(timestamp) // 随机暂定跟时间戳相关
	rsp.Nonce = proto.Int32(nonce)
	rsp.Timestamp = proto.Uint32(timestamp)
	return 0, nil
}

// UserNotifyWarmUpDataUpdate
func UserNotifyWarmUpDataUpdate(ctx context.Context, meetingId uint64, appId uint32, appUid string,
	templateInfo *model.TemplateInfo) {
	var coverListNeedTranCosIds []string
	GetOtherCosIdsFromTemplateInfo(ctx, templateInfo, &coverListNeedTranCosIds)
	var allNeedTransCosIds []string
	allNeedTransCosIds = append(allNeedTransCosIds, coverListNeedTranCosIds...)
	var warmUpImgCosIds []string
	var warmUpVideoCosIds []string
	var warmUpVideoStreamCosIds []string
	warmUpData := &pb.WarmUpData{}
	if GetWarmUpDataFromTemplateInfo(ctx, templateInfo, warmUpData) {
		GetPicCosIdsFromWarmUpData(ctx, warmUpData, &warmUpImgCosIds)
		GetVideoCosIdsFromWarmUpData(ctx, warmUpData, &warmUpVideoCosIds, &warmUpVideoStreamCosIds)
	}
	cosIdUrlMap := make(map[string]string)
	GetAndStoreDownloadUrls(ctx, util.RawImageType, allNeedTransCosIds, warmUpImgCosIds, warmUpVideoCosIds,
		warmUpVideoStreamCosIds, cosIdUrlMap)
	outWarmUpData := &pb.WarmUpData{}
	if ConvertWarmUpData(ctx, templateInfo, cosIdUrlMap, outWarmUpData, false) {
		ReplaceWarmUpImgWithCover(ctx, templateInfo, cosIdUrlMap, outWarmUpData)
	}
	rpc.NotifyWarmUpDataUpdate(ctx, meetingId, appId, appUid, outWarmUpData)
}

// 处理暖场视频审核和相关状态更
func ProcessWarmUpDataAudit(ctx context.Context, req *pb.UpdateTemplateReq, warmUpData *pb.WarmUpData,
	oldTemplateInfo *model.TemplateInfo, needPushWarmUpUpdate *bool) string {
	// 使用图片获取无视频信息时不用处理; 暖场视频暂只支持一个视频
	if warmUpData.GetUint32WarmupUseType() == util.WarmUpShowImg || len(warmUpData.WarmupVideoList) != 1 {
		log.InfoContextf(ctx, "ProcessWarmUpDataAudit,use pic or no video, no need process,type:%v, lenVideo:%v",
			warmUpData.GetUint32WarmupUseType(), len(warmUpData.WarmupImgList))
		return ""
	}
	videoCosId := warmUpData.WarmupVideoList[0].GetCosId()
	buf := []byte(oldTemplateInfo.WarmUpData)
	oldWarmUpData := pb.WarmUpData{}
	err := json.Unmarshal(buf, &oldWarmUpData)
	oldStatus := uint32(util.VideoNotAudit)
	if err != nil {
		util.ReportOne(util.JsonUnmarshalFail)
		oldStatus = util.VideoNotAudit
		log.ErrorContextf(ctx, "parse warmup data fail, data:%+v,error: %+v", oldTemplateInfo.WarmUpData, err)
	}
	if len(oldWarmUpData.WarmupVideoList) != 1 {
		oldStatus = util.VideoNotAudit
		log.InfoContextf(ctx, "ProcessWarmUpDataAudit, old warmup data size not 1: %+v", oldWarmUpData)
	} else {
		oldCosId := oldWarmUpData.WarmupVideoList[0].GetCosId()
		if oldCosId == videoCosId {
			oldStatus = oldWarmUpData.WarmupVideoList[0].GetStatus()
		}
	}
	decodedCosId := string(util.GetBase64Decoded(ctx, videoCosId))
	// uuid由template_id + "&" + cos_id构成
	auditUuid := req.GetTemplateId() + "&" + decodedCosId
	newStatus := warmUpData.WarmupVideoList[0].GetStatus()
	log.InfoContextf(ctx, "ProcessWarmUpDataAudit,oldData:%+v,newData:%v,", oldWarmUpData, warmUpData)
	if newStatus == util.VideoPassAudit || newStatus == util.VideoAuditFail {
		if oldStatus != newStatus {
			*needPushWarmUpUpdate = true
		}
		return ""
	}
	if oldStatus != util.VideoPassAudit && oldStatus != util.VideoAuditFail {
		log.DebugContextf(ctx, "ProcessWarmUpDataAudit, needPushWarmUpUpdate")
		// 视频未审核时，说明为新上传，需要push暖场物料更新
		*needPushWarmUpUpdate = true
		url := GetDownloadUrl(ctx, util.NotPicImageType, decodedCosId, util.DownloadUseCdn)
		callBackSwitch := config_rainbow.GetCallBackConf().WarmUpShowVideoAction
		log.InfoContextf(ctx, "ctlrpcImgSafetyAudit,WarmUpShowVideoAction: %v", callBackSwitch)
		if callBackSwitch {
			warmVideoReq := &rpc.SafetyAuditReq{
				Uuid:              auditUuid,
				AppId:             req.GetUint32AppId(),
				AppUid:            req.GetStrAppUid(),
				MeetingId:         req.GetMeetingId(),
				Url:               url,
				Scenes:            WarmUpShowVideoScene,
				Action:            WarmUpShowVideoAction,
				CallbackTarget:    CallbackTarget,
				CallbackUri:       WarmUpShowVideoCallBackUri,
				CallbackEnvName:   util.GetEnvName(),
				CallbackNameSpace: util.GetNameSpace(),
				CallbackProtoType: CallbackProtoType,
			}
			rpc.WarmVideoSafetyAudit(ctx, warmVideoReq)
		} else {
			rpc.VideoSafetyAudit(ctx, auditUuid, req.GetUint32AppId(), req.GetStrAppUid(),
				req.GetMeetingId(), url)
		}
		return ""
	}
	warmUpData.WarmupVideoList[0].Status = proto.Uint32(oldStatus)
	warmUpDataStr, err := util.GetSerializedJsonStr(ctx, warmUpData)
	if err != nil {
		util.ReportOne(util.JsonMarshalFail)
		log.InfoContextf(ctx, "get warmUpDataStr fail, warmUpData:%v, err:%v", warmUpData, err)
		return ""
	}
	return warmUpDataStr
}

// FillUpdataTemplateDirectInfo
func FillUpdataTemplateDirectInfo(req *pb.UpdateTemplateReq, meetingId uint64, appId uint32,
	templateInfo *model.TemplateInfo) {
	templateInfo.Sponsor = req.GetSponsor()
	templateInfo.CoverName = req.GetCoverName()
	templateInfo.CoverUrl = req.GetCoverUrl()
	templateInfo.Description = req.GetDescription()
	templateInfo.AppId = strconv.FormatUint(uint64(appId), 10)
	templateInfo.AppUid = req.GetStrAppUid()
	templateInfo.MeetingId = strconv.FormatUint(meetingId, 10)
}

// IsUpdateRequestValid 判断请求是否合法，ret=0 表示正常，非0异常
func IsUpdateRequestValid(ctx context.Context, req *pb.UpdateTemplateReq,
	templateInfo *model.TemplateInfo) (ret int32, err error) {

	templateId := req.GetTemplateId()
	appId := req.GetUint32AppId()
	appUId := req.GetStrAppUid()
	meetingId := req.GetMeetingId()
	if templateId == "" {
		metrics.IncrCounter("warn_template_IsUpdateRequestValid_template_id_empty", 1)
		util.ReportOne(36354616) //[IsUpdateRequestValid]模板id为空
		log.ErrorContextf(ctx, "update template_id empty, no need set redis")
		return -6, err
	}
	if meetingId == 0 {
		metrics.IncrCounter("warn_template_IsUpdateRequestValid_meeting_id_empty", 1)
		util.ReportOne(36354617) //[IsUpdateRequestValid]输入会议id为空
		log.ErrorContextf(ctx, "req meetingId empty")
		return int32(errpb.ERROR_CODE_WEMEET_TEMPLATE_TEMPLATE_PARAMS_INVALID), err
	}
	//优先校验会议id，问题是不同场次的会议的template_id相同。
	//怀疑是更新模板的时候传错了template_id，拿另外一个场次的template_id传过来了
	meetingInfo, _, _, err := rpc.GetMeetingInfo(ctx, meetingId)
	if err != nil {
		metrics.IncrCounter("fail_template_IsUpdateRequestValid_get_info", 1)
		util.ReportOne(36354618) //[IsUpdateRequestValid]获取会议信息失败
		log.ErrorContextf(ctx, "GetMeetingInfo fail, meetingId:[%v], err: %v",
			meetingId, err)
		return int32(errpb.ERROR_CODE_WEMEET_TEMPLATE_NOT_GET_MEETING_INFO), err
	}
	if meetingInfo.GetStrMeetingTemplateId() != templateId {
		metrics.IncrCounter("warn_template_IsUpdateRequestValid_template_not_match", 1)
		util.ReportOne(36354619) //[IsUpdateRequestValid]会议模板不匹配
		log.ErrorContextf(ctx, "Meeting template Not match, GetStrMeetingTemplateId:[%v], "+
			"templateId: %v", meetingInfo.GetStrMeetingTemplateId(), templateId)
		return int32(errpb.ERROR_CODE_WEMEET_TEMPLATE_MEETING_TEMPLATE_NOT_MATCH), err
	}

	coverItems := req.GetCoverList()
	err = util.CheckCoverListFormat(ctx, coverItems)
	if err != nil {
		util.ReportOne(util.CoverUrlCosIdInvalid)
		log.ErrorContextf(ctx, "CoverList cosId invalid, templateId:%v", templateId)
		return -4, err
	}
	coverListStr, err := util.GetSerializedJsonStr(ctx, coverItems)
	if err != nil {
		util.ReportOne(util.JsonMarshalFail)
		log.InfoContextf(ctx, "get coverListStr failed, coverItems:%v, err:%v", coverItems, err)
	}
	templateInfo.CoverList = coverListStr
	decodedCoverUrl := string(util.GetBase64Decoded(ctx, templateInfo.CoverUrl))
	cosIdValid := util.IsValidCosId(ctx, decodedCoverUrl)
	if !cosIdValid {
		util.ReportOne(util.CoverUrlCosIdInvalid)
		err = errors.New("cover url cos id invalid, cosId:" + templateInfo.CoverUrl)
		log.ErrorContextf(ctx, "CoverUrl hit cos id format filter, cosId:%v,", templateInfo.CoverUrl)
		return -4, err
	}
	hitSensitive, word := rpc.CheckLongTextSensitiveData(ctx, templateInfo.Sponsor, false,
		appId, appUId, WebinarTemplateSponsor)
	if hitSensitive {
		util.ReportOne(util.SponsorCoverNameHitSensitive)
		err = errors.New("sponsor or cover name has sensitive word:" + word)
		log.ErrorContextf(ctx, "hit sensitive,sponsor:%v,sensitive word:%v", templateInfo.Sponsor, word)
		return -1, err
	}
	descriptionHitSensitive, words := rpc.CheckLongTextSensitiveData(ctx, templateInfo.Description, true,
		appId, appUId, WebinarTemplateDesc)
	if descriptionHitSensitive {
		util.ReportOne(util.DescriptionHitSensitive)
		err = errors.New("description has sensitive words:" + words)
		log.ErrorContextf(ctx, "description hit sensitive,sensitive words:%v, templateId:%v", words, templateId)
		return -3, err
	}
	// 详情压缩后再存redis
	compressed, err := util.GetCompressed(ctx, templateInfo.Description)
	if err != nil {
		util.ReportOne(util.CompressDescriptionFail)
		log.ErrorContextf(ctx, "description compress error:%v, templateId:%v", err, templateId)
		return -5, err
	}
	templateInfo.Description = compressed
	return 0, nil
}

// HandleUpdateInviteSwitch.. 处理暖场邀请开关
func HandleUpdateInviteSwitch(ctx context.Context, req *pb.UpdateTemplateReq) (int32, error) {
	if req.GetInviteSwitchState() != util.InviteStateClose && req.GetInviteSwitchState() != util.InviteStateOpen {
		attr.AttrAPI(35765656, 1)
		log.ErrorContextf(ctx, "InviteSwitchState is not 0 or 1, req:%+v", req)
		return util.ErrSwitchStateIllegal, errors.New("SwitchState is illegal")
	}
	key := util.MakeWarmUpInviteSwitchKey(req.GetMeetingId())
	err := rpc.RDUpdateInviteSwitch(ctx, key, req.GetInviteSwitchState())
	if err != nil {
		attr.AttrAPI(35765657, 1)
		log.ErrorContextf(ctx, "rpc RDUpdateInviteSwitch failed. meetingId:%+v, err:%+v", req.GetMeetingId(), err)
		return util.ErrUpdateInviteSwitch, err
	}
	attr.AttrAPI(35765658, 1)
	log.InfoContextf(ctx, "HandleUpdateInviteSwitch succ. meetingId:%+v", req.GetMeetingId())
	return 0, nil
}

// judgeOtherCaseNeedPush ...
func judgeOtherCaseNeedPush(ctx context.Context, req *pb.UpdateTemplateReq, oldTemplateInfo *model.TemplateInfo,
	needPushWarmUpUpdate *bool) {
	if oldTemplateInfo.WarmUpData == "" {
		if req.GetWarmUpData().GetUint32PlayMusicType() != 0 {
			log.InfoContextf(ctx, "judgeOtherCaseNeedPush need push, because new choose music")
			*needPushWarmUpUpdate = true
		}
		return
	}
	oldWarmUpData := pb.WarmUpData{}
	buf := []byte(oldTemplateInfo.WarmUpData)
	err := json.Unmarshal(buf, &oldWarmUpData)
	if err != nil {
		log.ErrorContextf(ctx, "judgeOtherCaseNeedPush dont push, because warmup json Unmarshal failed")
		return
	}
	// 变更音乐类型
	if oldWarmUpData.GetUint32PlayMusicType() != req.GetWarmUpData().GetUint32PlayMusicType() {
		log.InfoContextf(ctx, "judgeOtherCaseNeedPush need push, because music type change")
		*needPushWarmUpUpdate = true
		return
	}
	// 切回图片类型，原来勾选了音乐
	if oldWarmUpData.GetUint32WarmupUseType() == 1 && oldWarmUpData.GetUint32PlayMusicType() != 0 &&
		req.GetWarmUpData().GetUint32WarmupUseType() == 0 {
		log.InfoContextf(ctx, "judgeOtherCaseNeedPush need push, change to img type and has music")
		*needPushWarmUpUpdate = true
		return
	}
	// 删除图片
	if req.GetWarmUpData().GetUint32WarmupUseType() == 0 && len(req.GetWarmUpData().GetWarmupImgList()) == 0 &&
		len(oldWarmUpData.GetWarmupImgList()) == 1 {
		log.InfoContextf(ctx, "judgeOtherCaseNeedPush need push, because img delete")
		*needPushWarmUpUpdate = true
		return
	}
	// 删除视频
	if req.GetWarmUpData().GetUint32WarmupUseType() == 1 && len(req.GetWarmUpData().GetWarmupVideoList()) == 0 &&
		len(oldWarmUpData.GetWarmupVideoList()) == 1 {
		log.InfoContextf(ctx, "judgeOtherCaseNeedPush need push, because video delete")
		*needPushWarmUpUpdate = true
		return
	}
	if oldWarmUpData.GetUint32WarmupUseType() != req.GetWarmUpData().GetUint32WarmupUseType() {
		log.InfoContextf(ctx, "judgeOtherCaseNeedPush need push, change user type")
		*needPushWarmUpUpdate = true
		return
	}
	return
}

// WarmUpInfoPreprocess ...
func WarmUpInfoPreprocess(ctx context.Context, req *pb.UpdateTemplateReq, warmUpData *pb.WarmUpData,
	oldTemplateInfo *model.TemplateInfo, needCheckImg *bool, needDelayQueryM3U8Status *bool) {
	buf := []byte(oldTemplateInfo.WarmUpData)
	oldWarmUpData := pb.WarmUpData{}
	err := json.Unmarshal(buf, &oldWarmUpData)
	if err != nil {
		attr.AttrAPI(35907888, 1)
		log.ErrorContextf(ctx, "parse warmup data fail, data:%+v, error: %+v", oldTemplateInfo.WarmUpData, err)
	}

	// 老的视频的审核状态
	if len(oldWarmUpData.GetWarmupVideoList()) == 1 && len(warmUpData.GetWarmupVideoList()) == 1 {
		oldVideoCosId := oldWarmUpData.GetWarmupVideoList()[0].GetCosId()
		newVideoCosId := warmUpData.GetWarmupVideoList()[0].GetCosId()
		videoStatus := uint32(util.VideoNotAudit)
		if oldVideoCosId == newVideoCosId {
			videoStatus = oldWarmUpData.GetWarmupVideoList()[0].GetStatus()
		}
		warmUpData.WarmupVideoList[0].Status = proto.Uint32(videoStatus)
	}

	//视频转码相关
	if len(warmUpData.GetWarmupVideoList()) == 1 {
		m3u8TransStatus := uint32(util.M3U8TransNoReady)
		updateTimeStamp := uint64(time.Now().Unix())
		if len(oldWarmUpData.GetWarmupVideoList()) == 1 {
			if oldWarmUpData.GetWarmupVideoList()[0].GetCosId() == warmUpData.GetWarmupVideoList()[0].GetCosId() {
				m3u8TransStatus = oldWarmUpData.GetWarmupVideoList()[0].GetVideoTransStatus()
				updateTimeStamp = oldWarmUpData.GetWarmupVideoList()[0].GetUpdateTimeStamp()
			} else { // 视频发生变化
				log.InfoContextf(ctx, "WarmUpInfoPreprocess video cosId changed, so needQueryM3U8Status, "+
					"req:%+v, cacheTemplateInfo:%+v", req, oldTemplateInfo)
				*needDelayQueryM3U8Status = true
			}
		} else if len(oldWarmUpData.GetWarmupVideoList()) == 0 { //新上传视频
			log.InfoContextf(ctx, "WarmUpInfoPreprocess find user new upload video, so needQueryM3U8Status, "+
				"req:%+v, cacheTemplateInfo:%+v", req, oldTemplateInfo)
			*needDelayQueryM3U8Status = true
		}
		warmUpData.WarmupVideoList[0].VideoTransStatus = proto.Uint32(m3u8TransStatus)
		warmUpData.WarmupVideoList[0].UpdateTimeStamp = proto.Uint64(updateTimeStamp)
	}

	// 图片相关
	if len(warmUpData.GetWarmupImgList()) != 1 {
		log.InfoContextf(ctx, "HandleWarmUpImgAudit imgList is not 1, req:%+v", req)
		return
	}
	status := uint32(util.NotAudit)
	imgCosId := warmUpData.WarmupImgList[0].GetUrl()
	// 之前上传过图片
	if len(oldWarmUpData.WarmupImgList) == 1 {
		oldCosId := oldWarmUpData.WarmupImgList[0].GetUrl()
		if oldCosId == imgCosId { // 是相同的图片
			status = oldWarmUpData.WarmupImgList[0].GetStatus()
		}
		// 物料类型是视频，而且上次传的图片审核未通过
		if warmUpData.GetUint32WarmupUseType() == util.WarmUpShowVideo && oldCosId == "" &&
			oldWarmUpData.GetWarmupImgList()[0].GetName() == "" &&
			oldWarmUpData.WarmupImgList[0].GetStatus() == util.AuditFail {
			attr.AttrAPI(35907570, 1)
			status = util.AuditFail
			warmUpData.WarmupImgList[0].Name = proto.String("")
			warmUpData.WarmupImgList[0].Url = proto.String("")
		}
	}
	// 未送审的时候进行送审
	if status == util.NotAudit {
		*needCheckImg = true
	}
	warmUpData.WarmupImgList[0].Status = proto.Uint32(status)
}

// HandleUserSafeCheck ...
func HandleUserSafeCheck(ctx context.Context, req *pb.UpdateTemplateReq) (int32, error) {
	if req.GetWarmUpData() == nil {
		return 0, nil
	}
	if len(req.GetWarmUpData().GetWarmupImgList()) == 0 && len(req.GetWarmUpData().GetWarmupVideoList()) == 0 {
		return 0, nil
	}
	log.InfoContextf(ctx, "UpdateTemplate HandleUserSafeCheck do GetSafe logic. req:%+v", req)
	action := IMGAction
	scenes := IMGScene
	if req.GetWarmUpData().GetUint32WarmupUseType() == util.WarmUpShowVideo {
		action = VideoAction
		scenes = VideoScene
	}
	rst, tips := rpc.GetSafe(ctx, action, req, "", scenes)
	if rst != 0 {
		attr.AttrAPI(35907576, 1)
		log.ErrorContextf(ctx, "UpdateTemplate HandleUserSafeCheck failed, req:%+v", req)
		return rst, errors.New(tips)
	}
	return 0, nil
}

// AuditImage ...
func AuditImage(ctx context.Context, req *pb.UpdateTemplateReq) {
	imgCosId := req.GetWarmUpData().GetWarmupImgList()[0].GetUrl()
	decodedCosId := string(util.GetBase64Decoded(ctx, imgCosId))
	auditUuid := req.GetTemplateId() + "&" + decodedCosId // uuid由template_id + "&" + cos_id构成
	url := GetDownloadUrl(ctx, util.RawImageType, decodedCosId, util.DownloadUseCdn)
	callBackSwitch := config_rainbow.GetCallBackConf().GateWayImgScenesAction
	log.InfoContextf(ctx, "ctlrpcImgSafetyAudit GateWayImgScenesAction: %v", callBackSwitch)
	if callBackSwitch {
		safeReq := &rpc.SafetyAuditReq{
			Uuid:              auditUuid,
			AppId:             req.GetUint32AppId(),
			AppUid:            req.GetStrAppUid(),
			MeetingId:         req.GetMeetingId(),
			Url:               url,
			Scenes:            GateWayImgScenes,
			Action:            GateWayImgScenesAction,
			CallbackTarget:    CallbackTarget,
			CallbackUri:       GateWayImgScenesCallBackUri,
			CallbackEnvName:   util.GetEnvName(),
			CallbackNameSpace: util.GetNameSpace(),
			CallbackProtoType: CallbackProtoType,
		}
		rpc.ImgSafetyAuditV2(ctx, safeReq)
	} else {
		rpc.ImgSafetyAudit(ctx, auditUuid, req.GetUint32AppId(), req.GetStrAppUid(), req.GetMeetingId(), url, GateWayImgScenes)
	}
}

// HandleMeetingDescriptionESLogic ...
func HandleMeetingDescriptionESLogic(ctx context.Context, req *pb.UpdateTemplateReq,
	oldTemplateInfo *model.TemplateInfo) error {
	description := req.GetDescription()
	decodedDescription := util.GetBase64Decoded(ctx, description)
	if len(decodedDescription) == 0 {
		log.InfoContextf(ctx, "HandleMeetingDescriptionESLogic get template Info:%+v", oldTemplateInfo)
		decompressed, err := util.GetDeCompressed(ctx, oldTemplateInfo.Description)
		if err != nil {
			attr.AttrAPI(36337036, 1)
			log.ErrorContextf(ctx, "HandleMeetingDescriptionESLogic DeCompressed failed. tplID:%+v, err:%+v",
				oldTemplateInfo.TemplateId, err)
		}
		log.InfoContextf(ctx, "HandleMeetingDescriptionESLogic get DeCompressed description:%+v, tplID:%+v",
			decompressed, oldTemplateInfo.TemplateId)
		if decompressed != "" { //原来是有内容的 删除es文档
			es.DelIntroductionToES(ctx, fmt.Sprint(req.GetMeetingId()))
		}
		return nil
	}
	extractor := &t_extractor.TSimpleExtractor{}
	elements, err := extractor.Parse(string(decodedDescription))
	if err != nil {
		log.ErrorContextf(ctx, "HandleMeetingDescriptionESLogic parse text failed, text: %+v, error: %+v",
			decodedDescription, err)
		return err
	}
	longTexts := elements.TextFrags
	log.InfoContextf(ctx, "HandleMeetingDescriptionESLogic meetingId:%+v, longTexts:%+v",
		req.GetMeetingId(), longTexts)
	text := strings.Join(longTexts, "")
	log.InfoContextf(ctx, "HandleMeetingDescriptionESLogic meetingId:%+v, get text:%+v",
		req.GetMeetingId(), text)
	//获取 meetingInfo
	meetingInfo, _, _, err := rpc.GetMeetingInfo(ctx, req.GetMeetingId())
	if err != nil {
		attr.AttrAPI(36337037, 1)
		log.ErrorContextf(ctx, "HandleSaveItinerary failed, meetingId:%+v, err:%+v", req.GetMeetingId(), err)
		return nil
	}
	info := FillESIntroductionInfo(ctx, meetingInfo, text)
	err = es.UpsertIntroductionToES(ctx, fmt.Sprint(meetingInfo.GetUint64MeetingId()), info)
	if err != nil {
		attr.AttrAPI(36337038, 1)
		log.ErrorContextf(ctx, "HandleMeetingDescriptionESLogic update ES failed, meetingId:%+v, err:%+v",
			req.GetMeetingId(), err)
	}
	return nil
}

// FillESIntroductionInfo ...
func FillESIntroductionInfo(ctx context.Context, meetingInfo *cachePb.MeetingInfo, text string) *es.Introduction {
	introduction := &es.Introduction{
		MeetingId:           fmt.Sprint(meetingInfo.GetUint64MeetingId()),
		AppId:               fmt.Sprint(meetingInfo.GetUint32CreatorSdkappid()),
		AppUid:              meetingInfo.GetStrCreatorAppUid(),
		MeetingIntroduction: text,
	}
	now := time.Now().Format("2006-01-02 15:04:05")

	introduction.CreateTime = now
	introduction.UpdateTime = now
	// 会议预定时间
	orderTm := time.Unix(int64(meetingInfo.GetUint32OrderTime()), 0)
	strOrderTime := orderTm.Format("2006-01-02 15:04:05")
	// 会议预定开始时间
	orderStartTm := time.Unix(int64(meetingInfo.GetUint32OrderStartTime()), 0)
	strOrderStartTime := orderStartTm.Format("2006-01-02 15:04:05")
	// 会议预定结束时间
	orderEndTm := time.Unix(int64(meetingInfo.GetUint32OrderEndTime()), 0)
	strOrderEndTime := orderEndTm.Format("2006-01-02 15:04:05")
	//会议主题
	meetingSubject, _ := base64.StdEncoding.DecodeString(string(meetingInfo.GetBytesMeetingSubject()))
	introduction.MeetingOrderTime = strOrderTime
	introduction.MeetingOrderStartTime = strOrderStartTime
	introduction.MeetingOrderEndTime = strOrderEndTime
	introduction.MeetingSubject = string(meetingSubject)
	return introduction
}

// CheckUserHasWarmUpPower ...
func CheckUserHasWarmUpPower(ctx context.Context, req *pb.UpdateTemplateReq) bool {
	warmPowerSwitch := config_rainbow.GetWarmPowerSwitchConfConfig()
	if warmPowerSwitch == "close" {
		return true
	}
	isNeedCheckPower := rpc.GetUserSettings(ctx, req.GetStrAppUid(), req.GetUint32AppId())
	if !isNeedCheckPower { //不在特定企业 不校验权限点
		return true
	}
	hasChangeWarm := false
	// 上传了图片
	if req.GetWarmUpData().GetUint32WarmupUseType() == util.WarmUpShowImg &&
		len(req.GetWarmUpData().GetWarmupImgList()) > 0 {
		hasChangeWarm = true
	}
	// 上传了视频
	if req.GetWarmUpData().GetUint32WarmupUseType() == util.WarmUpShowVideo &&
		len(req.GetWarmUpData().GetWarmupVideoList()) > 0 {
		hasChangeWarm = true
	}
	// 勾选了音乐
	if req.GetWarmUpData().GetUint32PlayMusicType() != 0 {
		hasChangeWarm = true
	}
	if hasChangeWarm {
		warmUpPower := "warm_up"
		hasPower := rpc.CheckUserHasPower(ctx, req.GetUint32AppId(), req.GetStrAppUid(), warmUpPower)
		if !hasPower {
			metrics.IncrCounter("CheckUserHasPower_No_Warm_Power", 1)
			log.ErrorContextf(ctx, "CheckUserHasWarmUpPower user dont have warmUP power, req:%+v", req)
			return false
		}
	}
	return true
}

// HandleWarmUpVideoDelayJob ...
func HandleWarmUpVideoDelayJob(ctx context.Context, req *pb.UpdateTemplateReq) error {
	cosId := ""
	if len(req.GetWarmUpData().GetWarmupVideoList()) == 1 {
		cosId = req.GetWarmUpData().GetWarmupVideoList()[0].GetCosId()
	} else { // cosId为空的情况
		metrics.IncrCounter("HandleWarmUpVideoDelayJob.VideoCosIdNot1", 1)
		log.InfoContextf(ctx, "HandleWarmUpVideoDelayJob but video cosId not 1, req:%+v", req)
		return nil
	}
	msg := &model.TemplateMqMsg{
		MeetingId:  req.GetMeetingId(),
		TemplateId: req.GetTemplateId(),
		CosId:      cosId,
		MsgType:    0,
		TryCount:   1,
	}
	byteMsg, _ := json.Marshal(msg)
	delaySpan := config_rainbow.GetMusicConfConfig().DelaySpan
	err := rpc.SendDelayJob(ctx, byteMsg, delaySpan, req.GetMeetingId())
	if err != nil {
		metrics.IncrCounter("rpc.SendDelayJob.fail", 1)
		log.ErrorContextf(ctx, "HandleWarmUpVideoDelayJob SendDelayJob failed, req:%+v, err:%+v", req, err)
		return err
	}
	log.InfoContextf(ctx, "HandleWarmUpVideoDelayJob SendDelayJob succ, req:%+v, delaySpan:%+v", req, delaySpan)
	return nil
}

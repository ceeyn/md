package service

import (
	"context"
	"encoding/json"
	"errors"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"meeting_template/model"
	"meeting_template/rpc"
	"meeting_template/util"
	"strconv"
	"strings"

	"git.code.oa.com/trpc-go/trpc-go/log"
	safePb "git.code.oa.com/trpcprotocol/wemeet/wemeet_safe_gateway"
	"github.com/golang/protobuf/proto"
)

// VideoCensorCallback ...
func (s *WemeetMeetingTemplateOidbServiceImpl) VideoCensorCallback(ctx context.Context, req *safePb.GetVideoCallbackReq,
	rsp *safePb.GetVideoCallbackRsp) error {
	bgTime := util.NowMs()
	util.ReportOne(util.VideoCensorCallback) //[VideoCensorCallback]请求
	rst, err := HandleVideoCensorResult(ctx, req)
	rsp.ErrorCode = proto.Int32(rst)
	cost := util.NowMs() - bgTime
	if err != nil {
		util.ReportOne(util.VideoCensorCallbackFail) //[VideoCensorCallback]请求失败
		rsp.ErrorMsg = proto.String(err.Error())
		log.ErrorContextf(ctx, "VideoCentsorCallback fail req:%v, rsp:%v, cost:%vms", req, rsp, cost)
	} else {
		rsp.ErrorMsg = proto.String("ok")
		log.InfoContextf(ctx, "VideoCentsorCallback ok  req:%v, rsp:%v, cost:%vms", req, rsp, cost)
	}
	return nil
}

// HandleVideoCensorResult
func HandleVideoCensorResult(ctx context.Context, req *safePb.GetVideoCallbackReq) (ret int32, err error) {
	resultCode := req.GetResultCode()
	if resultCode != util.VideoAuditResultPass && resultCode != util.VideoAuditResultFail {
		err = errors.New(
			"video callback resultCode invalid, resultCode:" + strconv.FormatUint(uint64(resultCode), 10))
		return 1, err
	}
	defer ProcessVideoAuditResult(ctx, resultCode, req.GetSer().GetTraceId())
	return 0, nil
}

// ProcessVideoAuditResult
func ProcessVideoAuditResult(ctx context.Context, resultCode uint32, uuid string) {
	if uuid == "" {
		util.ReportOne(util.VideoCensorCallbakTraceIdEmpty)
		log.ErrorContextf(ctx, "GetVideoCallbackReq trace id empty")
		return
	}
	infos := strings.Split(uuid, "&")
	if len(infos) != 2 {
		util.ReportOne(util.VideoCensorCallbakTraceIdFormatError)
		log.ErrorContextf(ctx, "GetVideoCallbackReq trace id format error, trace_id:%v", uuid)
		return
	}
	templateId := infos[0]
	videoCosId := infos[1]
	templateInfo := &model.TemplateInfo{}
	templateInfo, err := GetTemplateInfoSingleFlight(ctx, templateId)
	if err != nil {
		log.ErrorContextf(ctx, "ProcessVideoAuditResult, get templateInfo from redis failed, templateId:%v",
			templateId)
		return
	}
	if templateInfo.WarmUpData == "" {
		util.ReportOne(util.VideoCensorCallbakWarmupDataEmpty)
		log.ErrorContextf(ctx, "ProcessVideoAuditResult, warmup data empty, templateId:%v!", templateId)
		return
	}
	buf := []byte(templateInfo.WarmUpData)
	warmUpData := pb.WarmUpData{}
	err = json.Unmarshal(buf, &warmUpData)
	if err != nil {
		util.ReportOne(util.JsonUnmarshalFail)
		log.ErrorContextf(ctx, "ProcessVideoAuditResult, json parse warmup fail, warmUpData:%+v,error:%+v",
			templateInfo.WarmUpData, err)
		return
	}
	// 更新视频审核状态
	ModifyWarmUpDataVideoStatus(ctx, resultCode, videoCosId, &warmUpData)
	warmUpDataStr, err := util.GetSerializedJsonStr(ctx, warmUpData)
	if err != nil {
		util.ReportOne(util.JsonMarshalFail)
		log.InfoContextf(ctx, "ProcessVideoAuditResult, get warmUpDataStr failed, warmUpData::%v, err:%v",
			warmUpData, err)
	} else {
		templateInfo.WarmUpData = warmUpDataStr
	}
	err = SetTemplateInfo(ctx, templateInfo)
	if err != nil {
		log.ErrorContextf(ctx, "ProcessVideoAuditResult set templateInfo fail,templateId:%v", templateId)
	}
	meetingId, err := strconv.ParseUint(templateInfo.MeetingId, 10, 64)
	if err != nil {
		util.ReportOne(util.PaseUintMeetingIdFail)
		log.ErrorContextf(ctx, "ProcessVideoAuditResult PaseUintMeetingIdFailed")
	}
	appId, err := strconv.ParseUint(templateInfo.AppId, 10, 32)
	if err != nil {
		util.ReportOne(util.PaseUintAppIdFail)
		log.ErrorContextf(ctx, "ProcessVideoAuditResult PaseUintAppIdFailed")
	}
	uint32AppId := uint32(appId)
	UserNotifyWarmUpDataUpdate(ctx, meetingId, uint32AppId, templateInfo.AppUid, templateInfo)

	// 审核不通过时发通知
	if resultCode == util.VideoAuditResultFail {
		rpc.SendMessageToMsgBox(ctx, templateInfo.AppUid, uint32AppId, meetingId)
	}
}

// ModifyWarmUpDataVideoStatus
func ModifyWarmUpDataVideoStatus(ctx context.Context, resultCode uint32, videoCosId string, data *pb.WarmUpData) {
	//  暂只支持单个视频
	if len(data.GetWarmupVideoList()) != 1 {
		util.ReportOne(util.WarmUpVideoDataFormatError)
		log.InfoContextf(ctx, "ModifyWarmUpDataVideoStatus, warmUpData video size invalid, size: %v",
			len(data.GetWarmupVideoList()))
		return
	}
	decodedCosId := string(util.GetBase64Decoded(ctx, data.GetWarmupVideoList()[0].GetCosId()))
	if decodedCosId != videoCosId {
		util.ReportOne(util.CensorCallbackCosIdExpired)
		log.InfoContextf(ctx, "ModifyWarmUpDataVideoStatus, cosId not same, from callback:%v,from redis:%v",
			videoCosId, decodedCosId)
		return
	}
	if resultCode == util.VideoAuditResultPass {
		util.ReportOne(util.CensorCallbackAuditPass)
		data.GetWarmupVideoList()[0].Status = proto.Uint32(util.VideoPassAudit)
	}
	if resultCode == util.VideoAuditResultFail {
		util.ReportOne(util.CensorCallbackAuditFail)
		data.GetWarmupVideoList()[0].Status = proto.Uint32(util.VideoAuditFail)
	}
}

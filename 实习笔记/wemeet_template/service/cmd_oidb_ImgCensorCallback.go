package service

import (
	"context"
	"encoding/json"
	"errors"
	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	safePb "git.code.oa.com/trpcprotocol/wemeet/wemeet_safe_gateway"
	"google.golang.org/protobuf/proto"
	"meeting_template/model"
	"meeting_template/util"
	"strconv"
	"strings"
)

// ImgCensorCallback ...
func (s *WemeetMeetingTemplateOidbServiceImpl) ImgCensorCallback(ctx context.Context, req *safePb.GetImageCallbackReq,
	rsp *safePb.GetImageCallbackRsp) error {
	rst, err := HandleImgCensorResult(ctx, req)
	rsp.ErrorCode = proto.Int32(rst)
	if err != nil {
		attr.AttrAPI(35907579, 1)
		rsp.ErrorMsg = proto.String(err.Error())
		log.ErrorContextf(ctx, "ImgCensorCallback fail req:%+v, rsp:%+v", req, rsp)
	} else {
		attr.AttrAPI(35907580, 1)
		rsp.ErrorMsg = proto.String("ok")
		log.InfoContextf(ctx, "ImgCensorCallback ok  req:%+v, rsp:%+v", req, rsp)
	}
	return nil
}

// HandleImgCensorResult ...
func HandleImgCensorResult(ctx context.Context, req *safePb.GetImageCallbackReq) (ret int32, err error) {
	resultCode := req.GetResultCode()
	if resultCode != util.VideoAuditResultPass && resultCode != util.VideoAuditResultFail {
		err = errors.New(
			"img callback resultCode invalid, resultCode:" + strconv.FormatUint(uint64(resultCode), 10))
		return 1, err
	}
	defer ProcessImgAuditResult(ctx, resultCode, req.GetSer().GetTraceId())
	return 0, nil
}

// ProcessImgAuditResult ...
func ProcessImgAuditResult(ctx context.Context, resultCode uint32, uuid string) {
	if uuid == "" {
		util.ReportOne(util.VideoCensorCallbakTraceIdEmpty)
		log.ErrorContextf(ctx, "GetImgCallbackReq trace id empty")
		return
	}
	infos := strings.Split(uuid, "&")
	if len(infos) != 2 {
		util.ReportOne(util.VideoCensorCallbakTraceIdFormatError)
		log.ErrorContextf(ctx, "GetImgCallbackReq trace id format error, trace_id:%v", uuid)
		return
	}
	templateId := infos[0]
	cosId := infos[1]
	templateInfo := &model.TemplateInfo{}
	templateInfo, err := GetTemplateInfoSingleFlight(ctx, templateId)
	if err != nil {
		log.ErrorContextf(ctx, "ProcessImgAuditResult, get templateInfo from redis failed, templateId:%v",
			templateId)
		return
	}
	if templateInfo.WarmUpData == "" {
		util.ReportOne(util.VideoCensorCallbakWarmupDataEmpty)
		log.ErrorContextf(ctx, "ProcessImgAuditResult, warmup data empty, templateId:%v!", templateId)
		return
	}
	buf := []byte(templateInfo.WarmUpData)
	warmUpData := pb.WarmUpData{}
	err = json.Unmarshal(buf, &warmUpData)
	if err != nil {
		util.ReportOne(util.JsonUnmarshalFail)
		log.ErrorContextf(ctx, "ProcessImgAuditResult, json parse warmup fail, warmUpData:%+v,error:%+v",
			templateInfo.WarmUpData, err)
		return
	}
	// 更新图片审核状态
	ModifyWarmUpDataImgStatus(ctx, resultCode, cosId, &warmUpData)
	warmUpDataStr, err := util.GetSerializedJsonStr(ctx, warmUpData)
	if err != nil {
		log.InfoContextf(ctx, "ProcessImgAuditResult, get warmUpDataStr failed, warmUpData::%v, err:%v",
			warmUpData, err)
	} else {
		templateInfo.WarmUpData = warmUpDataStr
	}
	err = SetTemplateInfo(ctx, templateInfo)
	if err != nil {
		log.ErrorContextf(ctx, "ProcessImgAuditResult set templateInfo fail,templateId:%v", templateId)
	}
	meetingId, err := strconv.ParseUint(templateInfo.MeetingId, 10, 64)
	if err != nil {
		log.ErrorContextf(ctx, "ProcessImgAuditResult PaseUintMeetingIdFailed, templateInfo:%+v", templateInfo)
	}
	appId, err := strconv.ParseUint(templateInfo.AppId, 10, 32)
	if err != nil {
		util.ReportOne(util.PaseUintAppIdFail)
		log.ErrorContextf(ctx, "ProcessImgAuditResult PaseUintAppIdFailed")
	}
	uint32AppId := uint32(appId)
	UserNotifyWarmUpDataUpdate(ctx, meetingId, uint32AppId, templateInfo.AppUid, templateInfo)
}

// ModifyWarmUpDataImgStatus ...
func ModifyWarmUpDataImgStatus(ctx context.Context, resultCode uint32, cosId string, data *pb.WarmUpData) {
	//  暂只支持单个图片
	if len(data.GetWarmupImgList()) != 1 {
		log.InfoContextf(ctx, "ModifyWarmUpDataImgStatus, warmUpData img size invalid, size: %v",
			len(data.GetWarmupVideoList()))
		return
	}
	decodedCosId := string(util.GetBase64Decoded(ctx, data.GetWarmupImgList()[0].GetUrl()))
	if decodedCosId != cosId {
		log.InfoContextf(ctx, "ModifyWarmUpDataImgStatus, cosId not same, from callback:%v, from redis:%v",
			cosId, decodedCosId)
		return
	}
	if resultCode == util.ImgAuditResultPass {
		data.GetWarmupImgList()[0].Status = proto.Uint32(util.VideoPassAudit)
	}
	if resultCode == util.ImgAuditResultFail {
		data.GetWarmupImgList()[0].Status = proto.Uint32(util.VideoAuditFail)
		data.GetWarmupImgList()[0].Url = proto.String("")
		data.GetWarmupImgList()[0].Name = proto.String("")
	}
}

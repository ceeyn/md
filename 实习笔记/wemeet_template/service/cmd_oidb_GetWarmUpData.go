package service

import (
	"context"
	"encoding/base64"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	"meeting_template/config/config_rainbow"
	"meeting_template/model"
	"meeting_template/rpc"
	"strconv"
	"strings"
	"time"

	"meeting_template/util"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"github.com/golang/protobuf/proto"

	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
)

const SLOW_MUSIC = 1

// GetWarmUpData ... oidb协议 客户端调用 application服务触发请求
func (s *WemeetMeetingTemplateOidbServiceImpl) GetWarmUpData(ctx context.Context, req *pb.GetWarmUpDataReq,
	rsp *pb.GetWarmUpDataRsp) error {
	bgTime := util.NowMs()
	util.ReportOne(util.GetWarmUpDataReq) //[GetWarmUpData]请求
	rst, err := HandleGetWarmUpData(ctx, req, rsp)
	rsp.ErrorCode = proto.Int32(rst)
	cost := util.NowMs() - bgTime
	if err != nil {
		util.ReportOne(util.GetWarmUpDataFail) //[GetWarmUpData]请求失败
		rsp.ErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "GetWarmUpData oidb fail, req:%v, rsp:%v, cost:%vms", req, rsp, cost)
	} else {
		rsp.ErrorMessage = proto.String("ok")
		log.InfoContextf(ctx, "GetWarmUpData oidb ok, req:%v, rsp:%v, cost:%vms", req, rsp, cost)
	}
	//兼容客户端，视频没处理好的话用图片物料类型
	warmData := rsp.GetData()
	if warmData.GetUint32WarmupUseType() == util.WarmUpShowVideo && len(warmData.GetWarmupVideoList())>0 {
		videoItem := warmData.GetWarmupVideoList()[0]
		if videoItem.GetUrl() == "" && videoItem.GetStreamUrl() == "" {
			warmData.Uint32WarmupUseType = proto.Uint32(util.WarmUpShowImg)
		}
	}
	log.InfoContextf(ctx, "GetWarmUpData oidb ok, finally return req:%v, rsp:%v, cost:%vms", req, rsp, cost)
	return nil
}

// HandleGetWarmUpData
func HandleGetWarmUpData(ctx context.Context, req *pb.GetWarmUpDataReq,
	rsp *pb.GetWarmUpDataRsp) (ret int32, err error) {
	rsp.Data = &pb.WarmUpData{}
	if req.GetTemplateId() == "" {
		util.ReportOne(util.GetTemplateIdEmpty)
		log.InfoContextf(ctx, "tmId empty, return null data")
		return 0, nil
	}
	templateInfo := &model.TemplateInfo{}
	templateInfo, err = GetTemplateInfoSingleFlight(ctx, req.GetTemplateId())
	log.InfoContextf(ctx, "HandleGetWarmUpData, req:%v, templateId:%v", req, req.GetTemplateId())
	if err != nil {
		rsp.ErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "GetTemplateInfo from redis failed, req:%v, rsp:%v", req, rsp)
		return -1, err // redis读取失败，返回失败
	}
	imageType := req.GetPictureType()
	if imageType == "" {
		imageType = util.RawImageType // 前端未赋值时，暖场场景下默认获取原图
	}
	var allNeedTransCosIds []string
	var coverListNeedTranCosIds []string
	GetOtherCosIdsFromTemplateInfo(ctx, templateInfo, &coverListNeedTranCosIds)
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
	GetAndStoreDownloadUrls(ctx, imageType, allNeedTransCosIds, warmUpImgCosIds, warmUpVideoCosIds,
		warmUpVideoStreamCosIds, cosIdUrlMap)
	if ConvertWarmUpData(ctx, templateInfo, cosIdUrlMap, rsp.Data, false) {
		ReplaceWarmUpImgWithCover(ctx, templateInfo, cosIdUrlMap, rsp.Data)
		ProcessVidoSatus(ctx, templateInfo, rsp.Data)
	}
	FillWarmUpOtherRspData(ctx, req, rsp)
	log.InfoContextf(ctx, "Get WarmUpData templateId:[%v], rsp:%+v", req.GetTemplateId(), rsp)
	return 0, err
}

// 填充其它回包数据
func FillWarmUpOtherRspData(ctx context.Context, req *pb.GetWarmUpDataReq, rsp *pb.GetWarmUpDataRsp) {
	timestamp := uint32(time.Now().Unix())
	nonce := int32(timestamp)
	rsp.Nonce = proto.Int32(nonce)
	rsp.Timestamp = proto.Uint32(timestamp)
	if req.GetMeetingId() != 0 {
		// 暖场中邀请开关按钮
		key := util.MakeWarmUpInviteSwitchKey(req.GetMeetingId())
		switchState, _ := rpc.RDGetInviteSwitch(ctx, key)
		rsp.EnableInvite = proto.Uint32(uint32(switchState))
	}
}

// ConvertWarmUpData ... 转换暖场视频数据， 成功返回true, 失败返回false, 区分一下是客户端来的请求还是自web cgi的查询请求
func ConvertWarmUpData(ctx context.Context, templateInfo *model.TemplateInfo,
	cosIdUrlMap map[string]string, outWarmUpData *pb.WarmUpData, isFromCgiGetReq bool) bool {
	// 音乐列表的数据先填充上, 音乐列表不管怎么都是需要的
	FillWarmUpMusicData(ctx, outWarmUpData)
	warmUpData := pb.WarmUpData{}
	if !GetWarmUpDataFromTemplateInfo(ctx, templateInfo, &warmUpData) {
		return false
	}
	log.InfoContextf(ctx, "ConvertWarmUpData, warm up data:[%+v], img:%+v, video:%+v, type:%+v",
		warmUpData, warmUpData.WarmupImgList, warmUpData.WarmupVideoList, *warmUpData.Uint32WarmupUseType)
	// 物料类型
	outWarmUpData.Uint32WarmupUseType = proto.Uint32(warmUpData.GetUint32WarmupUseType())
	// 播放的音乐类型
	outWarmUpData.Uint32PlayMusicType = proto.Uint32(warmUpData.GetUint32PlayMusicType())
	var rawImgCosIds []string
	var imageUrls []string
	for _, val := range warmUpData.WarmupImgList {
		decodedRaw := util.GetBase64Decoded(ctx, *val.Url)
		rawImgCosIds = append(rawImgCosIds, string(decodedRaw))
		imgItem := &pb.WarmUpImgItem{}
		imgItem.OriginCosId = proto.String(*val.Url)
		imgItem.Name = val.Name
		imgItem.Status = val.Status    // 审核状态
		imgItem.Url = proto.String("") // 赋值空的url
		outWarmUpData.WarmupImgList = append(outWarmUpData.WarmupImgList, imgItem)
	}
	log.InfoContextf(ctx, "ConvertWarmUpData,  outWarmUpData:%v", outWarmUpData)
	GetUrlsFromDataMap(ctx, cosIdUrlMap, rawImgCosIds, &imageUrls)
	tranImgSucc := (len(warmUpData.WarmupImgList) == len(imageUrls))
	log.InfoContextf(ctx, "ConvertWarmUpData,length of rawImgCosIds|imageUrls:%v|%v, tranImgSucc:%v",
		len(rawImgCosIds), len(imageUrls), tranImgSucc)
	if tranImgSucc {
		for idx, val := range outWarmUpData.WarmupImgList {
			// 兜底 没有回调的话，还是展示出来
			if isFromCgiGetReq || val.GetStatus() != util.AuditFail { // 图片状态是审核通过的话，返回url链接，这里主要是客户端用到
				val.Url = proto.String(base64.StdEncoding.EncodeToString([]byte(imageUrls[idx])))
			}
			//val.Url = proto.String(base64.StdEncoding.EncodeToString([]byte(imageUrls[idx])))
		}
	}
	var rawVideoCosIds []string
	var rawStreamVideoCosIds []string
	var videoUrls []string
	var streamVideoUrls []string
	for _, val := range warmUpData.WarmupVideoList {
		decodedRaw := string(util.GetBase64Decoded(ctx, *val.CosId))
		rawVideoCosIds = append(rawVideoCosIds, decodedRaw)
		rawStreamVideoCosIds = append(rawStreamVideoCosIds,
			strings.Replace(decodedRaw, ".mp4", ".m3u8", 1))
		videoItem := &pb.WarmUpVideoItem{}
		videoItem.Name = val.Name
		videoItem.CosId = val.CosId
		videoItem.Status = val.Status
		videoItem.VideoTransStatus  = proto.Uint32(val.GetVideoTransStatus())   // 新加转码状态
		videoItem.UpdateTimeStamp   = proto.Uint64(val.GetUpdateTimeStamp())    // 新加变更时间
		outWarmUpData.WarmupVideoList = append(outWarmUpData.WarmupVideoList, videoItem)
	}
	GetUrlsFromDataMap(ctx, cosIdUrlMap, rawVideoCosIds, &videoUrls)
	GetUrlsFromDataMap(ctx, cosIdUrlMap, rawStreamVideoCosIds, &streamVideoUrls)
	tranVideoSucc := len(warmUpData.WarmupVideoList) == len(videoUrls)
	tranStreamVideoSucc := len(warmUpData.WarmupVideoList) == len(streamVideoUrls)
	log.InfoContextf(ctx, "ConvertWarmUpData,length of rawVideoCosIds|videoUrls:%v|%v, transVideoSucc:%v, tmpId:%v",
		len(rawVideoCosIds), len(videoUrls), tranVideoSucc, templateInfo.TemplateId)
	log.InfoContextf(ctx, "ConvertWarmUpData,length of rawVideoCosIds|streamVideaUrls:%v|%v, result:%v, tmpId:%v",
		len(rawVideoCosIds), len(videoUrls), tranStreamVideoSucc, templateInfo.TemplateId)

	if tranVideoSucc && tranStreamVideoSucc{
		for idx, val := range outWarmUpData.WarmupVideoList {
			val.Url = proto.String("")            // mp4
			val.StreamUrl = proto.String("")      // m3u8
			// 审核通过且转码成功，才给出去mp4和m3u8 链接。 否则mp4和m3u8链接一起置为空，展示图片
			if JudgeCanShowVideoUrl(ctx, val, templateInfo.TemplateId) {
				val.Url = proto.String(base64.StdEncoding.EncodeToString([]byte(videoUrls[idx])))
				val.StreamUrl = proto.String(base64.StdEncoding.EncodeToString([]byte(streamVideoUrls[idx])))
			}
		}
	}
	return true
}

// ReplaceWarmUpImgWithCover
func ReplaceWarmUpImgWithCover(ctx context.Context, templateInfo *model.TemplateInfo,
	cosIdUrlMap map[string]string, warmUpData *pb.WarmUpData) {
	useVideo := (warmUpData.GetUint32WarmupUseType() == util.WarmUpShowVideo && 0 < len(warmUpData.WarmupVideoList))
	log.InfoContextf(ctx, "ReplaceWarmUpImgCover,useVideo:%v,templateId:%v", useVideo, templateInfo.TemplateId)
	if 0 < len(warmUpData.WarmupImgList) || useVideo {
		return
	}
	// 暖场视频展示使用不了视频且未上传图片，使用第一张封面图片
	var imgUrl string
	if templateInfo.CoverList == "" {
		if templateInfo.CoverUrl == "" {
			return
		}
		decodedCosId := util.GetBase64Decoded(ctx, templateInfo.CoverUrl)
		imgUrl = base64.StdEncoding.EncodeToString([]byte(cosIdUrlMap[string(decodedCosId)]))
	} else {
		coverItems := make([]pb.CoverItem, 0)
		if !GetCoverItemsFromTemplateInfo(ctx, templateInfo, &coverItems) {
			return
		}
		log.InfoContextf(ctx, "covert warmup data, cover Items:[%v]", coverItems)
		var needTransRawCosIds []string
		var needTransCuttedCosIds []string
		var rawImageUrls []string
		var cuttedImageUrls []string
		for _, val := range coverItems {
			decodedRaw := util.GetBase64Decoded(ctx, *val.RawUrl)
			needTransRawCosIds = append(needTransRawCosIds, string(decodedRaw))
			decodedCutted := util.GetBase64Decoded(ctx, *val.CuttedUrl)
			needTransCuttedCosIds = append(needTransCuttedCosIds, string(decodedCutted))
		}
		GetUrlsFromDataMap(ctx, cosIdUrlMap, needTransRawCosIds, &rawImageUrls)
		GetUrlsFromDataMap(ctx, cosIdUrlMap, needTransCuttedCosIds, &cuttedImageUrls)
		if len(rawImageUrls) != len(cuttedImageUrls) || len(rawImageUrls) != len(coverItems) {
			log.InfoContextf(ctx, "list data invalid, not fill list data!")
			util.ReportOne(util.CoverListRawCuttedNumDiff)
			return
		}
		if len(coverItems) > 0 {
			imgUrl = base64.StdEncoding.EncodeToString([]byte(cuttedImageUrls[0]))
		}
	}
	imgItem := pb.WarmUpImgItem{}
	imgItem.Url = proto.String(imgUrl)
	warmUpData.WarmupImgList = append(warmUpData.WarmupImgList, &imgItem)
}

// GetDownloadUrl
func GetDownloadUrl(ctx context.Context, imageType string, cosId string, useCdn string) string {
	var needTransCosIds []string
	needTransCosIds = append(needTransCosIds, cosId)
	var imageUrls []string
	GetUrlsFromCosIds(ctx, imageType, useCdn, needTransCosIds, &imageUrls)
	if len(imageUrls) != 1 {
		log.InfoContextf(ctx, "GetDownloadUrl  failed!")
		return ""
	}
	return imageUrls[0]
}

// 判断视频审核状态，为"未审核"时送审即可
func ProcessVidoSatus(ctx context.Context, templateInfo *model.TemplateInfo, warmUpData *pb.WarmUpData) {
	if warmUpData.GetUint32WarmupUseType() == util.WarmUpShowImg || len(warmUpData.WarmupVideoList) != 1 {
		log.InfoContextf(ctx, "Process video status, use pic or no video, not need process!")
		return
	}
	status := warmUpData.WarmupVideoList[0].GetStatus()
	if status == util.VideoNotAudit {
		cosId := warmUpData.WarmupVideoList[0].GetCosId()
		decodedCosId := string(util.GetBase64Decoded(ctx, cosId))
		auditUuid := (templateInfo.TemplateId + "&" + decodedCosId)
		log.InfoContextf(ctx, "ProcessVidoSatus,templateId:%v,autidUuid:%v", templateInfo.TemplateId, auditUuid)
		appId, err := strconv.ParseUint(templateInfo.AppId, 10, 32)
		if err != nil {
			util.ReportOne(util.PaseUintAppIdFail)
			log.ErrorContextf(ctx, "ProcessVideoStatus PaseUintAppIdFailed")
		}
		meetingId, err := strconv.ParseUint(templateInfo.MeetingId, 10, 64)
		if err != nil {
			util.ReportOne(util.PaseUintMeetingIdFail)
			log.ErrorContextf(ctx, "ProcessVideoStatus PaseUintMeetingIdFailed")
		}
		url := GetDownloadUrl(ctx, util.NotPicImageType, decodedCosId, util.DownloadUseCdn)
		callBackSwitch := config_rainbow.GetCallBackConf().WarmUpShowVideoAction
		log.InfoContextf(ctx,"ctlrpcImgSafetyAudit WarmUpShowVideoAction: %v",callBackSwitch)
		if callBackSwitch {
			warmVideoReq := &rpc.SafetyAuditReq{
				Uuid:              auditUuid,
				AppId:             uint32(appId),
				AppUid:            templateInfo.AppUid,
				MeetingId:         meetingId,
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
		}else{
			rpc.VideoSafetyAudit(ctx, auditUuid, uint32(appId), templateInfo.AppUid, meetingId, url)
		}

	}
}

// FillWarmUpMusicData ...
func FillWarmUpMusicData(ctx context.Context, warmUpData *pb.WarmUpData) {
	softMusic := &pb.WarmUpMusicItem{}
	softMusic.Name = proto.String("有")
	softMusic.MusicType = proto.Uint32(SLOW_MUSIC) // 有
	softMusic.Url = proto.String("")

	musicUrlMap := make(map[string]string)
	musicCfg := config_rainbow.GetMusicConfConfig() // 七彩石配置
	slowMusicCosId := musicCfg.SlowMusic

	musicCosIds := []string{slowMusicCosId}
	GetMapDataFromCosIds(ctx, util.NotPicImageType, util.DownloadUseCdn, musicCosIds, musicUrlMap)
	log.InfoContextf(ctx, "FillWarmUpMusicData get musicUrlMap:%+v", musicUrlMap)

	if val, ok := musicUrlMap[slowMusicCosId]; ok {
		base64Url := base64.StdEncoding.EncodeToString([]byte(val))
		softMusic.Url = proto.String(base64Url)
	}

	musicList := []*pb.WarmUpMusicItem{softMusic}
	warmUpData.WarmupMusicList = musicList

}

// JudgeCanShowVideoUrl 判断能否展示视频url
func JudgeCanShowVideoUrl(ctx context.Context, videoItemInfo *pb.WarmUpVideoItem, tplId string) bool {
	// 审核不通过，是不展示视频
	if videoItemInfo.GetStatus() != util.VideoPassAudit {
		log.InfoContextf(ctx,"JudgeCanShowVideoUrl video not pass Audit, tplId:%+v, videoItemInfo:%+v",
			tplId, videoItemInfo)
		return false
	}

	if videoItemInfo.GetVideoTransStatus() == util.M3U8TransDone {
		log.InfoContextf(ctx,"JudgeCanShowVideoUrl video have trans success, tplId:%+v, videoItemInfo:%+v",
			tplId, videoItemInfo)
		return true
	}

	// 延时消息还没消费到，客户端先进暖场，这时候触发的查询发现转码成功了，给出视频链接。
	decodedMp4CosId := string(util.GetBase64Decoded(ctx, videoItemInfo.GetCosId()))
	m3u8CosId := strings.Replace(decodedMp4CosId, ".mp4", ".m3u8", 1)
	isM3U8Exist := rpc.JudgeCosResourceIsExist(ctx, m3u8CosId)    //查m3u8的转码状态
	if isM3U8Exist {
		metrics.IncrCounter("m3u8Exist.sum",1)
		log.InfoContextf(ctx, "JudgeCanShowVideoUrl query m3u8 resource isExist, tplId:%+v, videoItemInfo:%+v",
			tplId, videoItemInfo)
		return true
	}

	// 之前上传的视频，兼容一下
	if videoItemInfo.GetVideoTransStatus() == 0  {
		metrics.IncrCounter("videoItemInfo.NoVideoTransStatus",1)
		log.InfoContextf(ctx,"JudgeCanShowVideoUrl video not have TransStatus, tplId:%+v, videoItemInfo:%+v",
			tplId, videoItemInfo)
		return true
	}
	if videoItemInfo.GetUpdateTimeStamp() == 0 {
		metrics.IncrCounter("videoItemInfo.NoUpdateTimeStamp",1)
		log.InfoContextf(ctx,"JudgeCanShowVideoUrl video not have UpdateTimeStamp, tplId:%+v, videoItemInfo:%+v",
			tplId, videoItemInfo)
		return true
	}
	// 兜底逻辑，一直没有查询到转码完成
	timeInterval := uint64(time.Now().Unix()) - videoItemInfo.GetUpdateTimeStamp()
	delaySpan, maxRetryCnt := config_rainbow.GetMusicConfConfig().DelaySpan, config_rainbow.GetMusicConfConfig().MaxRetryCnt
	if videoItemInfo.GetVideoTransStatus() == util.M3U8TransNoReady && timeInterval > uint64(delaySpan*maxRetryCnt) {
		metrics.IncrCounter("videoItemInfo.IntervalOverMax",1)
		log.InfoContextf(ctx,"JudgeCanShowVideoUrl video not trans success, but over max time. tplId:%+v, " +
			"videoItemInfo:%+v, delaySpan:%+v, maxRetryCnt%+v", tplId, videoItemInfo, delaySpan, maxRetryCnt)
		return true
	}
    metrics.IncrCounter("JudgeCanShowVideoUrl.VideoNotTrans",1)
	log.InfoContextf(ctx,"JudgeCanShowVideoUrl video dont have trans success, no show videoUrl, tplId:%+v, " +
		"videoItemInfo:%+v", tplId, videoItemInfo)
	return false
}
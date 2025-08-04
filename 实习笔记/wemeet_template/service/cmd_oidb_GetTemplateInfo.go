package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"meeting_template/model"
	"meeting_template/rpc"
	"meeting_template/util"
	t_extractor "meeting_template/util/t-helper/t-html/t-extractor"
	t_replacer "meeting_template/util/t-helper/t-html/t-replacer"

	// "gonum.org/v1/gonum/mat"
	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
	commonUploadPb "git.code.oa.com/trpcprotocol/wemeet/common_upload"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/golang/protobuf/proto"
)

const CoverPlayInterval = 3 // 图片轮播间隔暂设置为3s
const ReqCommonUpladFeature = "wemeet_webinar"

// GetTemplateInfo ...
func (s *WemeetMeetingTemplateOidbServiceImpl) GetTemplateInfo(ctx context.Context, req *pb.GetTemplateInfoReq,
	rsp *pb.GetTemplateInfoRsp) error {
	bgTime := util.NowMs()
	util.ReportOne(util.GetTemplateInfoReq) //[GetTemplateInfo]请求

	rst, err := HandleGetTemplateInfo(ctx, req, rsp)
	rsp.ErrorCode = proto.Int32(rst)
	cost := util.NowMs() - bgTime
	if err != nil {
		util.ReportOne(util.GetTemplateInfoFail) //[GetTemplateInfo]请求失败
		rsp.ErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "GetTemplateInfo fail req:%v, rsp:%v, cost:%vms", req, rsp, cost)
	} else {
		rsp.ErrorMessage = proto.String("ok")
		log.InfoContextf(ctx, "GetTemplateInfo ok  req:%v, rsp:%v, cost:%vms", req, rsp, cost)
	}
	return nil
}

// DoGetTemplateInfo ...
func DoGetTemplateInfo(ctx context.Context, templateId string,
	pictureType string) (data *pb.TemplateInfo, ret int32, err error) {
	bgTime := util.NowMs()
	data = &pb.TemplateInfo{}
	data.WarmUpData = &pb.WarmUpData{}
	if templateId == "" {
		util.ReportOne(util.GetTemplateIdEmpty)
		log.InfoContextf(ctx, "tmId empty, return null data")
		return data, 0, nil
	}
	ttInfo := &model.TemplateInfo{}
	ttInfo, err = GetTemplateInfoSingleFlight(ctx, templateId)
	getRedisCost := util.NowMs() - bgTime
	var getDownloadUrlErr string
	if err != nil {
		log.ErrorContextf(ctx, "GetTemplateInfo from redis failed, templateId:%v", templateId)
		return nil, -1, err // redis读取失败，返回失败
	}
	imageType := pictureType
	if imageType == "" {
		imageType = util.MiddleImageType // 前端未赋值时，默认获取中图
	}
	decompressed, err := util.GetDeCompressed(ctx, ttInfo.Description)
	log.DebugContextf(ctx, "DoGetTemplateInfo description:%+v, decompressed:%+v, err:%+v",
		ttInfo.Description, decompressed, err)
	if err != nil {
		util.ReportOne(util.DeCompressDescriptionFail)
		log.ErrorContextf(ctx, " DeCompress Description failed, err:%v", err)
		return nil, -2, err // redis读取详情解压出错，返回失败
	}

	//Description详情里面(老的)的图片提取出来
	decodedDescription := util.GetBase64Decoded(ctx, decompressed)
	decodedCoverUrl := util.GetBase64Decoded(ctx, ttInfo.CoverUrl)
	var needTransCosIds []string
	needTransCosIds = append(needTransCosIds, string(decodedCoverUrl))
	GetCosIdsFromDescription(ctx, decompressed, decodedDescription, &getDownloadUrlErr, data, &needTransCosIds)
	var allNeedTransCosIds []string
	allNeedTransCosIds = append(allNeedTransCosIds, needTransCosIds...)
	var coverListNeedTranCosIds []string
	GetOtherCosIdsFromTemplateInfo(ctx, ttInfo, &coverListNeedTranCosIds)
	allNeedTransCosIds = append(allNeedTransCosIds, coverListNeedTranCosIds...)

	var warmUpImgCosIds []string
	var warmUpVideoCosIds []string
	var warmUpVideoStreamCosIds []string
	warmUpData := pb.WarmUpData{}
	if GetWarmUpDataFromTemplateInfo(ctx, ttInfo, &warmUpData) {
		GetPicCosIdsFromWarmUpData(ctx, &warmUpData, &warmUpImgCosIds)
		GetVideoCosIdsFromWarmUpData(ctx, &warmUpData, &warmUpVideoCosIds, &warmUpVideoStreamCosIds)
	}
	cosIdUrlMap := make(map[string]string)
	GetAndStoreDownloadUrls(ctx, imageType, allNeedTransCosIds, warmUpImgCosIds, warmUpVideoCosIds,
		warmUpVideoStreamCosIds, cosIdUrlMap)
	var imageUrls []string
	GetUrlsFromDataMap(ctx, cosIdUrlMap, needTransCosIds, &imageUrls)
	FillDescription(ctx, ttInfo, &imageUrls, &getDownloadUrlErr, &decompressed, &decodedCoverUrl,
		&decodedDescription, data)
	//新的放cover list的图片
	FillCoverListData(ctx, ttInfo, data, cosIdUrlMap)
	if ConvertWarmUpData(ctx, ttInfo, cosIdUrlMap, data.WarmUpData, true) { // 暖场场景都使用原图
		ReplaceWarmUpImgWithCover(ctx, ttInfo, cosIdUrlMap, data.WarmUpData)
	}
	data.RequestId = proto.String(templateId)
	data.Sponsor = proto.String(ttInfo.Sponsor)
	data.CoverName = proto.String(ttInfo.CoverName)
	data.CoverPlayInterval = proto.Uint32(CoverPlayInterval)
	totalCost := util.NowMs() - bgTime
	log.InfoContextf(ctx, "tmId:[%v],d:%+v,getRedis:%vms,total:%vms", templateId, data, getRedisCost, totalCost)
	if getDownloadUrlErr != "" {
		err = errors.New(getDownloadUrlErr)
	}
	return data, 0, nil
}

// HandleGetTemplate
func HandleGetTemplateInfo(ctx context.Context, req *pb.GetTemplateInfoReq,
	rsp *pb.GetTemplateInfoRsp) (ret int32, err error) {

	rsp.Data, ret, err = DoGetTemplateInfo(ctx, req.GetTemplateId(), req.GetPictureType())

	timestamp := uint32(time.Now().Unix())
	nonce := int32(timestamp)
	rsp.Nonce = proto.Int32(nonce)
	rsp.Timestamp = proto.Uint32(timestamp)

	// 暖场中邀请开关按钮
	if req.GetMeetingId() != 0 {
		key := util.MakeWarmUpInviteSwitchKey(req.GetMeetingId())
		switchState, _ := rpc.RDGetInviteSwitch(ctx, key)
		rsp.EnableInvite = proto.Uint32(uint32(switchState))
	} else {
		attr.AttrAPI(35906218, 1)
	}

	return ret, err
}

// 将cosId转换为对应的下载链接, useCdn "0"不走cdn获取下载链接， "1"走cdn获取下载链接
func GetUrlsFromCosIds(ctx context.Context, imageType string, useCdn string, cosIds []string, urls *[]string) {
	var req commonUploadPb.BatchTempUrlReq
	req.Feature = ReqCommonUpladFeature
	req.ImageType = imageType
	req.IsCdn = useCdn
	for _, cosId := range cosIds {
		req.CosIds = append(req.CosIds, cosId)
	}
	rsp, err := rpc.BatchQueryTempUrl(ctx, &req)
	if err != nil {
		util.ReportOne(util.TransCosIdToUrlFail) //[GetUrlsFromCosIds] 异常
		log.ErrorContextf(ctx, "GetUrlsFromCosIds failed, error: %+v", err)
		return
	}
	cosIdMap := rsp.GetData().GetCosIdMap()
	for _, cosId := range cosIds {
		if cosIdMap[cosId] != "" {
			if cosId == "" {
				*urls = append(*urls, "")
			} else {
				*urls = append(*urls, cosIdMap[cosId])
			}
		}
	}
	log.InfoContextf(ctx, "GetUrlsFromCosIds req:[%v],  rsp: %+v, urls size: %+v", req, rsp, len(*urls))
}

// 从详情文本里提取图片cos id
func GetCosIdsFromDescription(ctx context.Context, decompressed string, decodedDescription []byte,
	getDownloadUrlErr *string, data *pb.TemplateInfo, needTransCosIds *[]string) {
	if len(decodedDescription) == 0 {
		return
	}
	extractor := &t_extractor.TSimpleExtractor{}
	elements, err := extractor.Parse(string(decodedDescription))
	if err != nil {
		util.ReportOne(util.HtmlParseDescriptionFail)
		log.ErrorContextf(ctx, "parse description failed, description: %+v,  error: %+v",
			decodedDescription, err)
		*getDownloadUrlErr += ("Parse decription failed!")
		data.Description = proto.String(decompressed)
	} else {
		imageCosIds := elements.Images
		*needTransCosIds = append(*needTransCosIds, imageCosIds...)
	}
}

// 填充封面列表数据
func FillCoverListData(ctx context.Context, templateInfo *model.TemplateInfo,
	data *pb.TemplateInfo, cosIdUrlMap map[string]string) {
	coverItems := make([]pb.CoverItem, 0)
	if !GetCoverItemsFromTemplateInfo(ctx, templateInfo, &coverItems) {
		log.InfoContextf(ctx, "GetCoverItemsFromTemplateInfo failed, not fill coverlist")
		return
	}
	log.InfoContextf(ctx, "FillCoverListData, cover Items:[%v]", coverItems)
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
	for idx, val := range coverItems {
		coverItem := pb.CoverItem{}
		coverItem.CropInfo = val.CropInfo
		coverItem.RawUrl = proto.String(base64.StdEncoding.EncodeToString([]byte(rawImageUrls[idx])))
		coverItem.CuttedUrl = proto.String(base64.StdEncoding.EncodeToString([]byte(cuttedImageUrls[idx])))
		cutUrlCosId := util.GetBase64Decoded(ctx, *val.CuttedUrl)
		coverItem.CuttedUrlCosId = proto.String(string(cutUrlCosId))
		log.InfoContextf(ctx, "FillCoverListData loop,templateId:[%v], coverItem:[+%v]",
			templateInfo.TemplateId, coverItem)
		data.CoverList = append(data.CoverList, &coverItem)
	}
	log.InfoContextf(ctx, "After FillCoverListData, templateId:[%v], data:[+%v]", templateInfo.TemplateId, data)
}

// GetOtherCosIdsFromTemplateInfo
func GetOtherCosIdsFromTemplateInfo(ctx context.Context, templateInfo *model.TemplateInfo, needTransCosIds *[]string) {
	decodedCoverUrl := util.GetBase64Decoded(ctx, templateInfo.CoverUrl)
	*needTransCosIds = append(*needTransCosIds, string(decodedCoverUrl))
	coverItems := make([]pb.CoverItem, 0)
	if !GetCoverItemsFromTemplateInfo(ctx, templateInfo, &coverItems) {
		return
	}
	var needTransRawCosIds []string
	var needTransCuttedCosIds []string
	for _, val := range coverItems {
		decodedRaw := util.GetBase64Decoded(ctx, *val.RawUrl)
		needTransRawCosIds = append(needTransRawCosIds, string(decodedRaw))
		decodedCutted := util.GetBase64Decoded(ctx, *val.CuttedUrl)
		needTransCuttedCosIds = append(needTransCuttedCosIds, string(decodedCutted))
	}
	*needTransCosIds = append(*needTransCosIds, needTransRawCosIds...)
	*needTransCosIds = append(*needTransCosIds, needTransCuttedCosIds...)
}

// 从templateInfo获取coverListItems， list不为空返回true, 否则返回false
func GetCoverItemsFromTemplateInfo(ctx context.Context, templateInfo *model.TemplateInfo, items *[]pb.CoverItem) bool {
	if templateInfo.CoverList == "" {
		log.InfoContextf(ctx, "GetCoverItemsFromTemplateInfo, cover list empty, templateId:[%+v]",
			templateInfo.TemplateId)
		return false
	}
	buf := []byte(templateInfo.CoverList)
	err := json.Unmarshal(buf, items)
	if err != nil {
		util.ReportOne(util.JsonUnmarshalFail)
		log.ErrorContextf(ctx, "json parse cover list failed, coverList: %+v,  error: %+v",
			templateInfo.CoverList, err)
		return false
	}
	return true
}

// GetWarmUpDataFromTemplateInfo
func GetWarmUpDataFromTemplateInfo(ctx context.Context, templateInfo *model.TemplateInfo,
	warmUpData *pb.WarmUpData) bool {
	if templateInfo.WarmUpData == "" {
		log.InfoContextf(ctx, "ConvertWarmUpData, data empty, templateId:%v!", templateInfo.TemplateId)
		return false
	}
	buf := []byte(templateInfo.WarmUpData)
	err := json.Unmarshal(buf, warmUpData)
	if err != nil {
		util.ReportOne(util.JsonUnmarshalFail)
		log.ErrorContextf(ctx, "json parse warm up data failed, warmUpData: %+v,  error: %+v",
			templateInfo.WarmUpData, err)
		return false
	}
	return true
}

// 将cosId转换为对应的下载链接并存储起来
func GetMapDataFromCosIds(ctx context.Context, imageType string, useCdn string, cosIds []string,
	cosIdUrlMap map[string]string) {
	if len(cosIds) == 0 {
		log.InfoContextf(ctx, "GetMapDataFromCosIds cosIds empty!")
		return
	}
	var req commonUploadPb.BatchTempUrlReq
	req.Feature = ReqCommonUpladFeature
	req.ImageType = imageType
	req.IsCdn = useCdn
	for _, cosId := range cosIds {
		req.CosIds = append(req.CosIds, cosId)
	}
	rsp, err := rpc.BatchQueryTempUrl(ctx, &req)
	if err != nil {
		util.ReportOne(util.TransCosIdToUrlFail) //[GetUrlsFromCosIds] 异常
		log.ErrorContextf(ctx, "GetUrlsFromCosIds failed, error: %+v", err)
		return
	}
	cosIdMap := rsp.GetData().GetCosIdMap()
	for _, cosId := range cosIds {
		if cosIdMap[cosId] != "" {
			if cosId == "" {
				cosIdUrlMap[cosId] = ""
			} else {
				cosIdUrlMap[cosId] = cosIdMap[cosId]
			}
		}
	}
	log.InfoContextf(ctx, "rpc batch query temp url rsp: %+v, map size: %+v", rsp, len(cosIdUrlMap))
}

// GetUrlsFromDataMap
func GetUrlsFromDataMap(ctx context.Context, cosIdUrlMap map[string]string, cosIds []string, urls *[]string) {
	for _, cosId := range cosIds {
		*urls = append(*urls, cosIdUrlMap[cosId])
	}
	log.InfoContextf(ctx, "GetUrlsFromDataMap, in size:%+v,out size:%+v", len(cosIds), len(*urls))
}

// MergeMaps
func MergeMaps(src map[string]string, maps ...map[string]string) {
	for _, m := range maps {
		for k, v := range m {
			src[k] = v
		}
	}
}

// GetPicCosIdsFromWarmUpData
func GetPicCosIdsFromWarmUpData(ctx context.Context, warmUpData *pb.WarmUpData, needTransCosIds *[]string) {
	for _, val := range warmUpData.WarmupImgList {
		decodedRaw := util.GetBase64Decoded(ctx, *val.Url)
		*needTransCosIds = append(*needTransCosIds, string(decodedRaw))
	}
}

// GetVideoCosIdsFromWarmUpData
func GetVideoCosIdsFromWarmUpData(ctx context.Context, warmUpData *pb.WarmUpData, needTransCosIds *[]string,
	needTransStreamCosIds *[]string) {
	for _, val := range warmUpData.WarmupVideoList {
		decodedRaw := string(util.GetBase64Decoded(ctx, *val.CosId))
		*needTransCosIds = append(*needTransCosIds, decodedRaw)
		*needTransStreamCosIds = append(*needTransStreamCosIds,
			strings.Replace(decodedRaw, ".mp4", ".m3u8", 1))
	}
}

// FillDescription
func FillDescription(ctx context.Context, ttInfo *model.TemplateInfo, imageUrls *[]string,
	getDownloadUrlErr *string, decompressed *string, decodedCoverUrl *[]byte,
	decodedDescription *[]byte, data *pb.TemplateInfo) {
	if len(*imageUrls) == 0 {
		data.CoverUrl = proto.String(ttInfo.CoverUrl)
		data.Description = proto.String(*decompressed)
	} else {
		if ttInfo.CoverUrl == "" || len(*decodedCoverUrl) == 0 {
			data.CoverUrl = proto.String(ttInfo.CoverUrl) // codId为空，或base64解码失败，使用转换前的
		} else {
			data.CoverUrl = proto.String(base64.StdEncoding.EncodeToString([]byte((*imageUrls)[0])))
		}
		if 1 < len(*imageUrls) {
			images := (*imageUrls)[1:]
			replacer := &t_replacer.TSimpleReplacer{}
			description, err := replacer.ReplaceImageSrc(string(*decodedDescription), images)
			if err != nil { // 图片下载链接替换失败，使用原cosId，这时图片无法正常展示
				util.ReportOne(util.ReplaceDiscriptionImgFail)
				log.ErrorContextf(ctx, "replace description fail, description:%+v,error: %+v", decompressed, err)
				*getDownloadUrlErr += fmt.Sprintf("|replace desp fail:%+v,error: %+v", decodedDescription, err)
				data.Description = proto.String(*decompressed)
			} else { // 富文本为不完整的html文本，go的html包会补全，增加<html><head></head><body>前缀和</body></html>后缀，需去除
				log.InfoContextf(ctx, "replace description images succeed, description: %+v", description)
				data.Description = proto.String(base64.StdEncoding.EncodeToString(
					[]byte(description[25:(len(description) - 14)])))
			}
		} else {
			data.Description = proto.String(*decompressed)
		}
	}
}

// GetAndStoreDownloadUrls
func GetAndStoreDownloadUrls(ctx context.Context, imageType string, allNeedTransCosIds []string,
	warmUpImgCosIds []string, warmUpVideoCosIds []string, warmUpVideoStreamCosIds []string,
	cosIdUrlMap map[string]string) error {
	warmUpPicMap := make(map[string]string)
	warmUpVideoMap := make(map[string]string)
	warmUpVideoStreamMap := make(map[string]string)
	err := trpc.GoAndWait(
		func() error {
			GetMapDataFromCosIds(ctx, imageType, util.DownloadUseCdn, allNeedTransCosIds, cosIdUrlMap)
			return nil
		},
		func() error {
			GetMapDataFromCosIds(ctx, util.RawImageType, util.DownloadUseCdn, warmUpImgCosIds, warmUpPicMap)
			return nil
		},
		func() error {
			GetMapDataFromCosIds(ctx, util.NotPicImageType, util.DownloadUseCdn, warmUpVideoCosIds, warmUpVideoMap)
			return nil
		},
		func() error {
			GetMapDataFromCosIds(ctx, util.NotPicImageType, util.DownloadNotUseCdn,
				warmUpVideoStreamCosIds, warmUpVideoStreamMap)
			return nil
		},
	)
	MergeMaps(cosIdUrlMap, warmUpPicMap, warmUpVideoMap, warmUpVideoStreamMap)
	return err
}

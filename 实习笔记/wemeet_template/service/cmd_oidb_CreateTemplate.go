package service

import (
	"context"
	"errors"
	"git.code.oa.com/trpc-go/trpc-go/log"
	uuid "git.code.oa.com/wesee_ugc/go.uuid"
	"github.com/golang/protobuf/proto"
	"meeting_template/model"
	"meeting_template/rpc"
	"meeting_template/util"
	"strconv"
	"time"

	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
)

//
type WemeetMeetingTemplateOidbServiceImpl struct{}

// CreateTemplate ...
func (s *WemeetMeetingTemplateOidbServiceImpl) CreateTemplate(ctx context.Context, req *pb.CreateTemplateReq,
	rsp *pb.CreateTemplateRsp) error {
	util.ReportOne(util.CreateTemplateReq) //[CreateTemplate]请求
	rst, err := HandleCreateTemplate(ctx, req, rsp)
	rsp.ErrorCode = proto.Int32(rst)
	rsp.StrNonceId = proto.String(req.GetStrNonceId())
	if rsp.GetStrNonceId() == "" {
		rsp.StrNonceId = proto.String(uuid.NewV4().String())
	}
	if err != nil {
		util.ReportOne(util.CreateTemplateFail) //[CreateTemplate]请求失败
		rsp.ErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "CreateTemplate fail req:%v, rsp:%v ", req, rsp)
	} else {
		rsp.ErrorMessage = proto.String("ok")
		log.InfoContextf(ctx, "CreateTemplate ok  req:%v, rsp:%v", req, rsp)
	}
	return nil
}

//HandleCreateTemplate ...
func HandleCreateTemplate(ctx context.Context, req *pb.CreateTemplateReq,
	rsp *pb.CreateTemplateRsp) (ret int32, err error) {

	templateInfo := &model.TemplateInfo{}
	FillTemplateDirectInfo(req, templateInfo)
	coverItems := req.GetCoverList()
	templateInfo.CoverList = ""
	if coverItems != nil {
		err = util.CheckCoverListFormat(ctx, coverItems)
		if err != nil {
			util.ReportOne(util.CoverUrlCosIdInvalid)
			return -4, err
		}
		templateInfo.CoverList, err = util.GetSerializedJsonStr(ctx, coverItems)
		if err != nil {
			util.ReportOne(util.JsonMarshalFail)
			log.InfoContextf(ctx, "CreateTemplate, get coverListStr fail,coverItems:%v,err:%v", coverItems, err)
		}
	}
	warmUpData := req.GetWarmUpData()
	templateInfo.WarmUpData = ""
	if warmUpData != nil {
		err = util.CheckWarmUpDataFormat(ctx, warmUpData)
		if err != nil {
			util.ReportOne(util.CoverUrlCosIdInvalid)
			return -4, err
		}
		templateInfo.WarmUpData, err = util.GetSerializedJsonStr(ctx, warmUpData)
		if err != nil {
			util.ReportOne(util.JsonMarshalFail)
			log.InfoContextf(ctx, "CreateTemplate,get warmUpDataStr fail,warmUpData:%v,err:%v", warmUpData, err)
		}
	}
	log.InfoContextf(ctx, "CreateTemplate, CoverListStr: %v, warmUpDataStr:%v",
		templateInfo.CoverList, templateInfo.WarmUpData)
	appId := req.GetUint32AppId()
	appUId := req.GetStrAppUid()
	decodedCoverUrl := string(util.GetBase64Decoded(ctx, templateInfo.CoverUrl))
	cosIdValid := util.IsValidCosId(ctx, decodedCoverUrl)
	if !cosIdValid {
		util.ReportOne(util.CoverUrlCosIdInvalid)
		err = errors.New("cover url cos id invalid, cosId:" + templateInfo.CoverUrl)
		log.ErrorContextf(ctx, "hit cos id format filter, cosId:%v", templateInfo.CoverUrl)
		return -4, err
	}
	hitSensitive, word := rpc.CheckLongTextSensitiveData(ctx, templateInfo.Sponsor, false,
		appId, appUId, WebinarTemplateSponsor)
	if hitSensitive {
		err = errors.New("sponsor or cover name has sensitive word, " + word)
		log.ErrorContextf(ctx, "hit sensitive filter,sponsor:%v,sensitive word:%v", templateInfo.Sponsor, word)
		return -1, err
	}
	descriptionHitSensitive, words := rpc.CheckLongTextSensitiveData(ctx, templateInfo.Description, true,
		appId, appUId, WebinarTemplateDesc)
	if descriptionHitSensitive {
		err = errors.New("description has sensitive words:" + words)
		log.ErrorContextf(ctx, "description hit sensitive filter, sensitive words:%v", words)
		return -3, err
	}
	// 详情压缩后再存redis
	compressed, err := util.GetCompressed(ctx, templateInfo.Description)
	log.DebugContextf(ctx, "HandleCreateTemplate description:%+v, compressed:%+v, err:%+v",
		templateInfo.Description, compressed, err)
	if err != nil {
		log.ErrorContextf(ctx, "description compress error:%v", err)
		return -5, err
	}
	templateInfo.Description = compressed
	err = SetTemplateInfo(ctx, templateInfo)
	if err != nil {
		log.ErrorContextf(ctx, "CreateTemplate fail templateInfo:%+v", templateInfo)
		return -2, err
	}
	// 随机暂定跟时间戳相关
	currentTime := time.Now().Unix()
	rsp.TemplateId = proto.String(templateInfo.TemplateId)
	rsp.Nonce = proto.Int32(int32(currentTime))
	rsp.Timestamp = proto.Uint32(uint32(currentTime))
	return 0, nil
}

// FillTemplateDirectInfo
func FillTemplateDirectInfo(req *pb.CreateTemplateReq, templateInfo *model.TemplateInfo) {
	templateInfo.Sponsor = req.GetSponsor()
	templateInfo.CoverName = req.GetCoverName()
	templateInfo.CoverUrl = req.GetCoverUrl()
	templateInfo.Description = req.GetDescription()
	appId := req.GetUint32AppId()
	templateInfo.AppId = strconv.FormatUint(uint64(appId), 10)
	templateInfo.AppUid = req.GetStrAppUid()
	//创建template的时候，如果请求有带meeting id，那么写入templateInfo
	if req.GetUint64MeetingId() != 0 {
		templateInfo.MeetingId = strconv.FormatUint(req.GetUint64MeetingId(), 10)
	}
}

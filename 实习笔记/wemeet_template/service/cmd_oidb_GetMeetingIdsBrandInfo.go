package service

import (
	"context"
	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/golang/protobuf/proto"
	"meeting_template/es"
	"meeting_template/util"
	"strconv"
)

// GetMeetingIdsBrandInfo ... oidb协议
func (s *WemeetMeetingTemplateOidbServiceImpl) GetMeetingIdsBrandInfo(ctx context.Context,
	req *pb.GetMeetingIdsBrandInfoReq, rsp *pb.GetMeetingIdsBrandInfoRsp) (err error) {
	attr.AttrAPI(36337026, 1) //[GetMeetingIdsBrandInfo]请求
	rst, err := HandleGetMeetingIdsBrandInfo(ctx, req, rsp)
	rsp.ErrorCode = proto.Uint32(rst)
	if err != nil {
		attr.AttrAPI(36337027, 1) //[GetMeetingIdsBrandInfo]请求失败
		rsp.ErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "GetMeetingIdsBrandInfo fail req:%v, rsp:%v", req, rsp)
	} else {
		attr.AttrAPI(36337028, 1) //[GetMeetingIdsBrandInfo]请求成功
		rsp.ErrorMessage = proto.String("ok")
		log.InfoContextf(ctx, "GetMeetingIdsBrandInfo ok  req:%v, rsp:%v", req, rsp)
	}
	return nil
}

// HandleGetMeetingIdsBrandInfo ...
func HandleGetMeetingIdsBrandInfo(ctx context.Context, req *pb.GetMeetingIdsBrandInfoReq,
	rsp *pb.GetMeetingIdsBrandInfoRsp) (uint32, error) {

	meetingIds := req.GetMeetingIds()
	if len(meetingIds) > 50 {
		return util.ErrLenMeetingId, nil
	}

	strIdList := []string{}
	for _, id := range meetingIds {
		strId := strconv.FormatUint(id, 10)
		strIdList = append(strIdList, strId)
	}

	//会议的介绍列表
	introductionList, err := es.BatchSearchMeetingIntroduction(ctx, strIdList)
	if err != nil {
		attr.AttrAPI(36337029, 1)
		log.ErrorContextf(ctx, "HandleGetMeetingIdsBrandInfo BatchSearchMeetingIntroduction failed. ids:%+v, "+
			"err:%+v", strIdList, err)
		return util.ErrSearchIntroduction, err
	}
	meetingIdIntroductionMap := make(map[string]string)
	for _, introduction := range introductionList {
		meetingId := introduction.MeetingId
		meetingIdIntroductionMap[meetingId] = introduction.MeetingIntroduction
	}

	//会议的日程
	itineraryList, err := es.BatchSearchMeetingItineraryList(ctx, strIdList)
	if err != nil {
		attr.AttrAPI(36337030, 1)
		log.ErrorContextf(ctx, "HandleGetMeetingIdsBrandInfo BatchSearchMeetingItineraryList failed. ids:%+v, "+
			"err:%+v", strIdList, err)
		return util.ErrSearchItinerary, err
	}
	meetingIdItineraryMap := make(map[string][]*pb.WebinarItineraryInfo)
	for _, val := range itineraryList {
		itineraryInfo := &pb.WebinarItineraryInfo{}
		strMeetingId := val.MeetingId
		uint64MeetingId, _ := strconv.ParseUint(strMeetingId, 10, 64)
		itineraryInfo.MeetingId = proto.Uint64(uint64MeetingId)
		itineraryInfo.ItineraryName = proto.String(val.ItineraryName)
		itineraryInfo.ItineraryIntroduction = proto.String(val.ItineraryIntroduction)
		meetingIdItineraryMap[strMeetingId] = append(meetingIdItineraryMap[strMeetingId], itineraryInfo)
	}

	webinarBrandInfoList := []*pb.WebinarBrandInfoDetail{}

	for _, strId := range strIdList {
		webinarBrandInfoDetail := &pb.WebinarBrandInfoDetail{}
		uint64MeetingId, _ := strconv.ParseUint(strId, 10, 64)
		webinarBrandInfoDetail.MeetingId = proto.Uint64(uint64MeetingId) // 会议id
		// 会议介绍
		if strIntro, ok := meetingIdIntroductionMap[strId]; ok {
			webinarBrandInfoDetail.MeetingIntroduction = proto.String(strIntro)
		} else {
			webinarBrandInfoDetail.MeetingIntroduction = proto.String("")
		}
		//日程列表
		if itineraryList, ok := meetingIdItineraryMap[strId]; ok {
			webinarBrandInfoDetail.ItineraryInfoList = itineraryList
		} else {
			webinarBrandInfoDetail.ItineraryInfoList = []*pb.WebinarItineraryInfo{}
		}
		webinarBrandInfoList = append(webinarBrandInfoList, webinarBrandInfoDetail)
	}

	rsp.WebinarBrandInfoList = webinarBrandInfoList

	return 0, nil
}

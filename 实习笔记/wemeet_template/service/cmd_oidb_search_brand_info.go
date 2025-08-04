package service

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"meeting_template/es"
	"meeting_template/util"

	"git.code.oa.com/going/attr"
	"git.code.oa.com/meettrpc/meet_util"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/golang/protobuf/proto"
)

// SearchWebinarBrandInfo ...  oidb协议
func (s *WemeetMeetingTemplateOidbServiceImpl) SearchWebinarBrandInfo(ctx context.Context,
	req *pb.SearchWebinarBrandInfoReq, rsp *pb.SearchWebinarBrandInfoRsp) (err error) {
	attr.AttrAPI(36337025, 1) //[SearchWebinarBrandInfo]请求
	rst, err := HandleSearchWebinarBrandInfo(ctx, req, rsp)
	rsp.ErrorCode = proto.Uint32(rst)
	if err != nil {
		attr.AttrAPI(36337015, 1) //[SearchWebinarBrandInfo]请求失败
		rsp.ErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "SearchWebinarBrandInfo fail req:%v, rsp:%v", req, rsp)
	} else {
		attr.AttrAPI(36337016, 1) //[SearchWebinarBrandInfo]请求成功
		rsp.ErrorMessage = proto.String("ok")
		log.InfoContextf(ctx, "SearchWebinarBrandInfo ok  req:%v, rsp:%v", req, rsp)
	}
	return nil
}

// HandleSearchWebinarBrandInfo ...
func HandleSearchWebinarBrandInfo(ctx context.Context, req *pb.SearchWebinarBrandInfoReq,
	rsp *pb.SearchWebinarBrandInfoRsp) (uint32, error) {
	// page方面的参数校验
	if req.GetPage() <= 0 || req.GetPageSize() > 50 {
		attr.AttrAPI(36337017, 1)
		log.ErrorContextf(ctx, "HandleSearchWebinarBrandInfo page param failed, req:%+v", req)
		return util.ErrPageParam, errors.New("page param error")
	}

	if req.GetSearchType() == 1 || req.GetSearchType() == 2 { //日程搜索
		return HandleSearchByItinerary(ctx, req, rsp)
	}
	if req.GetSearchType() == 3 {
		return HandleSearchByIntroduction(ctx, req, rsp)
	}
	return 0, nil
}

// HandleSearchByItinerary ... 日程名称搜索
func HandleSearchByItinerary(ctx context.Context, req *pb.SearchWebinarBrandInfoReq,
	rsp *pb.SearchWebinarBrandInfoRsp) (uint32, error) {
	if req.GetSearchType() == 0 || req.GetSearchKey() == "" || req.GetSearchTimeType() == 0 {
		attr.AttrAPI(36337018, 1)
		log.ErrorContextf(ctx, "HandleSearchByItineraryName SearchType or SearchKey invalid, req:%+v", req)
		return util.ErrSearchKey, errors.New("SearchType or SearchKey invalid")
	}
	startTs := time.Unix(int64(req.GetStartTime()), 0)
	strStartTime := startTs.Format("2006-01-02 15:04:05")
	endTs := time.Unix(int64(req.GetEndTime()), 0)
	strEndTime := endTs.Format("2006-01-02 15:04:05")

	searchTimeType := req.GetSearchTimeType() //搜索的时间类型

	// 搜索
	total, itineraryList, err := es.FuzzyQueryItineraryInfoFromES(ctx,
		req.GetSearchKey(), int(req.GetPage()), int(req.GetPageSize()),
		req.GetSearchType(), searchTimeType, strStartTime, strEndTime)
	if err != nil {
		attr.AttrAPI(36337019, 1)
		log.ErrorContextf(ctx, "HandleSearchByItinerary failed, req:%+v, err:%+v", req, err)
		return util.ErrFuzzSearch, err
	}
	rsp.Total = proto.Uint32(total)

	wg := sync.WaitGroup{}
	wg.Add(len(itineraryList))
	var data sync.Map
	for i, val := range itineraryList {
		strMeetingId := val.MeetingId
		go func(newCtx context.Context, k int, meetingId string) {
			defer wg.Done()
			defer meet_util.DefPanicFun()

			webinarBrandInfo := &pb.WebinarBrandInfoDetail{}
			uint64MeetingId, _ := strconv.ParseUint(meetingId, 10, 64)
			webinarBrandInfo.MeetingId = proto.Uint64(uint64MeetingId)
			meetingIntroduction := GetMeetingIntroductionFromES(newCtx, meetingId)
			webinarBrandInfo.MeetingIntroduction = proto.String(meetingIntroduction)
			webinarItineraryList, _ := GetMeetingAllItineraryList(newCtx, meetingId)
			webinarBrandInfo.ItineraryInfoList = webinarItineraryList
			data.Store(k, webinarBrandInfo)
		}(trpc.CloneContext(ctx), i, strMeetingId)
	}
	wg.Wait()

	webinarBrandInfoList := []*pb.WebinarBrandInfoDetail{}
	for a := 0; a < len(itineraryList); a++ {
		if value, ok := data.Load(a); ok {
			data, _ := value.(*pb.WebinarBrandInfoDetail)
			if data == nil {
				log.InfoContextf(ctx, "HandleSearchByItinerary value.(*pb.WebinarBrandInfoDetail) nil")
				attr.AttrAPI(36337020, 1)
				continue
			}
			webinarBrandInfoList = append(webinarBrandInfoList, data)
		}
	}

	rsp.WebinarBrandInfoList = webinarBrandInfoList
	return 0, nil
}

// GetMeetingIntroductionFromES ... 从ES中获取某个会议的介绍
func GetMeetingIntroductionFromES(ctx context.Context, meetingId string) string {
	introductions, err := es.AccurateSearchMeetingIntroduction(ctx, meetingId)
	if err != nil {
		attr.AttrAPI(36338314, 1)
		return ""
	}
	if len(introductions) <= 0 {
		return ""
	}
	return introductions[0].MeetingIntroduction
}

// GetMeetingAllItineraryList ... 获取这场会议所有的日程列表
func GetMeetingAllItineraryList(ctx context.Context, meetingId string) ([]*pb.WebinarItineraryInfo, error) {
	webinarItineraryList := []*pb.WebinarItineraryInfo{}
	itineraryList, err := es.AccurateSearchMeetingAllItinerary(ctx, meetingId)
	if err != nil {
		attr.AttrAPI(36337021, 1)
		return webinarItineraryList, err
	}
	for _, val := range itineraryList {
		webinarItineraryInfo := &pb.WebinarItineraryInfo{}
		uint64MeetingId, _ := strconv.ParseUint(val.MeetingId, 10, 64)
		webinarItineraryInfo.MeetingId = proto.Uint64(uint64MeetingId)
		webinarItineraryInfo.ItineraryName = proto.String(val.ItineraryName)
		webinarItineraryInfo.ItineraryIntroduction = proto.String(val.ItineraryIntroduction)
		webinarItineraryList = append(webinarItineraryList, webinarItineraryInfo)
	}
	return webinarItineraryList, nil
}

// HandleSearchByIntroduction ...  搜索会议介绍
func HandleSearchByIntroduction(ctx context.Context, req *pb.SearchWebinarBrandInfoReq,
	rsp *pb.SearchWebinarBrandInfoRsp) (uint32, error) {
	if req.GetSearchKey() == "" || req.GetSearchTimeType() == 0 {
		attr.AttrAPI(36337022, 1)
		log.ErrorContextf(ctx, "HandleSearchByIntroduction searchKey empty, req:%+v", req)
		return util.ErrSearchKey, errors.New("searchKey empty")
	}
	startTs := time.Unix(int64(req.GetStartTime()), 0)
	strStartTime := startTs.Format("2006-01-02 15:04:05")
	endTs := time.Unix(int64(req.GetEndTime()), 0)
	strEndTime := endTs.Format("2006-01-02 15:04:05")

	searchTimeType := req.GetSearchTimeType() //搜索的时间类型

	// 搜索
	total, introductionList, err := es.FuzzyQueryIntroductionFromES(ctx, req.GetSearchKey(), int(req.GetPage()),
		int(req.GetPageSize()), searchTimeType, strStartTime, strEndTime)
	if err != nil {
		attr.AttrAPI(36337023, 1)
		log.ErrorContextf(ctx, "HandleSearchByIntroduction failed, req:%+v, err:%+v", req, err)
		return util.ErrFuzzSearch, err
	}

	rsp.Total = proto.Uint32(total)

	wg := sync.WaitGroup{}
	wg.Add(len(introductionList))
	var data sync.Map
	for i, val := range introductionList {
		strMeetingId := val.MeetingId
		tempIntroduction := val
		go func(newCtx context.Context, k int, meetingId string) {
			defer wg.Done()
			defer meet_util.DefPanicFun()

			webinarBrandInfo := &pb.WebinarBrandInfoDetail{}
			uint64MeetingId, _ := strconv.ParseUint(meetingId, 10, 64)
			webinarBrandInfo.MeetingId = proto.Uint64(uint64MeetingId)
			webinarBrandInfo.MeetingIntroduction = proto.String(tempIntroduction.MeetingIntroduction) //会议介绍这里直接获取了
			webinarItineraryList, _ := GetMeetingAllItineraryList(newCtx, meetingId)
			webinarBrandInfo.ItineraryInfoList = webinarItineraryList
			data.Store(k, webinarBrandInfo)
		}(trpc.CloneContext(ctx), i, strMeetingId)
	}
	wg.Wait()

	webinarBrandInfoList := []*pb.WebinarBrandInfoDetail{}
	for a := 0; a < len(introductionList); a++ {
		if value, ok := data.Load(a); ok {
			data, _ := value.(*pb.WebinarBrandInfoDetail)
			if data == nil {
				attr.AttrAPI(36337024, 1)
				continue
			}
			webinarBrandInfoList = append(webinarBrandInfoList, data)
		}
	}
	rsp.WebinarBrandInfoList = webinarBrandInfoList
	return 0, nil
}
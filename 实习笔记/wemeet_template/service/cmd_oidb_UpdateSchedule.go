package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"meeting_template/config/config_rainbow"
	"meeting_template/es"
	"meeting_template/rpc"
	"meeting_template/util"

	"git.code.oa.com/going/attr"
	"git.code.oa.com/meettrpc/meet_util"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"google.golang.org/protobuf/proto"
)

const (
	SCScheduleName = "sc_meet_schedule_name"
	SCScheduleDesc = "sc_meet_schedule_desc"
)

//UpdateScheduleData ...
func (s *WemeetMeetingTemplateOidbServiceImpl) UpdateScheduleData(ctx context.Context, req *pb.UpdateScheduleReq,
	rsp *pb.UpdateScheduleRsp) (err error) {
	attr.AttrAPI(35914576, 1) //[UpdateScheduleData]请求
	if req.GetAppUserKey() == nil || req.GetUint64MeetingId() == 0 {
		return errors.New("meetingId or appUserKey is empty")
	}
	isCreatorOrSuperAdmin := checkCreatorOrSuperAdmin(ctx, strconv.Itoa(int(req.GetAppUserKey().GetAppId())),
		req.GetAppUserKey().GetAppUserId(), strconv.FormatUint(req.GetUint64MeetingId(), 10))
	if !isCreatorOrSuperAdmin {
		return errors.New("permission denied，current login user is not creator or admin")
	}
	rst, err := HandleUpdateScheduleData(ctx, req, rsp)
	rsp.ErrorCode = proto.Int32(int32(rst))
	if err != nil {
		attr.AttrAPI(35914577, 1) //[HandleUpdateScheduleData]请求失败
		rsp.ErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "HandleUpdateScheduleData fail req:%v, rsp:%v", req, rsp)
	} else {
		attr.AttrAPI(35914578, 1) //[HandleUpdateScheduleData]请求成功
		rsp.ErrorMessage = proto.String("ok")
		log.InfoContextf(ctx, "HandleUpdateScheduleData ok  req:%v, rsp:%v", req, rsp)
	}
	return nil
}

// HandleUpdateScheduleData ...
func HandleUpdateScheduleData(ctx context.Context, req *pb.UpdateScheduleReq,
	rsp *pb.UpdateScheduleRsp) (int32, error) {
	if req.GetOperateType() == 1 { // 新增
		return HandleAddScheduleData(ctx, req, rsp)
	}
	if req.GetOperateType() == 2 { // 修改
		return HandleModifyScheduleData(ctx, req, rsp)
	}
	if req.GetOperateType() == 3 { // 删除
		return HandleDeleteScheduleData(ctx, req, rsp)
	}
	if req.GetOperateType() == 4 { // 预定相同配置的研讨会
		return HandleBatchSaveScheduleData(ctx, req, rsp)
	}
	return 0, nil
}

// HandleAddScheduleData ...
func HandleAddScheduleData(ctx context.Context, req *pb.UpdateScheduleReq, rsp *pb.UpdateScheduleRsp) (int32, error) {
	rst, err := CheckScheduleParam(ctx, req)
	if err != nil {
		attr.AttrAPI(35914584, 1)
		log.ErrorContextf(ctx, "")
		return rst, err
	}
	scheduleInfo := req.GetScheduleList()[0] // CheckScheduleParam中已经校验了slice长度
	incrKey := util.MakeScheduleIncrKey(req.GetUint64MeetingId())
	scheduleId, err := rpc.RDIncrForGetId(ctx, incrKey)
	if err != nil {
		attr.AttrAPI(35914585, 1)
		log.ErrorContextf(ctx, "rpc redis for get scheduleId failed, req:%v, err:%+v", req, err)
		return 3010, err
	}
	scheduleInfo.Uint32Id = proto.Uint32(uint32(scheduleId))
	log.InfoContextf(ctx, "HandleAddScheduleData GetWebinarSchedule Info, meetingId:%+v, scheduleInfo:%+v",
		req.GetUint64MeetingId(), scheduleInfo)
	key := util.MakeScheduleMainKey(req.GetUint64MeetingId())
	subKey := scheduleId
	val, err := proto.Marshal(scheduleInfo)
	err = rpc.RDHSetWebinarInfo(ctx, key, uint64(subKey), string(val))
	if err != nil {
		attr.AttrAPI(35914586, 1)
		log.ErrorContextf(ctx, "rpc redis save schedule info failed. req:%+v, err:%+v", req, err)
		return 3011, err
	}
	rsp.ScheduleId = proto.Uint32(uint32(scheduleId))
	// 写es
	if CanDoEsLogic() {
		go func(newCtx context.Context) {
			defer meet_util.DefPanicFun()
			if len(scheduleInfo.GetScheduleName()) > 0 || len(scheduleInfo.GetScheduleIndroduction()) > 0 {
				log.InfoContextf(newCtx, "enter es logic, meetingId:%+v, scheduleInfo:%+v",
					req.GetUint64MeetingId(), scheduleInfo)
				HandleSaveItinerary(newCtx, req, scheduleInfo, scheduleId)
			}
		}(trpc.CloneContext(ctx))
	}
	return 0, nil
}

// CheckScheduleParam ...
func CheckScheduleParam(ctx context.Context, req *pb.UpdateScheduleReq) (int32, error) {
	if len(req.GetScheduleList()) != 1 {
		attr.AttrAPI(35914587, 1)
		log.ErrorContextf(ctx, "CheckScheduleParam failed, len schedule not 1, req:%+v", req)
		return util.ERRLenIllegal, errors.New("len schedule not 1")
	}
	if req.GetUint64MeetingId() == 0 {
		attr.AttrAPI(35914588, 1)
		log.ErrorContextf(ctx, "CheckScheduleParam failed, meetingId is empty, req:%+v", req)
		return util.ERRMeetingIdIllegal, errors.New("meetingId is empty")
	}
	// 个数校验
	key := util.MakeScheduleMainKey(req.GetUint64MeetingId())
	rst, _ := rpc.RDHLenWebinarInfo(ctx, key)
	cfg := config_rainbow.GetParticipantConfConfig()
	if rst > int64(cfg.ScheduleMaxCount) {
		attr.AttrAPI(35914589, 1)
		return util.ERRCount, errors.New("over 50 limit")
	}
	scheduleInfo := req.GetScheduleList()[0] // 函数开始已经校验了slice长度
	// 长度校验
	if len([]rune(scheduleInfo.GetScheduleName())) > 20 || len([]rune(scheduleInfo.GetScheduleIndroduction())) > 140 {
		attr.AttrAPI(35914590, 1)
		log.ErrorContextf(ctx, "CheckScheduleParam failed, text length over limit, req:%+v", req)
		return util.ERRLength, errors.New("text length over limit")
	}
	// 敏感词校验
	appUser := req.GetAppUserKey()
	nameHit, _ := rpc.CheckHasSensitiveWords(ctx, req.GetUint64MeetingId(), appUser.GetAppId(), appUser.GetAppUserId(),
		scheduleInfo.GetScheduleName(), SCScheduleName)
	if nameHit {
		attr.AttrAPI(35914591, 1)
		log.ErrorContextf(ctx, "schedule strName hit sensitive, meetingId:%+v, strName:%+v",
			req.GetUint64MeetingId(), scheduleInfo.GetScheduleName())
		return util.ERRNameSensitive, errors.New("schedule strName hit sensitive")
	}
	descHit, _ := rpc.CheckHasSensitiveWords(ctx, req.GetUint64MeetingId(), appUser.GetAppId(), appUser.GetAppUserId(),
		scheduleInfo.GetScheduleIndroduction(), SCScheduleDesc)
	if descHit {
		attr.AttrAPI(35914592, 1)
		log.ErrorContextf(ctx, "schedule introduction hit sensitive, meetingId:%+v, introduction:%+v",
			req.GetUint64MeetingId(), scheduleInfo.GetScheduleIndroduction())
		return util.ERRDescSensitive, errors.New("schedule introduction hit sensitive")
	}
	return 0, nil
}

// HandleModifyScheduleData ...
func HandleModifyScheduleData(ctx context.Context, req *pb.UpdateScheduleReq,
	rsp *pb.UpdateScheduleRsp) (int32, error) {
	rst, err := CheckScheduleParam(ctx, req)
	if err != nil {
		attr.AttrAPI(35914607, 1)
		log.ErrorContextf(ctx, "HandleModifyScheduleData CheckScheduleParam failed. req:%+v, err:%+v", req, err)
		return rst, err
	}
	scheduleInfo := req.GetScheduleList()[0] // CheckScheduleParam中已经校验了slice长度
	log.InfoContextf(ctx, "HandleModifyParticipant GetScheduleInfo, meetingId:%+v, scheduleInfo:%+v",
		req.GetUint64MeetingId(), scheduleInfo)
	key := util.MakeScheduleMainKey(req.GetUint64MeetingId())
	subKey := scheduleInfo.GetUint32Id()
	val, err := proto.Marshal(scheduleInfo)
	err = rpc.RDHSetWebinarInfo(ctx, key, uint64(subKey), string(val))
	if err != nil {
		attr.AttrAPI(35914608, 1)
		log.ErrorContextf(ctx, "rpc redis modify schedule info failed. req:%+v, err:%+v", req, err)
		return 3011, err
	}
	rsp.ScheduleId = proto.Uint32(subKey)
	// 写es
	if CanDoEsLogic() {
		go func(newCtx context.Context) {
			defer meet_util.DefPanicFun()
			if len(scheduleInfo.GetScheduleName()) > 0 || len(scheduleInfo.GetScheduleIndroduction()) > 0 {
				log.InfoContextf(newCtx, "enter es logic, meetingId:%+v, scheduleInfo:%+v",
					req.GetUint64MeetingId(), scheduleInfo)
				HandleModifyItinerary(newCtx, req, scheduleInfo, int64(subKey))
			}
		}(trpc.CloneContext(ctx))
	}
	return 0, nil
}

// HandleDeleteScheduleData ...
func HandleDeleteScheduleData(ctx context.Context, req *pb.UpdateScheduleReq,
	rsp *pb.UpdateScheduleRsp) (int32, error) {
	fields := []string{}
	schedules := req.GetScheduleList()
	for i := 0; i < len(schedules); i++ {
		tempSchedule := schedules[i]
		fields = append(fields, fmt.Sprint(tempSchedule.GetUint32Id()))
	}
	key := util.MakeScheduleMainKey(req.GetUint64MeetingId())
	err := rpc.RDHDelWebinarInfo(ctx, key, fields)
	if err != nil {
		attr.AttrAPI(35914609, 1)
		log.ErrorContextf(ctx, "rpc redis del schedule failed. req:%+v, err:%+v", req, err)
		return 3012, err
	}
	//删除es
	if CanDoEsLogic() {
		go func(newCtx context.Context) {
			defer meet_util.DefPanicFun()
			if len(schedules) > 0 {
				scheduleId := schedules[0].GetUint32Id()
				docId := BuildItineraryESDocId(req.GetUint64MeetingId(), int64(scheduleId))
				es.DelItineraryToES(newCtx, docId)
			}
		}(trpc.CloneContext(ctx))
	}
	log.InfoContextf(ctx, "HandleDeleteScheduleData succ. req:%+v", req)
	return 0, nil
}

// HandleBatchSaveScheduleData ...
func HandleBatchSaveScheduleData(ctx context.Context, req *pb.UpdateScheduleReq,
	rsp *pb.UpdateScheduleRsp) (int32, error) {
	schedules := req.GetScheduleList()
	tempMap := make(map[uint64]string)
	var maxScheduleId uint32 = 0
	for i := 0; i < len(schedules); i++ {
		tempSchedule := schedules[i]
		subKey := tempSchedule.GetUint32Id()
		maxScheduleId = util.Max(maxScheduleId, subKey)
		subVal, _ := proto.Marshal(tempSchedule)
		tempMap[uint64(subKey)] = string(subVal)
	}
	key := util.MakeScheduleMainKey(req.GetUint64MeetingId())
	err := rpc.RDHMSETWebinarInfo(ctx, key, tempMap)
	if err != nil {
		attr.AttrAPI(35914610, 1)
		log.ErrorContextf(ctx, "HandleBatchSaveScheduleData rpc redis batch save schedule failed."+
			"req:%+v, err:%+v", req, err)
		return 3013, err
	}
	incrKey := util.MakeScheduleIncrKey(req.GetUint64MeetingId())
	err = rpc.RDSetIncrValue(ctx, incrKey, maxScheduleId)
	if err != nil {
		attr.AttrAPI(35914611, 1)
		log.ErrorContextf(ctx, "HandleBatchSaveScheduleData set maxScheduleId failed. req:%+v, err:%+v", req, err)
		return 3013, err
	}
	return 0, nil
}

// BuildItineraryESDocId ...
func BuildItineraryESDocId(meetingId uint64, scheduleId int64) string {
	id := fmt.Sprintf("%v_%v", meetingId, scheduleId)
	return id
}

// HandleSaveItinerary ...
func HandleSaveItinerary(ctx context.Context, req *pb.UpdateScheduleReq, scheduleInfo *pb.WebinarSchedule,
	scheduleId int64) error {
	itineraryIntroduction := &es.Itinerary{
		MeetingId:             fmt.Sprint(req.GetUint64MeetingId()),
		AppUid:                req.GetAppUserKey().GetAppUserId(),
		AppId:                 fmt.Sprint(req.GetAppUserKey().GetAppId()),
		ItineraryName:         scheduleInfo.GetScheduleName(),
		ItineraryIntroduction: scheduleInfo.GetScheduleIndroduction(),
	}
	//获取 meetingInfo
	meetingInfo, _, _, err := rpc.GetMeetingInfo(ctx, req.GetUint64MeetingId())
	if err != nil {
		attr.AttrAPI(36337039, 1)
		log.ErrorContextf(ctx, "HandleSaveItinerary failed, meetingId:%+v, err:%+v", req.GetUint64MeetingId(), err)
		return nil
	}
	docId := BuildItineraryESDocId(req.GetUint64MeetingId(), scheduleId)
	err = es.SaveItineraryToES(ctx, itineraryIntroduction, docId, meetingInfo)
	if err != nil {
		attr.AttrAPI(36337040, 1)
		log.ErrorContextf(ctx, "HandleSaveItinerary save es failed. meetingId:%+v, err:%+v",
			req.GetUint64MeetingId(), err)
		return err
	}
	log.InfoContextf(ctx, "HandleSaveItinerary succ, meetingId:%+v", req.GetUint64MeetingId())
	return nil
}

// HandleModifyItinerary ...
func HandleModifyItinerary(ctx context.Context, req *pb.UpdateScheduleReq, scheduleInfo *pb.WebinarSchedule,
	scheduleId int64) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	itineraryIntroduction := &es.Itinerary{
		MeetingId:             fmt.Sprint(req.GetUint64MeetingId()),
		AppUid:                req.GetAppUserKey().GetAppUserId(),
		AppId:                 fmt.Sprint(req.GetAppUserKey().GetAppId()),
		ItineraryName:         scheduleInfo.GetScheduleName(),
		ItineraryIntroduction: scheduleInfo.GetScheduleIndroduction(),
		UpdateTime:            now, // 修改时间
	}
	docId := BuildItineraryESDocId(req.GetUint64MeetingId(), scheduleId)
	err := es.ModifyItineraryToES(ctx, itineraryIntroduction, docId)
	if err != nil {
		log.ErrorContextf(ctx, "HandleModifyItinerary update es failed. meetingId:%+v, err:%+v",
			req.GetUint64MeetingId(), err)
		return err
	}
	log.InfoContextf(ctx, "HandleModifyItinerary succ, meetingId:%+v", req.GetUint64MeetingId())
	return nil
}

// CanDoEsLogic ...
func CanDoEsLogic() bool {
	esConf := config_rainbow.GetEsConfConfig()
	if esConf.EsSwitch == "open" {
		return true
	}
	return false
}

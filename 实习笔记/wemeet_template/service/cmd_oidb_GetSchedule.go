package service

import (
	"context"
	"sort"

	"meeting_template/rpc"
	"meeting_template/util"

	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"google.golang.org/protobuf/proto"
)

// GetScheduleData ...
func (s *WemeetMeetingTemplateOidbServiceImpl) GetScheduleData(ctx context.Context, req *pb.GetScheduleReq,
	rsp *pb.GetScheduleRsp) (err error) {
	attr.AttrAPI(35914626, 1) //[GetScheduleData]请求
	rst, err := HandleGetScheduleData(ctx, req, rsp)
	rsp.ErrorCode = proto.Int32(int32(rst))
	if err != nil {
		attr.AttrAPI(35914627, 1) //[HandleGetScheduleData]请求失败
		rsp.ErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "HandleGetScheduleData fail req:%v, rsp:%v", req, rsp)
	} else {
		attr.AttrAPI(35914628, 1) //[HandleGetScheduleData]请求成功
		rsp.ErrorMessage = proto.String("ok")
		log.InfoContextf(ctx, "HandleGetScheduleData ok  req:%v, rsp:%v", req, rsp)
	}
	return nil
}

// HandleGetScheduleData ...
func HandleGetScheduleData(ctx context.Context, req *pb.GetScheduleReq, rsp *pb.GetScheduleRsp) (int32, error) {
	key := util.MakeScheduleMainKey(req.GetUint64MeetingId())
	allSchedule, err := rpc.RDHValsWebinarInfo(ctx, key) //NOCA:RedisHashSlowCmd(EnsureSafe@yucachen)
	if err != nil {
		attr.AttrAPI(35914629, 1)
		log.ErrorContextf(ctx, "HandleGetScheduleData rpc redis get all schedule data failed. req:%+v, err:%+v",
			req, err)
		return 3020, err
	}
	scheduleList := []*pb.WebinarSchedule{}
	for _, schedule := range allSchedule {
		tempSchedule := &pb.WebinarSchedule{}
		err := proto.Unmarshal([]byte(schedule), tempSchedule)
		if err != nil {
			attr.AttrAPI(35915500, 1)
			continue
		}
		scheduleList = append(scheduleList, tempSchedule)
	}
	// 根据Id排序
	sort.Sort(ScheduleSlice(scheduleList))
	rsp.ScheduleList = scheduleList
	return 0, nil
}

type ScheduleSlice []*pb.WebinarSchedule

// Len ...
func (s ScheduleSlice) Len() int {
	return len(s)
}

// Less ...
func (s ScheduleSlice) Less(i, j int) bool {
	return s[i].GetUint32Id() < s[j].GetUint32Id() // ID从小到大排序
}

// Swap ...
func (s ScheduleSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
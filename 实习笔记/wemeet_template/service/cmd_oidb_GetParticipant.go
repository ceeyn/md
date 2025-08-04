package service

import (
	"context"
	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"google.golang.org/protobuf/proto"
	"meeting_template/rpc"
	"meeting_template/util"
	"sort"
)

// GetParticipant ...
func (s *WemeetMeetingTemplateOidbServiceImpl) GetParticipant(ctx context.Context, req *pb.GetParticipantReq,
	rsp *pb.GetParticipantRsp) (err error) {
	attr.AttrAPI(35914622, 1) //[GetParticipant]请求
	rst, err := HandleGetParticipant(ctx, req, rsp)
	rsp.Code = proto.Int32(int32(rst))
	if err != nil {
		attr.AttrAPI(35914623, 1) //[HandleGetParticipant]请求失败
		rsp.Message = proto.String(err.Error())
		log.ErrorContextf(ctx, "HandleGetParticipant fail req:%v, rsp:%v", req, rsp)
	} else {
		attr.AttrAPI(35914624, 1) //[HandleGetParticipant]请求成功
		rsp.Message = proto.String("ok")
		log.InfoContextf(ctx, "HandleGetParticipant ok  req:%v, rsp:%v", req, rsp)
	}
	return nil
}

// HandleGetParticipant ...
func HandleGetParticipant(ctx context.Context, req *pb.GetParticipantReq, rsp *pb.GetParticipantRsp) (int32, error) {
	key := util.MakeParticipantMainKey(req.GetMeetingId())
	allParticipant, err := rpc.RDHValsWebinarInfo(ctx, key)  //NOCA:RedisHashSlowCmd(EnsureSafe@yucachen)
	if err != nil {
		attr.AttrAPI(35914625,1)
		log.ErrorContextf(ctx,"HandleGetParticipant rpc redis get all participants failed. req:%+v, err:%+v", req, err)
		return 3020, err
	}
	participantList := []*pb.Participant{}
	for _, participant := range allParticipant{
		tempParticipant := &pb.Participant{}
		err := proto.Unmarshal([]byte(participant), tempParticipant)
		if err != nil {
			continue
		}
		participantList = append(participantList, tempParticipant)
	}

	// 获取cosId对应的URL
	cosIdUrlMap := make(map[string]string)
	imageType := util.MiddleImageType // 前端未赋值时，默认获取中图
	cosIds := []string{}
	for _, participant := range participantList {
		cosIds = append(cosIds, participant.GetCosId())
	}
	GetMapDataFromCosIds(ctx, imageType, util.DownloadUseCdn, cosIds, cosIdUrlMap)
	for _, participant := range participantList {
		if val, ok := cosIdUrlMap[participant.GetCosId()]; ok {
			participant.AvatarUrl = proto.String(val)
		}else {
			participant.AvatarUrl = proto.String("")
		}
	}
	log.InfoContextf(ctx, "HandleGetParticipant detail Info, req:%+v, participantList:%+v", req, participantList)

	// 根据userId排序
	sort.Sort(ParticipantSlice(participantList))
	rsp.ParticipantList = participantList
	return 0, nil
}


type ParticipantSlice []*pb.Participant

// Len ...
func (s ParticipantSlice) Len() int {
	return len(s)
}

//Less ...
func (s ParticipantSlice) Less(i, j int) bool {
	return s[i].GetId() < s[j].GetId()      // userID从小到大排序
}

// Swap ...
func (s ParticipantSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
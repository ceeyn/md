package service

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"meeting_template/rpc"
	"meeting_template/util"
	
	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	safePb "git.code.oa.com/trpcprotocol/wemeet/wemeet_safe_gateway"
	"google.golang.org/protobuf/proto"
)

const (
	CallBackAuditPass   = 0
	CallBackAuditFailed = 1
	AvatarPass          = 1
	AvatarNotPass       = 2
)

// AvatarSafeCallback ...
func (s *WemeetMeetingTemplateOidbServiceImpl) AvatarSafeCallback(ctx context.Context, req *safePb.GetImageCallbackReq,
	rsp *safePb.GetImageCallbackRsp) (err error) {
	attr.AttrAPI(35914630, 1)
	rst, err := HandleAvatarSafeCallback(ctx, req)
	rsp.ErrorCode = proto.Int32(rst)
	if err != nil {
		attr.AttrAPI(35914631, 1)
		rsp.ErrorMsg = proto.String(err.Error())
		log.ErrorContextf(ctx, "AvatarSafeCallback fail req:%+v, rsp:%+v", req, rsp)
	} else {
		attr.AttrAPI(35914632, 1)
		rsp.ErrorMsg = proto.String("ok")
		log.InfoContextf(ctx, "AvatarSafeCallback ok  req:%+v, rsp:%+v", req, rsp)
	}
	return nil
}

// HandleAvatarSafeCallback ...
func HandleAvatarSafeCallback(ctx context.Context, req *safePb.GetImageCallbackReq) (int32, error) {
	resultCode := req.GetResultCode()
	if resultCode != CallBackAuditPass && resultCode != CallBackAuditFailed {
		attr.AttrAPI(35915450, 1)
		return 1, errors.New("HandleAvatarSafeCallback code illegal")
	}
	meetingId := req.GetMeet().GetMeetingId()
	uuid := req.GetSer().GetTraceId()
	if uuid == "" {
		attr.AttrAPI(35914634, 1)
		log.ErrorContextf(ctx, "HandleAvatarSafeCallback trace id empty")
		return 2, errors.New("trace id empty")
	}
	infos := strings.Split(uuid, "&")
	if len(infos) != 2 {
		attr.AttrAPI(35914635, 1)
		log.ErrorContextf(ctx, "HandleAvatarSafeCallback trace id format error, trace_id:%v", uuid)
		return 3, errors.New("traceId format error")
	}
	userId := infos[0]
	cosId := infos[1]
	// 更新存储中的数据
	if err := UpdateParticipantAuditStatus(ctx, meetingId, userId, cosId, resultCode); err != nil {
		attr.AttrAPI(35914636, 1)
		log.ErrorContextf(ctx, "HandleAvatarSafeCallback UpdateParticipantAuditStatus failed. req:%+v", req)
		return 4, err
	}
	log.InfoContextf(ctx,"HandleAvatarSafeCallback sucess. req:%+v", req)
	return 0, nil
}

// UpdateParticipantAuditStatus 更新嘉宾头像的审核状态
func UpdateParticipantAuditStatus(ctx context.Context, meetingId uint64, userId string, cosId string,
	resultCode uint32) error {
	key := util.MakeParticipantMainKey(meetingId)
	val, err := rpc.RDHGetParticipantInfo(ctx, key, userId)
	if err != nil {
		attr.AttrAPI(35914637, 1)
		log.ErrorContextf(ctx, "UpdateParticipantAuditStatus rpc get ParticipantInfo failed. err:%+v", err)
		return err
	}
	participantInfo := &pb.Participant{}
	err = proto.Unmarshal([]byte(val), participantInfo)
	if err != nil {
		attr.AttrAPI(35914638, 1)
		log.ErrorContextf(ctx, "participantInfo.Unmarshal failed. err:%+v", err)
		return err
	}
	if cosId != participantInfo.GetCosId() {
		attr.AttrAPI(35914639, 1)
		log.ErrorContextf(ctx, "UpdateParticipantAuditStatus cosId dont same, participant:%+v, req:%+v",
			participantInfo, cosId)
		return errors.New("cosId dont same")
	}
	if resultCode == 0 { // 审核通过
		participantInfo.AvatarStatus = proto.Uint32(AvatarPass) // 状态通过
	}
	if resultCode == 1 { // 审核不通过
		participantInfo.AvatarStatus = proto.Uint32(AvatarNotPass) // 审核不通过
		participantInfo.CosId = proto.String("")                   // cosId置为空
		participantInfo.FileName = proto.String("")                // 文件名置为空
	}
	strParticipant, err := proto.Marshal(participantInfo)
	if err != nil {
		attr.AttrAPI(35914640, 1)
		log.ErrorContextf(ctx, "Marshal participantInfo failed")
		return errors.New("Marshal participantInfo failed")
	}
	subKey, _ := strconv.ParseUint(userId, 10, 64)
	err = rpc.RDHSetWebinarInfo(ctx, key, subKey, string(strParticipant))
	return err
}

package service

import (
	"context"
	"errors"
	"fmt"
	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"google.golang.org/protobuf/proto"
	"meeting_template/config/config_rainbow"
	"meeting_template/rpc"
	"meeting_template/util"
	"strconv"
)

const (
	SCParticipantName              = "sc_participant_name"
	SCParticipantProfile           = "sc_participant_profile"
	SCParticipantAvatar            = "sc_participant_avatar"
	SCParticipantAvatarAction      = "get_participant_avatar_check"
	SCParticipantAvatarCallBackUri = "/wemeet-template/AvatarSafeCallback"
)

// UpdateParticipant ...
func (s *WemeetMeetingTemplateOidbServiceImpl) UpdateParticipant(ctx context.Context, req *pb.UpdateParticipantReq,
	rsp *pb.UpdateParticipantRsp) (err error) {
	attr.AttrAPI(35914559, 1) //[UpdateParticipant]请求
	if req.GetAppUserKey() == nil || req.GetUint64MeetingId() == 0 {
		return errors.New("meetingId or appUserKey is empty")
	}
	isCreatorOrSuperAdmin := checkCreatorOrSuperAdmin(ctx, strconv.Itoa(int(req.GetAppUserKey().GetAppId())),
		req.GetAppUserKey().GetAppUserId(), strconv.FormatUint(req.GetUint64MeetingId(), 10))
	if !isCreatorOrSuperAdmin {
		return errors.New("permission denied，current login user is not creator or admin")
	}
	rst, err := HandleUpdateParticipant(ctx, req, rsp)
	rsp.ErrorCode = proto.Int32(int32(rst))
	if err != nil {
		attr.AttrAPI(35914560, 1) //[HandleUpdateParticipant]请求失败
		rsp.ErrorMessage = proto.String(err.Error())
		log.ErrorContextf(ctx, "HandleUpdateParticipant fail req:%v, rsp:%v", req, rsp)
	} else {
		attr.AttrAPI(35914561, 1) //[HandleUpdateParticipant]请求成功
		rsp.ErrorMessage = proto.String("ok")
		log.InfoContextf(ctx, "HandleUpdateParticipant ok  req:%v, rsp:%v", req, rsp)
	}
	return nil
}

// HandleUpdateParticipant ...
func HandleUpdateParticipant(ctx context.Context, req *pb.UpdateParticipantReq,
	rsp *pb.UpdateParticipantRsp) (int32, error) {
	if req.GetOperateType() == 1 { // 新增
		return HandleAddParticipant(ctx, req, rsp)
	}
	if req.GetOperateType() == 2 { // 修改
		return HandleModifyParticipant(ctx, req, rsp)
	}
	if req.GetOperateType() == 3 { // 删除
		return HandleDeleteParticipant(ctx, req, rsp)
	}
	if req.GetOperateType() == 4 { // 预定相同配置的研讨会
		return HandleBatchSaveParticipant(ctx, req, rsp)
	}
	return 0, nil
}

// HandleAddParticipant ...
func HandleAddParticipant(ctx context.Context, req *pb.UpdateParticipantReq,
	rsp *pb.UpdateParticipantRsp) (int32, error) {
	rst, err := CheckParticipantParam(ctx, req)
	if err != nil {
		attr.AttrAPI(35914562, 1)
		log.ErrorContextf(ctx, "HandleAddParticipant CheckParticipantParam failed, req:%+v", req)
		return rst, err
	}
	participant := req.GetParticipantList()[0]
	participant.AvatarStatus = proto.Uint32(util.NotAudit) // 添加的时候初始未审核
	if participant.GetCosId() == "" {
		participant.AvatarStatus = proto.Uint32(util.PassAudit) // 头像为空添加的时候审核通过
	}
	incrKey := util.MakeParticipantIncrKey(req.GetUint64MeetingId())
	userId, err := rpc.RDIncrForGetId(ctx, incrKey)
	if err != nil {
		attr.AttrAPI(35914563, 1)
		log.ErrorContextf(ctx, "rpc redis for get userId failed, req:%v, err:%+v", req, err)
		return 3010, err
	}
	participant.Id = proto.Uint32(uint32(userId))
	key := util.MakeParticipantMainKey(req.GetUint64MeetingId())
	subKey := userId
	val, err := proto.Marshal(participant)
	err = rpc.RDHSetWebinarInfo(ctx, key, uint64(subKey), string(val))
	if err != nil {
		attr.AttrAPI(35914564, 1)
		log.ErrorContextf(ctx, "rpc redis save participant info failed. req:%+v, err:%+v", req, err)
		return 3011, err
	}
	rsp.ParticipantId = proto.Uint32(uint32(userId))
	if participant.GetCosId() != "" {
		go func(newCtx context.Context) {
			ParticipantAvatarAudit(newCtx, req, userId)
		}(trpc.CloneContext(ctx))
	}
	return 0, nil
}

// CheckParticipantParam ...
func CheckParticipantParam(ctx context.Context, req *pb.UpdateParticipantReq) (int32, error) {
	if len(req.GetParticipantList()) != 1 {
		attr.AttrAPI(35914565, 1)
		log.ErrorContextf(ctx, "CheckAddParticipantParam failed, len participant not 1, req:%+v", req)
		return util.ERRLenIllegal, errors.New("len participant not 1")
	}
	if req.GetUint64MeetingId() == 0 {
		attr.AttrAPI(35914566, 1)
		log.ErrorContextf(ctx, "CheckAddParticipantParam failed, meetingId is empty, req:%+v", req)
		return util.ERRMeetingIdIllegal, errors.New("meetingId is empty")
	}
	// 个数校验
	key := util.MakeParticipantMainKey(req.GetUint64MeetingId())
	rst, _ := rpc.RDHLenWebinarInfo(ctx, key)
	cfg := config_rainbow.GetParticipantConfConfig()
	if rst > int64(cfg.ParticipantMaxCount) {
		attr.AttrAPI(35914567, 1)
		return util.ERRCount, errors.New("over 200 limit")
	}
	participant := req.GetParticipantList()[0]
	// cosId合法性校验
	if !util.IsValidCosId(ctx, participant.GetCosId()) {
		log.ErrorContextf(ctx, "CheckAddParticipantParam failed, cosId illegal, req:%+v", req)
		return util.ERRCosIdIllegal, errors.New("cosId illegal")
	}
	// 长度校验
	if len([]rune(participant.GetName())) > 30 || len([]rune(participant.GetGuestIntroduction())) > 80 {
		log.ErrorContextf(ctx, "CheckAddParticipantParam failed, text length over limit, req:%+v", req)
		return util.ERRLength, errors.New("text length over limit")
	}
	// 敏感词校验
	appUser := req.GetAppUserKey()
	nameHit, _ := rpc.CheckHasSensitiveWords(ctx, req.GetUint64MeetingId(), appUser.GetAppId(), appUser.GetAppUserId(),
		participant.GetName(), SCParticipantName)
	if nameHit {
		attr.AttrAPI(35914568, 1)
		log.ErrorContextf(ctx, "participant strName hit sensitive, meetingId:%+v, strName:%+v",
			req.GetUint64MeetingId(), participant.GetName())
		return util.ERRNameSensitive, errors.New("participant strName hit sensitive")
	}
	descHit, _ := rpc.CheckHasSensitiveWords(ctx, req.GetUint64MeetingId(), appUser.GetAppId(), appUser.GetAppUserId(),
		participant.GetGuestIntroduction(), SCParticipantProfile)
	if descHit {
		attr.AttrAPI(35914569, 1)
		log.ErrorContextf(ctx, "participant introduction hit sensitive, meetingId:%+v, introduction:%+v",
			req.GetUint64MeetingId(), participant.GetGuestIntroduction())
		return util.ERRDescSensitive, errors.New("participant introduction hit sensitive")
	}
	return 0, nil
}

// HandleModifyParticipant ...
func HandleModifyParticipant(ctx context.Context, req *pb.UpdateParticipantReq,
	rsp *pb.UpdateParticipantRsp) (int32, error) {
	rst, err := CheckParticipantParam(ctx, req)
	if err != nil {
		attr.AttrAPI(35914571, 1)
		log.ErrorContextf(ctx, "HandleModifyParticipant CheckParticipantParam failed. req:%+v, err:%+v", req, err)
		return rst, err
	}
	newParticipant := req.GetParticipantList()[0]
	key := util.MakeParticipantMainKey(req.GetUint64MeetingId())
	subKey := newParticipant.GetId()
	val, err := rpc.RDHGetParticipantInfo(ctx, key, fmt.Sprint(subKey))
	oldParticipant := &pb.Participant{}
	err = proto.Unmarshal([]byte(val), oldParticipant)
	if err != nil {
		attr.AttrAPI(35915509, 1)
		return util.ERRUnmarshal, errors.New("Unmarshal participant failed")
	}
	needCheckAvatar := true
	if oldParticipant.GetCosId() == newParticipant.GetCosId() {
		newParticipant.AvatarStatus = proto.Uint32(oldParticipant.GetAvatarStatus())
		if oldParticipant.GetAvatarStatus() != 0 {
			needCheckAvatar = false
		}
	} else {
		newParticipant.AvatarStatus = proto.Uint32(0) //修改的是头像的话，先把审核状态置为未审核
	}
	if newParticipant.GetCosId() == "" {
		newParticipant.AvatarStatus = proto.Uint32(1) //修改的时候头像为空，审核状态为已审核
		needCheckAvatar = false
	}
	strNewParticipant, err := proto.Marshal(newParticipant)
	err = rpc.RDHSetWebinarInfo(ctx, key, uint64(subKey), string(strNewParticipant))
	if err != nil {
		attr.AttrAPI(35914572, 1)
		log.ErrorContextf(ctx, "rpc redis modify participant info failed. req:%+v, err:%+v", req, err)
		return 3011, err
	}
	rsp.ParticipantId = proto.Uint32(subKey)
	if needCheckAvatar {
		ParticipantAvatarAudit(ctx, req, int64(subKey))
	}
	return 0, nil
}

// HandleDeleteParticipant ...
func HandleDeleteParticipant(ctx context.Context, req *pb.UpdateParticipantReq,
	rsp *pb.UpdateParticipantRsp) (int32, error) {
	fields := []string{}
	participants := req.GetParticipantList()
	for i := 0; i < len(participants); i++ {
		tempParticipant := participants[i]
		fields = append(fields, fmt.Sprint(tempParticipant.GetId()))
	}
	key := util.MakeParticipantMainKey(req.GetUint64MeetingId())
	err := rpc.RDHDelWebinarInfo(ctx, key, fields)
	if err != nil {
		attr.AttrAPI(35914573, 1)
		log.ErrorContextf(ctx, "rpc redis del participant failed. req:%+v, err:%+v", req, err)
		return 3012, err
	}
	log.InfoContextf(ctx, "HandleDeleteParticipant succ. req:%+v", req)
	return 0, nil
}

// HandleBatchSaveParticipant ...
func HandleBatchSaveParticipant(ctx context.Context, req *pb.UpdateParticipantReq,
	rsp *pb.UpdateParticipantRsp) (int32, error) {
	participantList := req.GetParticipantList()
	tempInsertMap := make(map[uint64]string)
	var maxUserId uint32 = 0
	for i := 0; i < len(participantList); i++ {
		tempParticipant := participantList[i]
		subKey := tempParticipant.GetId()
		tempParticipant.AvatarStatus = proto.Uint32(1) // 审核通过
		maxUserId = util.Max(maxUserId, subKey)
		subVal, _ := proto.Marshal(tempParticipant)
		tempInsertMap[uint64(subKey)] = string(subVal)
	}
	key := util.MakeParticipantMainKey(req.GetUint64MeetingId())
	err := rpc.RDHMSETWebinarInfo(ctx, key, tempInsertMap)
	if err != nil {
		attr.AttrAPI(35914574, 1)
		log.ErrorContextf(ctx, "HandleBatchSaveParticipant rpc redis batch save participant failed."+
			"req:%+v, err:%+v", req, err)
		return 3013, err
	}
	incrKey := util.MakeParticipantIncrKey(req.GetUint64MeetingId())
	err = rpc.RDSetIncrValue(ctx, incrKey, maxUserId)
	if err != nil {
		attr.AttrAPI(35914575, 1)
		log.ErrorContextf(ctx, "HandleBatchSaveParticipant set maxUserId failed. req:%+v, err:%+v", req, err)
		return 3014, err
	}
	return 0, nil
}

// ParticipantAvatarAudit ... 嘉宾头像送审
func ParticipantAvatarAudit(ctx context.Context, req *pb.UpdateParticipantReq, userId int64) {
	// 图片送审
	cosId := req.GetParticipantList()[0].GetCosId()
	auditUuid := fmt.Sprint(userId) + "&" + cosId // uuid由userId + "&" + cos_id构成
	url := GetDownloadUrl(ctx, util.RawImageType, cosId, util.DownloadUseCdn)
	callBackSwitch := config_rainbow.GetCallBackConf().SCParticipantAvatarAction
	log.InfoContextf(ctx, "rpcImgSafetyAudit SCParticipantAvatarAction: %v", callBackSwitch)
	if callBackSwitch {
		safeReq := &rpc.SafetyAuditReq{
			Uuid:              auditUuid,
			AppId:             req.GetAppUserKey().GetAppId(),
			AppUid:            req.GetAppUserKey().GetAppUserId(),
			MeetingId:         req.GetUint64MeetingId(),
			Url:               url,
			Scenes:            SCParticipantAvatar,
			Action:            SCParticipantAvatarAction,
			CallbackTarget:    CallbackTarget,
			CallbackUri:       SCParticipantAvatarCallBackUri,
			CallbackEnvName:   util.GetEnvName(),
			CallbackNameSpace: util.GetNameSpace(),
			CallbackProtoType: CallbackProtoType,
		}
		rpc.ImgSafetyAuditV2(ctx, safeReq)
	} else {
		rpc.ImgSafetyAudit(ctx, auditUuid, req.GetAppUserKey().GetAppId(), req.GetAppUserKey().GetAppUserId(),
			req.GetUint64MeetingId(), url, SCParticipantAvatar)
	}
}

package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"meeting_template/rpc"
	"meeting_template/util"

	"git.code.oa.com/meettrpc/meet_util"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/http"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	cachePb "git.code.oa.com/trpcprotocol/wemeet/common_meeting_cache"
	errpb "git.code.oa.com/trpcprotocol/wemeet/common_xcast_meeting_error_code"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	safePb "git.code.oa.com/trpcprotocol/wemeet/wemeet_safe_gateway"
	"github.com/golang/protobuf/proto"
	"meeting_template/config/config_rainbow"
	bg "meeting_template/material_control/background"
	ctlrpc "meeting_template/material_control/rpc"
)

const (
	GateWayImgBackgroundScenes = "sc_webinar_background_img" // 背景图片送审场景
	AppFromWemeetLayoutCenter  = "wemeet_layout_center"
	GateWayImgBackgroundActon  = "get_webinar_background_img_check"
	ImgBackgroundCallBackUri   = "/wemeet-template/VirtualBackgroundSafeCallback"
)

const (
	SetBackgroundOperateTypeAdd    = 1
	SetBackgroundOperateTypeUpdate = 2
	SetBackgroundOperateTypeDelete = 3
)

//SetVirtualBackground 设置虚拟背景
func (s *WemeetMeetingTemplateOidbServiceImpl) SetVirtualBackground(ctx context.Context,
	req *pb.SetVirtualBackgroundReq, rsp *pb.SetVirtualBackgroundRsp) error {
	metrics.IncrCounter("SetVirtualBackground Total", 1)
	start := time.Now()

	rst, err := handleSetVirtualBackground(ctx, req, rsp)
	if err != nil {
		metrics.IncrCounter("SetVirtualBackground fail", 1)
		rsp.ErrorCode = proto.Int32(rst)
		rsp.ErrorMessage = proto.String(err.Error())
	} else {
		metrics.IncrCounter("SetVirtualBackground success", 1)
		rsp.ErrorCode = proto.Int32(0)
		rsp.ErrorMessage = proto.String("ok")
	}
	log.InfoContextf(ctx, "SetVirtualBackground, cost:%v, rst:%v, err:%+v, req:%+v, rsp:%+v",
		time.Since(start), rst, err, req, rsp)
	return nil
}

// QueryVirtualBackgroundList 查询虚拟背景列表
func (s *WemeetMeetingTemplateOidbServiceImpl) QueryVirtualBackgroundList(ctx context.Context,
	req *pb.QueryVirtualBackgroundListReq, rsp *pb.QueryVirtualBackgroundListRsp) error {
	metrics.IncrCounter("QueryVirtualBackgroundList Total", 1)
	start := time.Now()

	rst, err := handleQueryVirtualBackgroundList(ctx, req, rsp)
	if err != nil {
		metrics.IncrCounter("QueryVirtualBackgroundList fail", 1)
		rsp.ErrorCode = proto.Int32(rst)
		rsp.ErrorMessage = proto.String(err.Error())
	} else {
		metrics.IncrCounter("QueryVirtualBackgroundList Succ", 1)
		rsp.ErrorCode = proto.Int32(0)
		rsp.ErrorMessage = proto.String("ok")
	}
	log.InfoContextf(ctx, "QueryVirtualBackgroundList, cost:%v, rst:%v, err:%+v, req:%+v, rsp:%+v",
		time.Since(start), rst, err, req, rsp)
	return nil
}

// GetVirtualBackgroundList 获取虚拟背景(内部服务调用，校验调用方)
func (s *WemeetMeetingTemplateOidbServiceImpl) GetVirtualBackgroundList(ctx context.Context,
	req *pb.GetVirtualBackgroundListReq, rsp *pb.GetVirtualBackgroundListRsp) (err error) {
	metrics.IncrCounter("GetVirtualBackgroundList.total", 1)
	start := time.Now()

	rst, err := handleGetVirtualBackgroundList(ctx, req, rsp)
	if err != nil {
		metrics.IncrCounter("getVirtualBackgroundList.failed", 1)
		rsp.ErrorCode = proto.Int32(rst)
		rsp.ErrorMessage = proto.String(err.Error())
	} else {
		metrics.IncrCounter("getVirtualBackgroundList.success", 1)
		rsp.ErrorCode = proto.Int32(0)
		rsp.ErrorMessage = proto.String("ok")
	}
	log.InfoContextf(ctx, "GetVirtualBackgroundList, cost:%v, rst:%v, err:%+v, req:%+v, rsp:%+v",
		time.Since(start), rst, err, req, rsp)
	return nil
}

// QueryVirtualBackgroundByID 通过ID查询虚拟背景
func (s *WemeetMeetingTemplateOidbServiceImpl) QueryVirtualBackgroundByID(ctx context.Context,
	req *pb.QueryVirtualBackgroundReq, rsp *pb.QueryVirtualBackgroundRsp) error {
	metrics.IncrCounter("QueryVirtualBackgroundByID Total", 1)
	start := time.Now()

	rst, err := handleQueryVirtualBackground(ctx, req, rsp)
	if err != nil {
		metrics.IncrCounter("QueryVirtualBackgroundByID fail", 1)
		rsp.ErrorCode = proto.Int32(rst)
		rsp.ErrorMessage = proto.String(err.Error())
	} else {
		metrics.IncrCounter("QueryVirtualBackgroundByID success", 1)
		rsp.ErrorCode = proto.Int32(0)
		rsp.ErrorMessage = proto.String("ok")
	}
	log.InfoContextf(ctx, "QueryVirtualBackgroundByID, cost:%v, rst:%v, err:%+v, req:%+v, rsp:%+v",
		time.Since(start), rst, err, req, rsp)
	return nil
}

// VirtualBackgroundSafeCallback 虚拟背景审核回调
func (s *WemeetMeetingTemplateOidbServiceImpl) VirtualBackgroundSafeCallback(ctx context.Context,
	req *safePb.GetImageCallbackReq, rsp *safePb.GetImageCallbackRsp) error {
	metrics.IncrCounter("VirtualBackgroundSafeCallback Total", 1)
	start := time.Now()

	rst, err := HandleVirtualBackgroundSafeCallback(ctx, req)
	if err != nil {
		metrics.IncrCounter("VirtualBackgroundSafeCallback fail", 1)
		rsp.ErrorCode = proto.Int32(rst)
		rsp.ErrorMsg = proto.String(err.Error())
	} else {
		metrics.IncrCounter("VirtualBackgroundSafeCallback success", 1)
		rsp.ErrorCode = proto.Int32(0)
		rsp.ErrorMsg = proto.String("ok")
	}
	log.InfoContextf(ctx, "VirtualBackgroundSafeCallback, cost:%v, rst:%v, err:%+v, req:%+v, rsp:%+v",
		time.Since(start), rst, err, req, rsp)
	return nil
}

// handleSetVirtualBackground ..
func handleSetVirtualBackground(ctx context.Context, req *pb.SetVirtualBackgroundReq,
	rsp *pb.SetVirtualBackgroundRsp) (code int32, err error) {

	//校验请求参数
	if req.GetUint64MeetingId() == 0 || req.GetUint32OperateType() == 0 {
		log.InfoContextf(ctx, "SetVirtualBackground Invalid Param, req:%+v", req)
		return util.InvalidParam, errors.New("invalid Param")
	}
	meetInfo, _, _, err := rpc.GetMeetingInfo(ctx, req.GetUint64MeetingId())
	if err != nil {
		metrics.IncrCounter("rpc.GetMeetingInfo fail", 1)
		log.ErrorContextf(ctx, "rpc.GetMeetingInfo fail, meetID:%v, err:%+v", req.GetUint64MeetingId(), err)
		return util.ERRQueryMeetInfo, err
	}
	//判断操作者身份
	if meetInfo.GetStrCreatorAppUid() != req.GetStrOperateAppuid() ||
		meetInfo.GetUint32CreatorSdkappid() != req.GetUint32OperateAppid() {
		hasPower := rpc.CheckMeetingPermission(ctx, req.GetUint32OperateAppid(), req.GetStrOperateAppuid(),
			meetInfo.GetUint32CreatorSdkappid())
		if !hasPower {
			metrics.IncrCounter("SetVirtualBackground Not Permission", 1)
			return int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_NO_PERMISSION),
				errors.New("not permission")
		}
	}

	switch req.GetUint32OperateType() {
	case SetBackgroundOperateTypeAdd: //添加
		return addVirtualBackground(ctx, meetInfo, req, rsp)
	case SetBackgroundOperateTypeUpdate: //更新
		return updateVirtualBackground(ctx, meetInfo, req, rsp)
	case SetBackgroundOperateTypeDelete: //删除
		return deleteVirtualBackground(ctx, req, rsp)
	default:
		return util.InvalidParam, errors.New("OperateType Invalid")
	}

	return 0, nil
}

// handleGetVirtualBackgroundList ..
func handleGetVirtualBackgroundList(ctx context.Context,
	req *pb.GetVirtualBackgroundListReq, rsp *pb.GetVirtualBackgroundListRsp) (int32, error) {
	// 校验参数
	if req.GetUint64MeetingId() == 0 {
		metrics.IncrCounter("GetVirtualBackgroundList invalid param", 1)
		return util.InvalidParam, errors.New("invalid param")
	}
	// 调用方校验
	if req.GetStrAppFrom() != AppFromWemeetLayoutCenter {
		metrics.IncrCounter("GetVirtualBackgroundList Not Permission", 1)
		return int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_NO_PERMISSION),
			errors.New("not permission")
	}

	// 查询背景图列表
	backgroundInfoList, code, err := getBackgroundInfoList(ctx, req.GetUint32CarryCond(), req.GetUint64MeetingId())
	if err != nil {
		return code, err
	}
	rsp.MsgBackgroundInfo = backgroundInfoList
	return 0, nil
}

// handleQueryVirtualBackgroundList ..
func handleQueryVirtualBackgroundList(ctx context.Context, req *pb.QueryVirtualBackgroundListReq,
	rsp *pb.QueryVirtualBackgroundListRsp) (int32, error) {

	//校验参数
	if req.GetUint64MeetingId() == 0 {
		metrics.IncrCounter("QueryVirtualBackgroundList invalid param", 1)
		return util.InvalidParam, errors.New("invalid param")
	}
	meetID := req.GetUint64MeetingId()
	//操作者权限校验
	meetInfo, _, _, err := rpc.GetMeetingInfo(ctx, meetID)
	if err != nil {
		metrics.IncrCounter("rpc.GetMeetingInfo fail", 1)
		log.ErrorContextf(ctx, "rpc.GetMeetingInfo fail, meetID:%v, err:%+v", meetID, err)
		return util.ERRQueryMeetInfo, err
	}
	//判断操作者身份
	if meetInfo.GetStrCreatorAppUid() != req.GetStrOperateAppuid() ||
		meetInfo.GetUint32CreatorSdkappid() != req.GetUint32OperateAppid() {
		hasPower := rpc.CheckMeetingPermission(ctx, req.GetUint32OperateAppid(), req.GetStrOperateAppuid(),
			meetInfo.GetUint32CreatorSdkappid())
		if !hasPower {
			metrics.IncrCounter("QueryVirtualBackgroundList Not Permission", 1)
			return int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_NO_PERMISSION),
				errors.New("not permission")
		}
	}
	// 查询背景图列表
	backgroundInfoList, code, err := getBackgroundInfoList(ctx, req.GetUint32CarryCond(), meetID)
	if err != nil {
		return code, err
	}
	rsp.MsgBackgroundInfo = backgroundInfoList
	return 0, nil
}

// getBackgroundInfoList 查询背景图列表
func getBackgroundInfoList(ctx context.Context, carryCond uint32, meetingID uint64) ([]*pb.BackgroundInfo,
	int32, error) {
	backgroundInfoList := make([]*pb.BackgroundInfo, 0)
	background := bg.NewBackground()

	//查询默认背景图
	if carryCond == 0 || carryCond == 2 {
		defBgIds, err := background.GetDefaultBackgroundSortSet(ctx)
		if err != nil {
			metrics.IncrCounter("GetDefaultBackgroundSortSet fail", 1)
			log.ErrorContextf(ctx, "GetDefaultBackgroundSortSet fail, meetingID:%v, err:%+v", meetingID, err)
			return backgroundInfoList, util.ErrQueryCacheData, errors.New("get defBackgroundIDs fail")
		}
		defBackgroundInfoList, err := background.GetBackground(ctx, defBgIds)
		if err != nil {
			metrics.IncrCounter("GetDefBackground fail", 1)
			log.ErrorContextf(ctx, "GetDefBackground fail, meetingID:%v, bgIDs:%+v", meetingID, defBgIds)
			return backgroundInfoList, util.ErrQueryCacheData, errors.New("get defBackgroundInfo fail")
		}
		backgroundInfoList = append(backgroundInfoList, defBackgroundInfoList...)
	}
	//查询会议背景图信息
	{
		meetBgIDs, err := background.GetMeetBackgroundSortSet(ctx, meetingID)
		if err != nil {
			metrics.IncrCounter("GetMeetBackgroundSortSet fail", 1)
			log.ErrorContextf(ctx, "GetMeetBackgroundSortSet fail, meetingID:%v, err:%+v", meetingID, err)
			return backgroundInfoList, util.ErrQueryCacheData, errors.New("get meetBackgroundIDs fail")
		}
		meetBackgroundInfoList, err := background.GetBackground(ctx, meetBgIDs)
		if err != nil {
			metrics.IncrCounter("GetBackground fail", 1)
			log.ErrorContextf(ctx, "GetBackground fail, meetingID:%v, bgIDs:%+v", meetingID, meetBgIDs)
			return backgroundInfoList, util.ErrQueryCacheData, errors.New("get meetBackgroundInfo fail")
		}
		//当图片状态一直是待审核时(时间间隔七彩石配置)，将图片状态改成审核超时
		for _, data := range meetBackgroundInfoList {
			if ((data.GetUint64IncrTime() + uint64(config_rainbow.GetBackgroundConfAuditTime())) <
				uint64(time.Now().Unix())) && (data.GetUint32PicStatus() == bg.PicStatusAuditing) {
				data.Uint32PicStatus = proto.Uint32(bg.PicStatusAuditTimeOut)
			}
		}
		backgroundInfoList = append(backgroundInfoList, meetBackgroundInfoList...)
	}
	return backgroundInfoList, 0, nil
}

func handleQueryVirtualBackground(ctx context.Context, req *pb.QueryVirtualBackgroundReq,
	rsp *pb.QueryVirtualBackgroundRsp) (int32, error) {

	//参数校验
	if len(req.GetInt64BackgroundId()) == 0 {
		log.ErrorContextf(ctx, "QueryVirtualBackground Id empty")
		return util.InvalidParam, errors.New("background id empty")
	}

	background := bg.NewBackground()

	backgroundInfoList, err := background.GetBackground(ctx, req.GetInt64BackgroundId())
	if err != nil {
		metrics.IncrCounter("background.GetBackground fail", 1)
		log.ErrorContextf(ctx, "background.GetBackground fail, ids:%+v, err:%+v",
			req.GetInt64BackgroundId(), err)
		return util.ErrQueryCacheData, err
	}
	//当图片状态一直是待审核时(时间间隔七彩石配置)，将图片状态改成审核超时
	for _, data := range backgroundInfoList {
		if ((data.GetUint64IncrTime() + uint64(config_rainbow.GetBackgroundConfAuditTime())) <
			uint64(time.Now().Unix())) && (data.GetUint32PicStatus() == bg.PicStatusAuditing) {
			data.Uint32PicStatus = proto.Uint32(bg.PicStatusAuditTimeOut)
		}
	}
	rsp.MsgBackgroundInfo = backgroundInfoList
	return 0, nil
}

// HandleVirtualBackgroundSafeCallback 虚拟背景图片审核回调
func HandleVirtualBackgroundSafeCallback(ctx context.Context, req *safePb.GetImageCallbackReq) (int32, error) {

	//校验参数
	if req.ResultCode == nil && req.GetSer().TraceId == nil {
		metrics.IncrCounter("VirtualBackgroundSafeCallback invalid param", 1)
		log.ErrorContextf(ctx, "VirtualBackgroundSafeCallback invalid param, req:%+v", req)
		return util.InvalidParam, errors.New("invalid param")
	}
	backgroundID, err := strconv.ParseInt(req.GetSer().GetTraceId(), 10, 64)
	if err != nil {
		metrics.IncrCounter("backgroundID invalid", 1)
		log.ErrorContextf(ctx, "VirtualBackgroundSafeCallback backgroundID invalid, req:%+v", req)
		return util.InvalidParam, errors.New("invalid param")
	}
	meetID := req.GetMeet().GetMeetingId()
	ids := []int64{backgroundID}
	background := bg.NewBackground()
	if meetID == 0 { //会议ID为0，为企业
		log.InfoContextf(ctx, "update default background info, id:%v", backgroundID)
	}
	backgroundList, err := background.GetBackground(ctx, ids)
	if err != nil {
		metrics.IncrCounter("background.GetBackground fail", 1)
		log.ErrorContextf(ctx, "background.GetBackground fail, meetID:%v, ids:%+v, err:%+v",
			meetID, ids, err)
		return util.ErrQueryCacheData, err
	}
	for _, val := range backgroundList {
		if req.GetResultCode() == 0 {
			val.Uint32PicStatus = proto.Uint32(bg.PicStatusNormal)
		} else {
			val.Uint32PicStatus = proto.Uint32(bg.PicStatusAuditFail)
		}
	}
	err = background.SetBackgroundList(ctx, backgroundList)
	if err != nil {
		metrics.IncrCounter("background.SetBackgroundList fail", 1)
		log.ErrorContextf(ctx, "background.SetBackgroundList fail, meetID:%v, background:%+v, err:%+v",
			meetID, backgroundList, err)
		return util.ErrSetBackground, err
	}

	return 0, nil
}

// addVirtualBackground ..
func addVirtualBackground(ctx context.Context, meetInfo *cachePb.MeetingInfo, req *pb.SetVirtualBackgroundReq,
	rsp *pb.SetVirtualBackgroundRsp) (int32, error) {

	//参数校验
	if req.GetMsgBackgroundInfo() == nil {
		metrics.IncrCounter("addVirtualBackground invalid param", 1)
		return util.InvalidParam, errors.New("invalid param")
	}

	background := bg.NewBackground()

	//图片添加超上限
	backgroundCnt, _ := background.GetMeetBackgroundSize(ctx, req.GetUint64MeetingId())
	if backgroundCnt > config_rainbow.GetBackgroundConfSize() {
		metrics.IncrCounter("GetMeetBackgroundSize Limit", 1)
		log.InfoContextf(ctx, "GetMeetBackgroundSize Limit, Cnt:%v", backgroundCnt)
		return util.ErrBackgroundLimit, errors.New("num limit")
	}

	backgroundList := make([]*pb.BackgroundInfo, 0, len(req.GetMsgBackgroundInfo()))
	backgroundIDList := make([]int64, 0, len(req.GetMsgBackgroundInfo()))
	for _, val := range req.GetMsgBackgroundInfo() {
		backgroundInfo := val
		//对背景图生成唯一索引
		backgroundID := background.GetBackgroundIndex(ctx)

		backgroundInfo.Int64BackgroundId = proto.Int64(backgroundID)
		backgroundInfo.Uint32PicStatus = proto.Uint32(bg.PicStatusAuditing)
		backgroundInfo.Uint64IncrTime = proto.Uint64(uint64(time.Now().Unix()))
		backgroundList = append(backgroundList, backgroundInfo)
		backgroundIDList = append(backgroundIDList, backgroundID)
		//保存图片
		log.InfoContextf(ctx, "SetBackgroundList, BackgroundInfo:%+v", backgroundInfo)
	}

	err := background.SetBackgroundList(ctx, backgroundList)
	if err != nil {
		metrics.IncrCounter("background.SetBackgroundList fail", 1)
		log.ErrorContextf(ctx, "background.SetBackgroundList fail, meetID:%v, backgroundList:%+v",
			meetInfo.GetUint64MeetingId(), backgroundList)
		return util.ErrSetBackground, err
	}

	//保存到Sort Set
	err = background.SetMeetBackgroundSortSet(ctx, meetInfo.GetUint64MeetingId(), backgroundIDList)
	if err != nil {
		metrics.IncrCounter("background.SetMeetBackgroundSortSet fail", 1)
		log.ErrorContextf(ctx, "SetMeetBackgroundSortSet fail, meetID:%v, err:%+v",
			meetInfo.GetUint64MeetingId(), err)
		return util.ErrSetBackground, err
	}

	//虚拟背景图片上传，cos信息通知layout
	backgroundCosNotify(ctx, meetInfo.GetUint64MeetingId(), backgroundList)

	//图片审核放在后面，因为vip客户会不过审核立马回调
	for _, info := range backgroundList {
		if info.GetStrPicUrl() == "" {
			valList := []*pb.BackgroundInfo{info}
			background.GetTempImageUrlByID(ctx, valList)
		}

		callBackSwitch := config_rainbow.GetCallBackConf().GateWayImgBackgroundActon
		log.InfoContextf(ctx, "ctlrpcImgSafetyAudit GateWayImgBackgroundActon: %v", callBackSwitch)
		if callBackSwitch {
			imgReq := &ctlrpc.ImgSafetyAuditReq{
				Uuid:              strconv.FormatInt(info.GetInt64BackgroundId(), 10),
				AppId:             meetInfo.GetUint32CreatorSdkappid(),
				AppUid:            meetInfo.GetStrCreatorAppUid(),
				MeetingId:         req.GetUint64MeetingId(),
				Url:               info.GetStrPicUrl(),
				Scenes:            GateWayImgBackgroundScenes,
				Action:            GateWayImgBackgroundActon,
				CallbackTarget:    CallbackTarget,
				CallbackUri:       ImgBackgroundCallBackUri,
				CallbackEnvName:   util.GetEnvName(),
				CallbackNameSpace: util.GetNameSpace(),
				CallbackProtoType: CallbackProtoType,
			}
			ctlrpc.ImgSafetyAuditV2(ctx, imgReq)
		} else {
			ctlrpc.ImgSafetyAudit(ctx, strconv.FormatInt(info.GetInt64BackgroundId(), 10),
				meetInfo.GetUint32CreatorSdkappid(), meetInfo.GetStrCreatorAppUid(),
				req.GetUint64MeetingId(), info.GetStrPicUrl(), GateWayImgBackgroundScenes)
		}
	}

	return 0, nil
}

// backgroundCosNotify 虚拟背景图片上传，cos信息通知layout
func backgroundCosNotify(ctx context.Context, meetingId uint64, backgrounds []*pb.BackgroundInfo) error {
	log.DebugContextf(ctx, "backgroundCosNotify, meetID:%+v, backgroundCosNotifySwitch:%+v, backgrounds:%+v",
		meetingId, config_rainbow.GetBackgroundCosNotifySwitch(), backgrounds)
	var err error
	if config_rainbow.GetBackgroundCosNotifySwitch() {
		// 异步通知，用于虚拟背景图片cos安全校验
		// 注意：虚拟背景图片cos资源清理不依赖于该通知调用成功率@thunderge
		newCtx := trpc.CloneContext(ctx)
		go func() {
			meet_util.DefPanicFun()
			err = rpc.GuestVbCosListNotify(newCtx, meetingId, backgrounds)
			// 失败重试
			if err != nil {
				err = rpc.GuestVbCosListNotify(newCtx, meetingId, backgrounds)
				if err != nil {
					log.ErrorContextf(ctx, "GuestVbCosListNotify fail, meetID:%+v, backgroundList: %+v, err:%+v",
						meetingId, backgrounds, err)
				}
			}
		}()
	}
	return err
}

func updateVirtualBackground(ctx context.Context, meetingInfo *cachePb.MeetingInfo, req *pb.SetVirtualBackgroundReq,
	rsp *pb.SetVirtualBackgroundRsp) (code int32, err error) {

	if req.GetMsgBackgroundInfo() == nil {
		metrics.IncrCounter("updateVirtualBackground invalid parameter", 1)
		log.ErrorContextf(ctx, "updateVirtualBackground invalid parameter, meetID:%v", req.GetUint64MeetingId())
		return util.InvalidParam, errors.New("invalid parameter")
	}

	background := bg.NewBackground()

	ids := make([]int64, 0, len(req.GetMsgBackgroundInfo()))
	for _, val := range req.GetMsgBackgroundInfo() {
		ids = append(ids, val.GetInt64BackgroundId())
	}
	backgroundList, err := background.GetBackground(ctx, ids)
	if err != nil {
		metrics.IncrCounter("background.GetBackground fail", 1)
		log.ErrorContextf(ctx, "background.GetBackground fail, meetID:%v, err:%+v",
			req.GetUint64MeetingId(), err)
		return util.ErrUpdateBackground, err
	}
	for _, val := range backgroundList {
		val.Uint32PicStatus = proto.Uint32(bg.PicStatusAuditing) //更新接口现只支持审核重试
	}
	err = background.SetBackgroundList(ctx, backgroundList)
	if err != nil {
		metrics.IncrCounter("background.SetBackgroundList fail", 1)
		log.ErrorContextf(ctx, "background.SetBackgroundList fail, meetID:%v, err:%+v",
			req.GetUint64MeetingId(), err)
		return util.ErrUpdateBackground, err
	}

	//虚拟背景图片上传，cos信息通知layout
	backgroundCosNotify(ctx, meetingInfo.GetUint64MeetingId(), backgroundList)

	//选择开关
	callBackSwitch := config_rainbow.GetCallBackConf().GateWayImgBackgroundActon
	log.InfoContextf(ctx, "ctlrpcImgSafetyAudit GateWayImgBackgroundActon: %v", callBackSwitch)
	for _, val := range backgroundList {
		if callBackSwitch {
			imgReq := &ctlrpc.ImgSafetyAuditReq{
				Uuid:              strconv.FormatInt(val.GetInt64BackgroundId(), 10),
				AppId:             meetingInfo.GetUint32CreatorSdkappid(),
				AppUid:            meetingInfo.GetStrCreatorAppUid(),
				MeetingId:         req.GetUint64MeetingId(),
				Url:               val.GetStrPicUrl(),
				Scenes:            GateWayImgBackgroundScenes,
				Action:            GateWayImgBackgroundActon,
				CallbackTarget:    CallbackTarget,
				CallbackUri:       ImgBackgroundCallBackUri,
				CallbackEnvName:   util.GetEnvName(),
				CallbackNameSpace: util.GetNameSpace(),
				CallbackProtoType: CallbackProtoType,
			}
			ctlrpc.ImgSafetyAuditV2(ctx, imgReq)
		} else {
			ctlrpc.ImgSafetyAudit(ctx, strconv.FormatInt(val.GetInt64BackgroundId(), 10),
				meetingInfo.GetUint32CreatorSdkappid(), meetingInfo.GetStrCreatorAppUid(),
				req.GetUint64MeetingId(), val.GetStrPicUrl(), GateWayImgBackgroundScenes)
		}
	}

	rsp.MsgBackgroundInfo = backgroundList
	return 0, nil
}

func deleteVirtualBackground(ctx context.Context, req *pb.SetVirtualBackgroundReq,
	rsp *pb.SetVirtualBackgroundRsp) (code int32, err error) {

	//校验参数
	if req.GetMsgBackgroundInfo() == nil {
		metrics.IncrCounter("deleteVirtualBackground invalid parameter", 1)
		log.ErrorContextf(ctx, "deleteVirtualBackground invalid parameter, meetID:%v", req.GetUint64MeetingId())
		return util.InvalidParam, errors.New("invalid parameter")
	}

	background := bg.NewBackground()

	IDs := make([]int64, 0, len(req.GetMsgBackgroundInfo()))
	for _, val := range req.GetMsgBackgroundInfo() {
		if val.GetInt64BackgroundId() == 0 {
			metrics.IncrCounter("BackgroundId invalid", 1)
			log.ErrorContextf(ctx, "BackgroundId invalid, meetID:%v", req.GetUint64MeetingId())
			continue
		}
		IDs = append(IDs, val.GetInt64BackgroundId())
	}

	err = background.DeleteMeetBackground(ctx, req.GetUint64MeetingId(), IDs)
	if err != nil {
		metrics.IncrCounter("background.DeleteMeetBackground fail", 1)
		log.ErrorContextf(ctx, "background.DeleteMeetBackground fail, meetID:%v, err:%+v",
			req.GetUint64MeetingId(), err)
		return util.ErrDeleteBackground, err
	}
	return 0, nil
}

func getUserInfo(ctx context.Context) (appId uint32, appUid string, err error) {
	header := http.Head(ctx)
	cookies := header.Request.Cookies()
	for _, val := range cookies {
		if val.Name == "app_uid" {
			appUid = val.Value
		} else if val.Name == "corp_id" {
			data, err := strconv.ParseUint(val.Value, 10, 64)
			if err != nil {
				return appId, appUid, err
			}
			appId = uint32(data)
		}
	}
	return appId, appUid, nil
}

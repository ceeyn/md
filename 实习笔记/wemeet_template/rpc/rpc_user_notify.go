package rpc

import (
	"context"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"google.golang.org/protobuf/proto"
	"meeting_template/util"
	"strconv"
	"time"

	xcastPb "git.code.oa.com/trpcprotocol/wemeet/common_xcast_im_conversion_logic"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	notifyPb "git.code.oa.com/trpcprotocol/wemeet/wemeet_user_notify"
)

const (
	NO_PLAY_MUSIC = 0
	PLAY_MUSIC = 1
)

// 通知用户暖场物料更新
func NotifyWarmUpDataUpdate(ctx context.Context, meetingId uint64, appId uint32, appUid string,
	warmUpData *pb.WarmUpData) {
	util.ReportOne(util.UpdateWarmUpUserNotify) //[UpdateWarmUpUserNotify]请求
	startTime := time.Now()
	req := &notifyPb.SendMsgToUserNotifyReq{}
	FillNotifyReq(ctx, meetingId, appId, appUid, warmUpData, req)

	//初始化请求代理
	proxy := notifyPb.NewWemeetUserNotifyOidbClientProxy()
	//发送请求
	rsp, err := proxy.SendMsgToUserNotify(ctx, req)
	cost := time.Since(startTime).Milliseconds()
	log.DebugContextf(ctx, "NotifyWarmUpDataUpdate,meetingId:%v,warmUpData:%v,req:%v, time cost:%v, rsp:%v",
		meetingId, warmUpData, req, cost, rsp)
	if err != nil || rsp == nil || rsp.GetErrorCode() != 0 {
		util.ReportOne(util.UpdateWarmUpUserNotifyFail) //[UpdateWarmUpUserNotify]请求失败
		log.ErrorContextf(ctx, "UpdateWarmUpUserNotify failed, err is:%v, rsp:%v", err, rsp)
	}
}

func FillNotifyReq(ctx context.Context, meetingId uint64, appId uint32, appUid string,
	warmUpData *pb.WarmUpData, req *notifyPb.SendMsgToUserNotifyReq) {
	req.PushHead = &notifyPb.TPushHeadInfo{}
	req.PushHead.MeetingId = proto.Uint64(meetingId)
	// 构造消息唯一ID
	uuid := (strconv.FormatUint(meetingId, 10) + "_" + strconv.FormatInt(util.NowMs(), 10))
	req.PushHead.IdUuid = proto.String(uuid)
	// 暖场模式内容push通知信令为2874
	req.MsgId = proto.Int32(util.WarmUpNotifyMsgId)
	sendObj := &notifyPb.TMsgSendObject{
		EnumSendType: notifyPb.USER_NOTIFY_SEND_TYPE_USER_NOTIFY_SEND_TYPE_IMMEDIATELY.Enum(),
	}
	sendObj.TargetUser = &notifyPb.TargetUser{}
	// SEND_OBJECT_WARMUP_MEMBERS = 1 << 23 //所有的暖场视频的成员
	sendObj.TargetUser.BitTarget = proto.Uint64(1 << 23)
	req.MsgSendObj = append(req.MsgSendObj, sendObj)
	req.SendDescription = proto.String(util.WarmUpNotifySendDescription)
	req.FmUser = &notifyPb.TUserInfo{}
	req.FmUser.AppId = proto.Uint32(appId)
	req.FmUser.AppUid = proto.String(appUid)
	xcastReq := &xcastPb.MeetingLogicConversionReqBody{}
	xcastReq.MsgS2CWarmUpMediaNotifyReqBody = &xcastPb.MeetingLogicS2CWarmupMediaNotifyReqBody{}
	for _, val := range warmUpData.WarmupImgList {
		xcastReq.MsgS2CWarmUpMediaNotifyReqBody.StrWarmUpPictureList = append(
			xcastReq.MsgS2CWarmUpMediaNotifyReqBody.StrWarmUpPictureList, *val.Url)
	}
	canShowVideo := false
	for _, val := range warmUpData.WarmupVideoList {
		if val.GetUrl() != "" && val.GetStreamUrl() != "" {
			canShowVideo = true
		}
		item := &xcastPb.WarmupMediaItem{}
		item.StrWebUrl = val.Url
		item.StrStreamUrl = val.StreamUrl
		xcastReq.MsgS2CWarmUpMediaNotifyReqBody.WarmUpMediaList = append(
			xcastReq.MsgS2CWarmUpMediaNotifyReqBody.WarmUpMediaList, item)
	}
	// 推送的物料类型
	xcastReq.MsgS2CWarmUpMediaNotifyReqBody.Uint32WarmUpType = proto.Uint32(warmUpData.GetUint32WarmupUseType())
	if warmUpData.GetUint32WarmupUseType() == util.WarmUpShowVideo && !canShowVideo {
		xcastReq.MsgS2CWarmUpMediaNotifyReqBody.Uint32WarmUpType = proto.Uint32(util.WarmUpShowImg)
	}
	// 音乐的push相关
	xcastReq.MsgS2CWarmUpMediaNotifyReqBody.Uint32PlayMusic = proto.Uint32(NO_PLAY_MUSIC)
	if warmUpData.GetUint32WarmupUseType() == 0 && warmUpData.GetUint32PlayMusicType() != 0 {
		xcastReq.MsgS2CWarmUpMediaNotifyReqBody.Uint32PlayMusic = proto.Uint32(PLAY_MUSIC)   // 图片音乐的物料类型，勾选了音乐，则播放音乐
		for _, item := range warmUpData.GetWarmupMusicList() {
			if item.GetMusicType() == warmUpData.GetUint32PlayMusicType() {
				xcastReq.MsgS2CWarmUpMediaNotifyReqBody.StrWarmUpMusicList = append(
					xcastReq.MsgS2CWarmUpMediaNotifyReqBody.StrWarmUpMusicList, item.GetUrl())
			}
		}
	}
	//序列化数据
	data, err := proto.Marshal(xcastReq)
	if err != nil {
		util.ReportOne(util.JsonMarshalFail)
		log.Errorf("MeetingLogicS2CWarmupMediaNotifyReqBody Marshal error[%v]", err.Error())
	} else {
		req.MsgBody = append(req.MsgBody, data...)
	}
}

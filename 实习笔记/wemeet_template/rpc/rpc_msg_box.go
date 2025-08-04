package rpc

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/http"
	"git.code.oa.com/trpc-go/trpc-go/log"
	msgBox "git.code.oa.com/trpcprotocol/wemeet/common_msgbox"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"meeting_template/config"
	"meeting_template/util"
	"strings"
	"time"
)

const (
	ContentTypeText   = 3    // 内容样式类型 (1=文本带参数, 2=图文, 3=纯文本)
	MessageTypeSystem = "11" // 系统通知类型为11
	TASKID            = "30" // taskID为30
)

// MessageV1 老版本消息定义
type MessageV1 struct {
	Type       interface{}    `json:"type,omitempty"`       // 内容样式类型 (1=文本带参数, 2=图文, 3=纯文本)
	Time       interface{}    `json:"time,omitempty"`       // 时间
	TaskId     interface{}    `json:"taskId,omitempty"`     // 任务id
	ExpId      string         `json:"expId,omitempty"`      // 灯塔实验id
	ExpVersion string         `json:"expVersion,omitempty"` // 灯塔实验版本
	Width      string         `json:"width,omitempty"`      // 宽
	Height     string         `json:"height,omitempty"`     // 高
	IconType   interface{}    `json:"iconType,omitempty"`   // 图标类型
	IconName   string         `json:"iconName,omitempty"`   // 图标名称
	IconUrl    string         `json:"iconUrl,omitempty"`    // 图标链接
	Msg        []*MessageData `json:"msg,omitempty"`        // 消息列表（支持多图文）
}

// MessageData 消息数据
type MessageData struct {
	Type    interface{} `json:"type,omitempty"`    // 类型，取值为image时，fileUrl表示图片链接
	FileUrl string      `json:"fileUrl,omitempty"` // 图片链接
	Name    string      `json:"name,omitempty"`    // 图片名
	Title   string      `json:"title,omitempty"`   // 标题
	Summary string      `json:"summary,omitempty"` // 概览
	Url     string      `json:"url,omitempty"`     // 跳转链接
}

// getSignData 生成签名
func getSignData() (string, string) {
	secret := config.Conf.MsgBoxConfig.Secret
	log.Infof("getSignData getsecret:%+v", secret)
	id := strings.Replace(uuid.New().String(), "-", "", -1)
	hash := md5.Sum([]byte(secret + ":" + id))
	return id, hex.EncodeToString(hash[:])
}

// MsgBoxClientProxy 消息中心client
var MsgBoxClientProxy = http.NewClientProxy("trpc.wemeet.common.message")

// SendMessageToMsgBox.. 发送消息到消息中心 和系统通知
func SendMessageToMsgBox(ctx context.Context, uid string, sdkAppid uint32, meetingId uint64) (uint32, error) {
	uuidStr, sign := getSignData()
	ts := time.Now().Format("2006-01-02 15:04")
	summary := BuildSummaryContent(ctx, meetingId)
	// 具体的消息内容
	msgData := &MessageData{
		Title:   "暖场视频提交处理失败", // 标题
		Summary: summary,      // 内容
	}
	// 消息结构体
	msgBoxInfo := MessageV1{
		Type:   proto.Int32(ContentTypeText), // 这里填的3 表示纯文本
		Time:   proto.String(ts),             // 当前时间
		TaskId: proto.String(TASKID),         // taskID: 30
		Msg:    []*MessageData{msgData},
	}
	// 打印一下
	log.InfoContextf(ctx, "SendMessageToMsgBox msgBoxInfo:%+v", msgBoxInfo)

	msgBoxData, err := json.Marshal(msgBoxInfo)
	if err != nil {
		attr.AttrAPI(35736054, 1)
		log.ErrorContextf(ctx, "json.(Marshal) msgBoxInfo failed! uid:%+v", uid)
	}

	// 信鸽弹窗
	targetUrl := "https://meeting.tencent.com/user-center/webinar/detail/" + fmt.Sprint(meetingId)
	tpnsParams := &msgBox.TpnsParams{
		Title:     proto.String("暖场视频提交处理失败"),
		Content:   proto.String("您上传的会议暖场视频处理失败，请前往腾讯会议官网用户中心修改视频后重试。"),
		TargetUrl: proto.String(targetUrl), //目标链接 会议详情
		AppOs:     proto.Int32(99),         // ios和安卓终端
	}

	req := &msgBox.PushMsgTaskReq{
		RedDot:       proto.Int32(2048), // 红点
		AppOs:        proto.Int32(255),  // 这个填什么 全平台：255
		CorpId:       proto.Uint32(sdkAppid),
		Uid:          proto.String(uid),                // 目标用户uid
		Type:         proto.String(MessageTypeSystem),  // type 填什么 11 系统消息
		Uuid:         proto.String(uuidStr),            // uuid
		Sign:         proto.String(sign),               // 签名
		Data:         proto.String(string(msgBoxData)), // 具体的消息内容
		Caller:       proto.Int32(3),                   // 这个填什么 3
		NeedTpnsPush: proto.Bool(true),                 // 需要信鸽通知
		TpnsParams:   tpnsParams,                       // 信鸽通知
	}

	rsp := &msgBox.PushMsgTaskRsp{}
	err = MsgBoxClientProxy.Post(ctx, "/msgbox.message/PushMsgTask", req, rsp)
	if err != nil {
		attr.AttrAPI(35736055, 1)
		log.ErrorContextf(ctx, "MsgBoxClientProxy post failed. req:%+v, rsp:%+v, err:%+v", req, rsp, err)
		return 310, err
	}
	if rsp.GetCode() != "0" {
		attr.AttrAPI(35736056, 1)
		log.ErrorContextf(ctx, "MsgBoxClientProxy rsp code not 0. req:%+v, rsp:%+v", req, rsp)
		return 311, errors.New("rsp GetCode not 0")
	}
	log.InfoContextf(ctx, "MsgBoxClientProxy post succ. req:%+v, rsp:%+v", req, rsp)
	return 0, nil
}

// 获取通知的消息
func BuildSummaryContent(ctx context.Context, meetingId uint64) string {
	var summaryTpl = "会议号：%+v\n会议主题：%+v\n会议时间：%+v\n您可以前往腾讯会议官网用户中心修改视频后重新提交。"
	meetingInfo, _, _, err := GetMeetingInfo(ctx, meetingId)
	if err != nil {
		attr.AttrAPI(35736057, 1) //[getMeetingInfo]获取会议信息失败
		log.ErrorContextf(ctx, "get meeting info from cache fail, err: %v", err)
		return ""
	}
	meetingCode := meetingInfo.GetStrMeetingCode()
	subject, _ := base64.StdEncoding.DecodeString(string(meetingInfo.GetBytesMeetingSubject()))
	startTs := util.FormatUintTime(meetingInfo.GetUint32OrderStartTime())

	summaryContent := fmt.Sprintf(summaryTpl, meetingCode, string(subject), startTs)
	return summaryContent
}

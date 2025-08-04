package util

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
)

const (
	//WebinarMeetingType  会议类型
	WebinarMeetingType          = 1
	WarmUpShowImg               = 0
	WarmUpShowVideo             = 1
	NotPicImageType             = "0"
	SmallImageType              = "1"
	MiddleImageType             = "2"
	RawImageType                = "3"
	DownloadNotUseCdn           = "0"
	DownloadUseCdn              = "1"
	VideoNotAudit               = 0
	VideoPassAudit              = 1
	VideoAuditFail              = 2
	VideoAuditUploadAction      = "1"
	VideoAuditResultPass        = 0
	VideoAuditResultFail        = 1
	WarmUpNotifyMsgId           = 2874
	WarmUpNotifySendDescription = "WarmUpS2CNotify"
	ErrSwitchStateIllegal       = 1400
	ErrUpdateInviteSwitch       = 1401
	InviteStateClose            = 0
	InviteStateOpen             = 1
	NotAudit                    = 0
	PassAudit                   = 1
	AuditFail                   = 2
	ImgAuditResultPass          = 0
	ImgAuditResultFail          = 1
	ERRLenIllegal               = 3000
	ERRMeetingIdIllegal         = 3001
	ERRCount                    = 3006
	ERRCosIdIllegal             = 3002
	ERRLength                   = 3003
	ERRNameSensitive            = 3004
	ERRDescSensitive            = 3005
	ERRUnmarshal                = 3007

	InvalidParam        = 4001
	ERRQueryMeetInfo    = 4002
	ErrDeleteBackground = 4003
	ErrQueryCacheData   = 4004
	ErrSetBackground    = 4005
	ErrSendHttp         = 4006
	ErrBackgroundLimit  = 4007
	ErrUpdateBackground = 4008
	//=========搜索的错误码==========
	ErrLenMeetingId             = 40000
	ErrSearchIntroduction       = 40001
	ErrSearchItinerary          = 40002
	ErrPageParam                = 19999
	ErrSearchKey                = 20000
	ErrFuzzSearch               = 20001
	RoleAuthMeetingModification = "meeting-modification"
	M3U8TransNoReady            = 1
	M3U8TransDone               = 2
)

const (
	//MeetingExpireTime 会议过期时间 2764800 30+2天
	MeetingExpireTime = 2764800
	//MeetingHalfAYearExpireTime 过期时间 180天
	MeetingHalfAYearExpireTime = 15552000
)

// Now 获取当前时间，单位毫秒
func Now() int64 {
	timestamp := time.Now().UnixNano() / 1e3
	return timestamp
}

// NowS 获取当前时间，单位秒
func NowS() uint32 {
	return uint32(time.Now().UnixNano() / 1e9)
}

// NowMs 获取当前时间，单位毫秒
func NowMs() int64 {
	now := time.Now().UnixNano() / 1e6
	return now
}

// GetBase64Decoded 获取base64解码后的详情
func GetBase64Decoded(ctx context.Context, raw string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		ReportOne(Base64DecodeDescriptionFail)
		log.ErrorContextf(ctx, "base64 decode  failed, raw string: %+v,  error: %+v",
			raw, err)
		decoded = []byte("")
	}
	return decoded
}

// IsCoverListCosIdValid 检测封面图片列表cos_id格式是否合法
func IsCoverListCosIdValid(ctx context.Context, coverList []pb.CoverItem) (bool, string) {
	for _, val := range coverList {
		if IsValidCosId(ctx, *val.RawUrl) && IsValidCosId(ctx, *val.CuttedUrl) {
			continue
		}
		return false, ("cos_id invalid, raw:" + *val.RawUrl + ",cutted:" + *val.CuttedUrl)
	}
	return true, ""
}

// IsValidCosId 检测图片cos_id格式是否合法
func IsValidCosId(ctx context.Context, cosId string) bool {
	// 图片字段为空认为是合法的
	if cosId == "" {
		return true
	}
	result, err := regexp.Match("^(wemeet_webinar|webinar_canvas)/.+", []byte(cosId))
	if err == nil {
		return result
	}
	log.ErrorContextf(ctx, "regexp match error: %+v, cosId:%+v, regard as valid", err, cosId)
	return true
}

// GetCompressed Compressed
func GetCompressed(ctx context.Context, raw string) (string, error) {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)
	writer.Write([]byte(raw))
	if err := writer.Close(); err != nil {
		log.ErrorContextf(ctx, "zlib compress error:%+v, raw string:%+v", err, raw)
		return "", err
	}
	return buf.String(), nil
}

// GetDeCompressed Compressed
func GetDeCompressed(ctx context.Context, compressed string) (string, error) {
	var in bytes.Buffer
	_, err := in.WriteString(compressed)
	if err != nil {
		log.ErrorContextf(ctx, "zlib decompress error:%+v, compressed string:%+v", err, compressed)
		return "", err
	}
	var out bytes.Buffer
	r, err := zlib.NewReader(&in)
	if err != nil {
		log.InfoContextf(ctx, "zlib decompress error:%+v, not compressed data, use raw", err)
		return compressed, nil
	}
	_, err = io.Copy(&out, r)
	if err != nil {
		log.ErrorContextf(ctx, "zlib decompress io copy error:%+v, compressed string:%+v", err, compressed)
		return "", err
	}
	return out.String(), nil
}

// GetSerializedJsonStr 获取序列化后的json字符串
func GetSerializedJsonStr(ctx context.Context, val interface{}) (string, error) {
	buf, err := json.Marshal(val)
	if err != nil {
		log.InfoContextf(ctx, "GetSerializedJsonStr, return:")
		return "", err
	}
	retString := string(buf)
	if "null" == retString {
		retString = ""
	}
	log.InfoContextf(ctx, "GetSerializedJsonStr, str buf:%v, retString:%v",
		string(buf), retString)
	return retString, nil
}

// CheckCoverListFormat 检查图片列表格式是否正常，异常返回err,正常返回nil
func CheckCoverListFormat(ctx context.Context, coverItems []*pb.CoverItem) error {
	for _, val := range coverItems {
		decodedRawUrl := string(GetBase64Decoded(ctx, val.GetRawUrl()))
		decodedCuttedUrl := string(GetBase64Decoded(ctx, val.GetCuttedUrl()))
		if !IsValidCosId(ctx, decodedRawUrl) || !IsValidCosId(ctx, decodedCuttedUrl) {
			err := errors.New("cover url cos id invalid, raw cosId:" + decodedRawUrl +
				", cutted cosId:" + decodedCuttedUrl)
			log.ErrorContextf(ctx, "hit cos id format filter, raw cosId:%v, cutted cosId:%v",
				decodedRawUrl, decodedCuttedUrl)
			return err
		}
	}
	return nil
}

// CheckWarmUpDataFormat 检查暖场物料格式是否正常，异常返回err,正常返回nil
func CheckWarmUpDataFormat(ctx context.Context, warmUpData *pb.WarmUpData) error {
	if warmUpData == nil {
		log.InfoContextf(ctx, "CheckWarmUpDataFormat, data empty,warmUpData:[%v]", warmUpData)
		return nil
	}
	if *warmUpData.Uint32WarmupUseType != WarmUpShowImg && *warmUpData.Uint32WarmupUseType != WarmUpShowVideo {
		err := errors.New("warmup_use_type invalid:" + strconv.Itoa(int(*warmUpData.Uint32WarmupUseType)))
		log.ErrorContextf(ctx, "warmup_use_type invalid, Uint32WarmupUseType :%v, ",
			warmUpData.Uint32WarmupUseType)
		return err
	}
	for _, val := range warmUpData.WarmupImgList {
		decodedCosId := string(GetBase64Decoded(ctx, *val.Url))
		if !IsValidCosId(ctx, decodedCosId) {
			err := errors.New("warmup img url cos id invalid, cosId:" + decodedCosId)
			log.ErrorContextf(ctx, "warmup img hit cos id format filter, raw cosId:%v", decodedCosId)
			return err
		}
	}
	for _, val := range warmUpData.WarmupVideoList {
		decodedCosId := string(GetBase64Decoded(ctx, *val.CosId))
		if !IsValidCosId(ctx, decodedCosId) {
			err := errors.New("warmup video url cos id invalid, cosId:" + decodedCosId)
			log.ErrorContextf(ctx, "warmup video hit cos id format filter, raw cosId:%v", decodedCosId)
			return err
		}
	}
	return nil
}

// FormatUintTime 将uint32时间戳转化为string
func FormatUintTime(timestamp uint32) string {
	formatTime := int64(timestamp)
	return fmt.Sprint(time.Unix(formatTime, 0).Format("2006-01-02 15:04:05"))
}

// MakeWarmUpInviteSwitchKey ...
func MakeWarmUpInviteSwitchKey(meetingId uint64) string {
	return fmt.Sprintf("warmup_invite_switch_%v", meetingId)
}

// MakeParticipantIncrKey ...
func MakeParticipantIncrKey(meetingId uint64) string {
	return fmt.Sprintf("webinar_participant_incr_%v", meetingId)
}

// MakeScheduleIncrKey ...
func MakeScheduleIncrKey(meetingId uint64) string {
	return fmt.Sprintf("webinar_schedule_incr_%v", meetingId)
}

// MakeParticipantMainKey ...
func MakeParticipantMainKey(meetingId uint64) string {
	return fmt.Sprintf("webinar_participant_main_%v", meetingId)
}

// MakeScheduleMainKey ...
func MakeScheduleMainKey(meetingId uint64) string {
	return fmt.Sprintf("webinar_schedule_main_%v", meetingId)
}

// Max ...
func Max(i, j uint32) uint32 {
	if i >= j {
		return i
	}
	return j
}

// Get32DaysExpireTimeDuration 获取会议过期时间戳
func Get32DaysExpireTimeDuration(meetingOrderEndTime uint32) uint32 {
	now := uint32(time.Now().Unix())
	expireDuation := uint32(MeetingExpireTime)
	if meetingOrderEndTime > now {
		expireDuation = meetingOrderEndTime + MeetingExpireTime - now
	}
	return expireDuation
}

// Get180DaysExpireTimeDuration 获取会议过期时间戳
func Get180DaysExpireTimeDuration(meetingOrderEndTime uint32) uint32 {
	now := uint32(time.Now().Unix())
	expireDuation := uint32(MeetingHalfAYearExpireTime)
	if meetingOrderEndTime > now {
		expireDuation = meetingOrderEndTime + MeetingHalfAYearExpireTime - now
	}
	return expireDuation
}

// GetNameSpace 获取NameSpace
func GetNameSpace() string {
	return trpc.GlobalConfig().Global.Namespace
}

// GetEnvName 获取EnvName
func GetEnvName() string {
	return trpc.GlobalConfig().Global.EnvName
}

// IsWebinarMeeting 判断是否是webinar会议
func IsWebinarMeeting(binaryMeetingType uint32) bool {
	return binaryMeetingType&0x01 == 1
}

// ContainsURL 判断是否包含网址
func ContainsURL(s string) bool {
	urlRegex := `http[s]?://[^\s]+|www\.[^\s]+`
	re := regexp.MustCompile(urlRegex)
	return re.MatchString(s)
}

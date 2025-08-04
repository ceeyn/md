package rpc

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"meeting_template/util"

	"git.code.oa.com/trpc-go/trpc-go/log"
	sensitivePb "git.code.oa.com/trpcprotocol/wemeet/wemeet_meet_sensitive"
	"github.com/golang/protobuf/proto"
	t_extractor "meeting_template/util/t-helper/t-html/t-extractor"
)

const SensitiveLabelPron = "Porn"     // 色情敏感词
const SensitiveLabelAd = "Ad"         // 广告敏感词
const SensitiveLabelAbuse = "Abuse"   // 谩骂敏感词
const SensitiveLabelPolity = "Polity" // 政治敏感词
const SuggestionBlock = "Block"
const SuggestionReview = "Review"
const BusinessFlag = "wemeet_template"

// 敏感词替换
func ReplaceSensitiveData(ctx context.Context, data string) (string, error) {
	// 一次校验最大5000字
	// 取出text进行校验
	bgTime := util.NowMs()
	var checkString strings.Builder
	var allUnAvailableString []string
	checkString.WriteString(data)
	for {
		if utf8.RuneCountInString(checkString.String()) >= 5000 {
			// 当前字符大于5000，截取前4999个字符发送rpc进行校验
			nowCheckString := []rune(checkString.String())
			unAvailableString, err := CheckDescriptionSensitiveData(ctx, string(nowCheckString[0:4999]))
			if err != nil {
				return "", err
			}
			allUnAvailableString = append(allUnAvailableString, unAvailableString...)
			// 剩余的字符放入checkString进行下一次校验
			checkString.Reset()
			checkString.WriteString(string(nowCheckString[4999:]))
		} else {
			// 当前字符小于5000，跳出循环
			break
		}
	}
	// 检测长度小于5000的剩余字符
	unAvailableString, err := CheckDescriptionSensitiveData(ctx, checkString.String())
	if err != nil {
		return "", err
	}
	allUnAvailableString = append(allUnAvailableString, unAvailableString...)
	// 敏感词替换
	for _, oldString := range allUnAvailableString {
		length := utf8.RuneCountInString(oldString)
		var newStringBuild strings.Builder
		for i := 0; i < length; i++ {
			newStringBuild.WriteString("*")
		}
		data = strings.ReplaceAll(data, oldString, newStringBuild.String())
	}
	cost := util.NowMs() - bgTime
	log.InfoContextf(ctx, "sensitive replace.allUnAvailableString:%v, cost:%vms",
		allUnAvailableString, cost)
	return data, nil
}

// 长文本校验敏感词
func CheckDescriptionSensitiveData(ctx context.Context, data string) ([]string, error) {
	util.ReportOne(util.GetTextModerationReq) //[CheckDescriptionSensitiveData]请求
	bgTime := util.NowMs()
	req := &sensitivePb.GetTextModerationReq{
		Content: proto.String(data),
		SysFlag: proto.String(BusinessFlag),
	}

	log.DebugContextf(ctx, "CheckDescriptionSensitiveData,content:%s,req:%v", data, req)
	//初始化请求代理
	proxy := sensitivePb.NewWemeetMeetSensitiveOidbClientProxy()
	//发送请求
	rsp, err := proxy.GetTextModeration(ctx, req)
	cost := util.NowMs() - bgTime
	log.InfoContextf(ctx, "CheckDescriptionSensitiveData,rsp:%v, cost:%vms", rsp, cost)
	if err != nil || rsp == nil || rsp.GetErrInfo() == nil {
		log.ErrorContextf(ctx, "check description sensitive failed, err is:%v", err)
		util.ReportOne(util.GetTextModerationFail) //[CheckDescriptionSensitiveData]请求失败
		return nil, errors.New("check description sensitive failed!")
	}
	if rsp.GetErrInfo().GetErrorCode() != 0 {
		log.ErrorContextf(ctx, "check idea sensitive failed, err code is:%v, err msg is:%v",
			rsp.GetErrInfo().GetErrorCode(), rsp.GetErrInfo().GetErrorMsg())
		util.ReportOne(util.GetTextModerationFail) //[CheckDescriptionSensitiveData]请求失败
		return nil, errors.New("check description sensitive failed, error code:" +
			string(rsp.GetErrInfo().GetErrorCode()))
	}
	var sensitiveKeyWords []string
	if rsp.GetKeyWordsDetail() != nil {
		for _, v := range rsp.GetKeyWordsDetail() {
			isSensitive := (v.GetSuggestion() == SuggestionBlock || v.GetSuggestion() == SuggestionReview) &&
				(v.GetLabel() == SensitiveLabelPron || v.GetLabel() == SensitiveLabelAd ||
					v.GetLabel() == SensitiveLabelAbuse || v.GetLabel() == SensitiveLabelPolity)
			if isSensitive {
				sensitiveKeyWords = append(sensitiveKeyWords, v.GetKeywords()...)
			}
		}
	}
	return sensitiveKeyWords, nil
}

// 长文本过信安敏感词检测
func CheckLongTextSensitiveData(ctx context.Context, data string, checkCosId bool, appId uint32,
	appUid, scenes string) (bool, string) {
	bgTime := util.NowMs()
	if data == "" {
		util.ReportOne(util.CheckTextEmpty)
		log.DebugContextf(ctx, "text empty, not hit sensitive")
		return false, ""
	}
	decodedDescription := util.GetBase64Decoded(ctx, data)
	if len(decodedDescription) == 0 {
		log.DebugContextf(ctx, "text base64 decocde failed, not hit sensitive")
		return false, ""
	}
	extractor := &t_extractor.TSimpleExtractor{}
	elements, err := extractor.Parse(string(decodedDescription))
	if err != nil {
		util.ReportOne(util.HtmlParseDescriptionFail)
		cost := util.NowMs() - bgTime
		log.ErrorContextf(ctx, "parse text failed, text: %+v,  error: %+v,not hit sensitive, cost:%vms",
			decodedDescription, err, cost)
		return false, ""
	}
	// cos_id不合法时，也应该判为检测不通过
	if checkCosId {
		imageCosIds := elements.Images
		for _, cosId := range imageCosIds {
			if !util.IsValidCosId(ctx, cosId) {
				util.ReportOne(util.DescriptionCosIdInvalid)
				cost := util.NowMs() - bgTime
				log.ErrorContextf(ctx, "invalid cos id: %+v, cost:%vms", cosId, cost)
				return true, "invalid cos id" + cosId
			}
		}
	}
	needCheckLongTexts := elements.TextFrags
	log.InfoContextf(ctx, "CheckLongTextSensitiveData() scenes:%v, elements.Images:%+v, elements.TextFrags:%+v",
		scenes, elements.Images, elements.TextFrags)
	text := strings.Join(needCheckLongTexts, " ")
	var checkString strings.Builder
	checkString.WriteString(text)
	for {
		if utf8.RuneCountInString(checkString.String()) >= 5000 {
			// 当前字符大于5000，截取前4999个字符发送rpc进行校验
			nowCheckString := []rune(checkString.String())
			hitSensitive, _ := CheckHasSensitiveWords(ctx, 0, appId, appUid, string(nowCheckString[0:4999]), scenes)
			if hitSensitive {
				cost := util.NowMs() - bgTime
				log.InfoContextf(ctx, "CheckLongTextSensitiveData() CheckHasSensitiveWords ,ret:%v, cost:%vms",
					hitSensitive, cost)
				return true, ""
			}
			// 剩余的字符放入checkString进行下一次校验
			checkString.Reset()
			checkString.WriteString(string(nowCheckString[4999:]))
		} else {
			// 当前字符小于5000，跳出循环
			break
		}
	}
	// 检测长度小于5000的剩余字符
	hitSensitive, _ := CheckHasSensitiveWords(ctx, 0, appId, appUid, checkString.String(), scenes)
	cost := util.NowMs() - bgTime
	log.InfoContextf(ctx, "CheckLongTextSensitiveData() CheckHasSensitiveWords ,ret:%v, cost:%vms", hitSensitive, cost)
	return hitSensitive, ""
}

// CheckSensitiveArray, 返回结果转换
func CheckSensitiveArray(data []string) (bool, string) {
	if len(data) != 0 {
		return true, strings.Join(data, ",")
	} else {
		return false, ""
	}
}

// CheckTextSensitiveData 过信安敏感词检测, 效果不符合预期
func CheckTextSensitiveData(ctx context.Context, data string, appId uint32, appUid string) bool {
	util.ReportOne(util.QueryMeetingSensitiveWords) //[CheckTextSensitive]请求
	startTime := time.Now()
	req := &sensitivePb.MeetingBillC2SQuerySensitiveWordReq{}
	req.RptStrWords = append(req.RptStrWords, data)
	req.SensitiveInfo = &sensitivePb.QueryMeetingSensitiveInfo{
		Uint32AppId: proto.Uint32(appId),
		StrAppUid:   proto.String(appUid),
		// 场景（1=昵称 2=入会昵称 3=会议主题), 这里约定传3
		Uint32Type: proto.Uint32(3),
		SysFlag:    proto.String(BusinessFlag),
		RequestId:  proto.String(BusinessFlag + strconv.FormatInt(util.Now(), 10)),
	}

	//初始化请求代理
	proxy := sensitivePb.NewWemeetMeetSensitiveOidbClientProxy()
	//发送请求
	rsp, err := proxy.QueryMeetingSensitiveWords(ctx, req)
	cost := time.Since(startTime).Milliseconds()
	log.DebugContextf(ctx, "CheckTextSensitiveData,content:%s,req:%v, time cost:%v, rsp:%v",
		data, req, cost, rsp)
	if err != nil || rsp == nil || len(rsp.GetRptSensitiveResult()) == 0 {
		util.ReportOne(util.QueryMeetingSensitiveWordsFail) //[CheckShortTextSensitive]请求失败
		log.ErrorContextf(ctx, "check short text sensitive failed, err is:%v", err)
		return false
	}
	if err != nil || rsp == nil || len(rsp.GetRptSensitiveResult()) == 0 {
		util.ReportOne(util.QueryMeetingSensitiveWordsFail) //[CheckShortTextSensitive]请求失败
		log.ErrorContextf(ctx, "check short text sensitive failed, rsp result size error:%v|%v", len(data),
			len(rsp.GetRptSensitiveResult()))
		return false
	}
	hitSensitive := false
	for idx, result := range rsp.GetRptSensitiveResult() {
		// 查询结果，0:非敏感词，1:敏感词
		if result.GetUint32ResultCode() == 1 {
			log.DebugContextf(ctx, "CheckTextSensitiveData got sensitive word:", data[idx])
			hitSensitive = true
			break
		}
	}
	return hitSensitive
}

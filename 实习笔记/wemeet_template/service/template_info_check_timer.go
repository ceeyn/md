package service

import (
	"context"
	"encoding/base64"

	"meeting_template/model"
	"meeting_template/rpc"
	"meeting_template/util"
	t_extractor "meeting_template/util/t-helper/t-html/t-extractor"
	t_replacer "meeting_template/util/t-helper/t-html/t-replacer"

	"git.code.oa.com/meettrpc/meet_util"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

// CheckTemplateInfoTimer
func CheckTemplateInfoTimer(ctx context.Context) error {
	defer meet_util.DefPanicFun()
	log.DebugContextf(ctx, "CheckTemplateInfoTimer start!")
	for {
		var needCheckIds []string
		// 获取需要检测的模板id列表，暂定一次获取1000个
		err := rpc.GetNeedCheckIds(ctx, &needCheckIds)
		if err != nil {
			util.ReportOne(util.RedisGetNeedCheckIdsFail)
			log.ErrorContextf(ctx, "get need check ids faild, err: %+v", err)
			break
		}
		log.DebugContextf(ctx, "CheckTemplateInfoTimer got ids, size:%v", len(needCheckIds))
		for _, id := range needCheckIds {
			templateInfo := &model.TemplateInfo{}
			templateInfo, err = GetTemplateInfoSingleFlight(ctx, id)
			if err != nil {
				log.ErrorContextf(ctx, "GetTemplateInfo failed, template id:%v, err:%v", id, err)
				continue
			}
			// 富文本为空时，不用检测，并将模板id从待检测有序集合里剔除
			if templateInfo.Description == "" {
				err := rpc.RemoveNeedCheck(ctx, id)
				log.DebugContextf(ctx, "decription empty, not need to check, id:%v, err:%v", id, err)
				continue
			}
			decodedDescription := util.GetBase64Decoded(ctx, templateInfo.Description)
			if len(decodedDescription) == 0 {
				err := rpc.RemoveNeedCheck(ctx, id) // 解码失败，说明数据格式有误，后续也不用再检测了
				log.DebugContextf(ctx, "decription decode error, not need to check, id:%v, err:%v", id, err)
				continue
			}
			extractor := &t_extractor.TSimpleExtractor{}
			elements, err := extractor.Parse(string(decodedDescription))
			if err != nil {
				util.ReportOne(util.HtmlParseDescriptionFail)
				log.ErrorContextf(ctx, "parse description failed, description: %+v,  error: %+v",
					decodedDescription, err)
				err := rpc.RemoveNeedCheck(ctx, id) // 解析html失败，说明数据格式有误，后续也不用再检测了
				log.DebugContextf(ctx, "parse description failed, not need to check, id:%v, err:%v", id, err)
				continue
			}
			needCheckLongTexts := elements.TextFrags
			var replacedLongTexts []string
			processed := true
			for _, text := range needCheckLongTexts {
				replaced, err := rpc.ReplaceSensitiveData(ctx, text)
				// 替换失败时，使用原数据, 不删除检测的模板id
				if err != nil {
					util.ReportOne(util.ReplaceDiscriptionSensitiveFail)
					log.ErrorContextf(ctx, "replace sensitive data failed, text: %+v,  error: %+v",
						text, err)
					replacedLongTexts = append(replacedLongTexts, text)
					processed = false
					continue
				}
				replacedLongTexts = append(replacedLongTexts, replaced)
			}
			replacer := &t_replacer.TSimpleReplacer{}
			description, err := replacer.ReplaceTextSrc(string(decodedDescription), replacedLongTexts)
			if err != nil {
				util.ReportOne(util.ReplaceDiscriptionTextFail)
				log.ErrorContextf(ctx, "replace description text failed, description: %+v,  error: %+v, "+
					"decodedDescription:%v", templateInfo.Description, err, decodedDescription)
				processed = false // 替换失败，下次重试，该次不删除检测的模板id
			} else {
				log.InfoContextf(ctx, "replace description images succeed, description: %+v", description)
				// 富文本为不完整的html文本，go的html包会补全，增加<html><head></head><body>前缀和</body></html>后缀，需要去除
				templateInfo.Description = base64.StdEncoding.EncodeToString(
					[]byte(description[25:(len(description) - 14)]))
			}
			err = SetTemplateInfo(ctx, templateInfo)
			log.ErrorContextf(ctx, "SetTemplateInfo  templateInfo:%+v, error:%v", templateInfo, err)
			if processed {
				err := rpc.RemoveNeedCheck(ctx, id)
				log.DebugContextf(ctx, "check completed, id:%v, err:%v", id, err)
			}
		}
		break
	}
	return nil
}

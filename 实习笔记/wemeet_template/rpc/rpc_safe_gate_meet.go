package rpc

import (
	"context"
	"errors"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	safePb "git.woa.com/wemeet-public/sdk-meet/safe/evt_ugc_config_join_welcome"
	wlcm "meeting_template/material_control/welcome"
	"strconv"
)

const (
	SAFE_ACTION   = "evt_ugc_config_join_welcome"
	SAFE_APP_FROM = "wemeet_template"
)

// BatchCheckHasSensitiveWords 批量校验敏感词，接口只返回命中安全策略的key
func BatchCheckHasSensitiveWords(ctx context.Context, meetingID uint64, appid uint32, appuid string, ip string, checkMap map[string]string) (map[string]string, error) {

	// 发送请求
	metrics.IncrCounter("Rpc.GetBatchSafe", 1)
	req := make([]*safePb.EvtUgcWbnJoinWelcome, 0)
	for key, word := range checkMap {
		checkItem := &safePb.EvtUgcWbnJoinWelcome{
			Action: SAFE_ACTION,
			User: &safePb.UserInfo{
				AppId:  appid,
				AppUid: appuid,
				Ip:     ip,
				Attrs: &safePb.UserAttrs{
					Text: word,
				},
			},
			Ser: &safePb.SerInfo{
				AppFrom: SAFE_APP_FROM,
				EchoStr: key,
			},
			Meet: &safePb.MeetingInfo{
				MeetingId: strconv.FormatUint(meetingID, 10),
			},
		}
		req = append(req, checkItem)
	}
	// 调用安全接口
	checkResult := make(map[string]string)
	// 失败重试一次
	for retry := 2; retry > 0; retry-- {
		rsp, err := safePb.GetBatchSafe(ctx, SAFE_ACTION, req)
		if err != nil {
			log.ErrorContextf(ctx, "GetBatchSafe failed, err: %+v", err.Error())
			continue
		}
		// 判断错误码
		if rsp.GetErrorCode() != 0 {
			log.ErrorContextf(ctx, "GetBatchSafe failed, errCode: %+v, errMsg: %+v", rsp.GetErrorCode(), rsp.GetErrorMsg())
			return nil, errors.New(rsp.GetErrorMsg())
		}
		log.InfoContextf(ctx, "GetBatchSafe rpc succ. meetingID:%+v, req:%+v, rsp:%+v", meetingID, req, rsp)
		for _, safeItem := range rsp.GetRsp() {
			if safeItem.GetResultCode() == 900 {
				// 命中安全的逻辑
				key := safeItem.GetSer().GetEchoStr()
				tips := safeItem.GetTips()
				if len(tips) == 0 {
					tips = wlcm.UPDATE_FAIL_DEFAULT_TIPS
				}
				checkResult[key] = tips
			}
		}
		break
	}
	log.InfoContextf(ctx, "BatchCheckHasSensitiveWords succ. meetingID:%+v, req:%+v, check_result:%+v", meetingID, req, checkResult)
	return checkResult, nil
}

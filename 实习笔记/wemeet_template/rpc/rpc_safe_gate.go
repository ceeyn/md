package rpc

import (
	"context"
	"errors"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	safe "git.code.oa.com/trpcprotocol/wemeet/wemeet_safe_gate"
	uuid "git.code.oa.com/wesee_ugc/go.uuid"
)

// CheckHasSensitiveWordsNew 新的接口：校验敏感词
func CheckHasSensitiveWordsNew(ctx context.Context, meetingId uint64, appId uint32, appUid string,
	word string, scenes string, action string) (bool, error) {

	req := &safe.GetSafeReq{
		Action: action,
		Meet: &safe.MeetingInfo{
			BinartMeetingType: 1, //webinar类型固定是1
			MeetingId:         meetingId,
		},
		FromUser: &safe.UserInfo{
			AppId:  appId,
			AppUid: appUid,
			//TinyId:     fmUser.GetUint64TinyId(),
			//InstanceId: fmUser.GetUint32InstanceId(),
			Attrs: map[string]string{
				"text": word,
			},
		},
		Ser: &safe.SerInfo{
			AppFrom: "wemeet_template",
			Scenes:  scenes,
			TraceId: uuid.NewV4().String(),
			//EchoStr: traceID,
		},
	}
	// 发送请求
	metrics.IncrCounter("Rpc.GetSafe", 1)
	proxy := safe.NewSafeTrpcClientProxy()
	rsp, err := proxy.GetSafe(ctx, req)
	if err != nil || rsp == nil || rsp.GetErrorCode() != 0 {
		metrics.IncrCounter("Rpc.GetSafe.Err", 1)
		log.ErrorContextf(ctx, "CheckHasSensitiveWordsNew template RpcGetSafe err, safeReq:%+v, rsp:%+v, err:%v",
			req, rsp, err)
		return false, errors.New("安全模块调用异常")
	}

	// 处理返回内容
	if err != nil {
		metrics.IncrCounter("Rpc.GetSafe.Err", 1) // [QueryMeetingSensitiveWords]查询敏感词信息失败
		log.ErrorContextf(ctx, "CheckHasSensitiveWordsNew template rpc failed. req:%+v,meetingId:%+v, err:%+v",
			req, meetingId, err)
		return false, err
	}
	metrics.IncrCounter("Rpc.GetSafe.Success", 1)

	//风险分，大于等于900，认为是命中恶意
	if rsp.GetScore() >= 900 {
		metrics.IncrCounter("Rpc.GetSafe.HighRisk", 1)
		log.ErrorContextf(ctx, "CheckHasSensitiveWordsNew template RpcGetSafe HighRisk！, safeReq:%v, rsp:%v", req, rsp)
		if rsp.GetTips() == "" {
			return true, errors.New("请检查修改内容后再重新保存")
		}
		return true, errors.New(rsp.GetTips())
	}
	//中风险
	if rsp.GetScore() >= 500 {
		metrics.IncrCounter("Rpc.GetSafe.MidRisk", 1) //中风险的放过(仅记录)
	}

	log.InfoContextf(ctx, "CheckHasSensitiveWordsNew-template rpc succ. meetingId:%+v,req:%+v, rsp:%+v",
		meetingId, req, rsp)

	return false, nil
}

// CheckHasSensitiveWordsReplace 替换旧的方法
func CheckHasSensitiveWordsReplace(ctx context.Context, meetingId uint64, appId uint32, appUid string,
	word string, oldScenes string) (bool, error) {

	if word == "" {
		return false, nil
	}

	//根据旧场景，确定新的action和scenes
	action, scenes := getActionAndSensesByOldScenes(oldScenes)
	if action == "" || scenes == "" {
		return false, errors.New("non register action or scenes")
	}

	return CheckHasSensitiveWordsNew(ctx, meetingId, appId, appUid, word, scenes, action)
}

// 不能新增。新的接入，请到信安@ethanqzheng,主动的申请
// getActionAndSensesByOldScenes 返回接入后的：接入分配的Action，接入分配的Scenes
func getActionAndSensesByOldScenes(scenes string) (string, string) {
	switch scenes {
	case "sc_webinar_template_sponsor":
		return "get_webinar_template_sponsor_check", "sc_webinar_template_sponsor"
	case "sc_webinar_template_desc":
		return "get_webinar_template_desc_check", "sc_webinar_template_desc"
	case "sc_meet_schedule_name":
		return "get_webinar_meet_schedule_name_check", "sc_meet_schedule_name"
	case "sc_meet_schedule_desc":
		return "get_webinar_meet_schedule_desc_check", "sc_meet_schedule_desc"
	case "sc_participant_name":
		return "get_webinar_participant_name_check", "sc_participant_name"
	case "sc_participant_profile":
		return "get_webinar_participant_profile_check", "sc_participant_profile"
	default:
		return "", ""
	}

}

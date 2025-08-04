package rpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"meeting_template/util"

	"git.code.oa.com/trpc-go/trpc-go/http"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
)

const cSerSendMsg = "trpc.http.FrontEnd.ProPlatform"

type nameBadgeResource struct {
	ResID         string   `json:"resId"`
	TaskContentID int64    `json:"taskContentId"`
	InstanceId    []uint32 `json:"instanceId"`
	Md5           string   `json:"md5"`
	Size          uint64   `json:"size"`
	Url           string   `json:"url"`
	Nonce         string   `json:"nonce"`
}

//ErrorInfo 错误信息
type ErrorInfo struct {
	ErrorMsg  string `json:"errorMsg"`
	ErrorCode int32  `json:"ret"` //0表示成功，-1表示失败
}

//PushNameBadgeZIPToProPlatformFrontEnd 推送名片样式到素材中心
func PushNameBadgeZIPToProPlatformFrontEnd(ctx context.Context, nameBadgeID string, taskContentID int64,
	size uint64, md5, url string) (int32, error) {

	//拼装请求信息
	req := &nameBadgeResource{
		ResID:         nameBadgeID,
		TaskContentID: taskContentID,
		InstanceId:    []uint32{1, 2, 3, 4, 6, 7, 8, 20, 21, 22},
		Md5:           md5,
		Size:          size,
		Url:           url,
		Nonce:         uuid.New().String(),
	}

	httpCli := http.NewClientProxy("client.http.material_system")
	path := fmt.Sprintf("/api/resource-update/save-default-cowork-resource")

	rsp := &ErrorInfo{}
	err := httpCli.Post(ctx, path, req, rsp)
	if err != nil {
		metrics.IncrCounter("save-default-cowork-resource fail", 1)
		log.ErrorContextf(ctx, "PushNameBadgeZIPToProPlatformFrontEnd fail, err:%+v", err)
		err = httpCli.Post(ctx, path, req, rsp)
		if err != nil {
			metrics.IncrCounter("save-default-cowork-resource again fail", 1)
			log.ErrorContextf(ctx, "PushNameBadgeZIPToProPlatformFrontEnd again fail, err:%+v, rsp:%+v", err, rsp)
			return util.ErrSendHttp, err
		}
	}
	if rsp.ErrorCode != 0 {
		metrics.IncrCounter("save-default-cowork-resource rsp fail", 1)
		log.ErrorContextf(ctx, "PushNameBadgeZIPToProPlatformFrontEnd rsp fail, req:%+v, rsp:%+v, err:%+v",
			req, rsp, err)
		return rsp.ErrorCode, errors.New(rsp.ErrorMsg)
	}
	metrics.IncrCounter("save-default-cowork-resource success", 1)
	log.InfoContextf(ctx, "PushNameBadgeZIPToProPlatformFrontEnd succ, req:%+v, rsp:%+v", req, rsp)
	return 0, nil
}

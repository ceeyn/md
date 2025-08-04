package rpc

import (
	"context"
	"errors"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	"meeting_template/util"

	commonUploadpb "git.code.oa.com/trpcprotocol/wemeet/common_upload"
)

const BusinessId = "wemeet_webinar"

// BatchQueryTempUrl rpc请求，获取分享信息
// commonUploadpb: http://trpc.rick.oa.com/rick/pb/view_protobuf?id=18780
func BatchQueryTempUrl(ctx context.Context,
	req *commonUploadpb.BatchTempUrlReq) (*commonUploadpb.BatchTempUrlRes, error) {
	bgTime := util.NowMs()
	util.ReportOne(util.BatchQueryTempUrl) //[BatchQueryTempUrl]请求
	proxy := commonUploadpb.NewCommonUploadOidbClientProxy()
	rsp, err := proxy.BatchQueryTempUrl(ctx, req)
	cost := util.NowMs() - bgTime
	if err != nil {
		util.ReportOne(util.BatchQueryTempUrlFail) //[BatchQueryTempUrl]rpc请求网络失败
		log.ErrorContextf(ctx, "rpc batch query temp url request fail, err: %v, req: %+v, rsp: %+v, cost:%vms",
			err, req, rsp, cost)
		return nil, err
	}
	// 错误码200为正常返回，其他为异常
	if rsp.GetCode() != 200 {
		util.ReportOne(util.BatchQueryTempUrlFail) //[BatchQueryTempUrl]rpc请求失败
		log.ErrorContextf(ctx, "rpc batch query temp url request fail, req: %+v, rsp: %+v, cost:%vms",
			req, rsp, cost)
		return nil, errors.New(rsp.GetMsg())
	}
	log.InfoContextf(ctx, "BatchQueryTempUrl ok  req:%v, rsp:%v, cost:%vms", req, rsp, cost)
	return rsp, nil
}

// JudgeCosResourceIsExist ...
func JudgeCosResourceIsExist(ctx context.Context, cosId string) bool {
    req := &commonUploadpb.QueryIsExistReq{
		BusinessId: BusinessId,
		CosId:      cosId,
	}
    proxy := commonUploadpb.NewCommonUploadTrpcClientProxy()
    rsp, err := proxy.TrpcQueryIsExist(ctx, req)
	if err != nil {
		log.ErrorContextf(ctx, "JudgeCosResourceIsExist rpc TrpcQueryIsExist failed, req:%+v, rsp:%+v, err:%+v",
			req, rsp, err)
		metrics.IncrCounter("TrpcQueryIsExist.failed", 1)
		return false
	}
	metrics.IncrCounter("TrpcQueryIsExist.succ",1)
	log.InfoContextf(ctx,"JudgeCosResourceIsExist rpc TrpcQueryIsExist succ, req:%+v, rsp:%+v",req, rsp.GetIsExist())
    return rsp.GetIsExist()
}
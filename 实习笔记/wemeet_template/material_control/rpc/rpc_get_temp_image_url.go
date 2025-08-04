package rpc

import (
	"context"
	"errors"

	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/wemeet/common_upload"
)

//BatchTempUrl 批量获取临时URL
func BatchTempUrl(ctx context.Context, cosIDs []string) (cosID2TempUrl map[string]string, err error) {
	proxy := pb.NewCommonUploadOidbClientProxy()
	BatchTempUrlReq := &pb.BatchTempUrlReq{
		CosIds:  cosIDs,
		Feature: "layout",
		IsCdn:   "0",
	}
	cosID2TempUrl = make(map[string]string)
	resp, err := proxy.BatchQueryTempUrl(ctx, BatchTempUrlReq)
	if err != nil || resp == nil {
		attr.AttrAPI(35732380, 1) //[BatchTempUrl]通过cosId获取tempUrl失败
		if err == nil {
			err = errors.New("empty resp")
		}
		log.ErrorContextf(ctx, "BatchTempUrl failed, err:%v, cosIDs:%v", err, cosIDs)
		return nil, err
	}
	if resp.GetData() != nil {
		cosID2TempUrl = resp.GetData().GetCosIdMap()
	}
	log.InfoContextf(ctx, "Batch Temp url finished, the data get is %v", cosID2TempUrl)
	return
}

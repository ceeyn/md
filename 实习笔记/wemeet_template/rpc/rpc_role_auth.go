package rpc

import (
	"context"
	"errors"
	"fmt"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	authPb "git.code.oa.com/trpcprotocol/tencent_meeting/common_service_role_auth"
	"github.com/golang/protobuf/proto"
)

// CheckIsSuperAdmin 校验是否是超管
func CheckIsSuperAdmin(ctx context.Context, appId, appUid string) (bool, error) {
	log.InfoContextf(ctx, "CheckIsAdmin appId: %+v appUid: %+v", appId, appUid)

	// 2、rpc 查询用户信息校验
	proxy := authPb.NewRoleAuthOidbClientProxy()

	req := &authPb.DescribeUserAuthReq{
		CorpId: proto.String(appId),
		Uid:    proto.String(appUid),
	}
	rsp, err := proxy.DescribeUserAuth(ctx, req)
	log.DebugContextf(ctx, "DescribeUserAuth rsp: %+v err: %+v", rsp, err)
	if err != nil {
		log.ErrorContextf(ctx, "DescribeUserAuth fail err: %+v", err)
		return false, errors.New("权限服务调用失败")
	}
	if rsp == nil {
		return false, nil
	}
	if rsp.GetRspCode() == nil || rsp.GetRspCode().GetCode() != "0" {
		log.ErrorContextf(ctx, "DescribeUserAuth rsp fail rsp: %+v", rsp)
		return false, nil
	}
	if rsp.Data == nil {
		return false, nil
	}
	if len(rsp.Data.GetAccessList()) > 0 {
		for _, authAccess := range rsp.Data.GetAccessList() {
			if authAccess.GetAuthCode() == "meeting-modification" {
				return true, nil
			}
		}
	}
	return false, nil
}

// QueryCorpInfo 查询企业信息
func QueryCorpInfo(ctx context.Context, appID uint32) (*authPb.DescribeCorpTagSingleReply, error) {
	req := &authPb.DescribeCorpTagSingleReq{
		CorpId: proto.String(fmt.Sprintf("%v", appID)),
	}
	rsp := &authPb.DescribeCorpTagSingleReply{}

	// 发送请求
	metrics.IncrCounter("QueryCorpInfo", 1)
	proxy := authPb.NewRoleAuthOidbClientProxy()
	rsp, err := proxy.DescribeCorpTag(ctx, req)

	// 处理返回内容
	if err != nil {
		metrics.IncrCounter("QueryCorpInfo err", 1) // [QueryCorpInfo]查询企业信息失败
		log.ErrorContextf(ctx, "(QueryUserDetail) rpc failed, appID:%+v, err:%+v", appID, err)
		return nil, err
	}
	if rsp.GetRspCode().GetCode() != "0" {
		log.ErrorContextf(ctx, "(QueryCorpTagInfo) Rsp code is not 0, appID: %+v, rsp: %+v", appID, rsp)
		return nil, errors.New(" Rsp code is not 0")
	}
	log.InfoContextf(ctx, "(QueryCorpTagInfo) success, req:%+v, rsp:%+v", req, rsp)

	return rsp, nil
}

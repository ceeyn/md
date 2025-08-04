package service

import (
	"context"
	"testing"

	trpc "git.code.oa.com/trpc-go/trpc-go"
	_ "git.code.oa.com/trpc-go/trpc-go/http"
	_ "git.code.oa.com/trpc-go/trpc-selector-cl5"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
)

var WemeetMeetingTemplateOidbService = &WemeetMeetingTemplateOidbServiceImpl{}

//go:generate mockgen -destination=stub/git.code.oa.com/trpcprotocol/wemeet/meeting_template/wemeet_meeting_template_oidb_mock.go -package=meeting_template -self_package=git.code.oa.com/trpcprotocol/wemeet/meeting_template git.code.oa.com/trpcprotocol/wemeet/meeting_template WemeetMeetingTemplateOidbClientProxy

func Test_WemeetMeetingTemplateOidb_CreateTemplate(t *testing.T) {

	// 开始写mock逻辑
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wemeetMeetingTemplateOidbClientProxy := pb.NewMockWemeetMeetingTemplateOidbClientProxy(ctrl)

	// 预期行为
	m := wemeetMeetingTemplateOidbClientProxy.EXPECT().CreateTemplate(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	m.DoAndReturn(func(ctx context.Context, req interface{}, opts ...interface{}) (interface{}, error) {

		r, ok := req.(*pb.CreateTemplateReq)
		if !ok {
			panic("invalid request")
		}

		rsp := &pb.CreateTemplateRsp{}
		err := WemeetMeetingTemplateOidbService.CreateTemplate(trpc.BackgroundContext(), r, rsp)
		return rsp, err
	})

	// 开始写单元测试逻辑
	req := &pb.CreateTemplateReq{}

	rsp, err := wemeetMeetingTemplateOidbClientProxy.CreateTemplate(trpc.BackgroundContext(), req)

	// 输出入参和返回 (检查t.Logf输出，运行 `go test -v`)
	t.Logf("WemeetMeetingTemplateOidb_CreateTemplate req: %v", req)
	t.Logf("WemeetMeetingTemplateOidb_CreateTemplate rsp: %v, err: %v", rsp, err)

	assert.Nil(t, err)
}

func Test_WemeetMeetingTemplateOidb_UpdateTemplate(t *testing.T) {

	// 开始写mock逻辑
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wemeetMeetingTemplateOidbClientProxy := pb.NewMockWemeetMeetingTemplateOidbClientProxy(ctrl)

	// 预期行为
	m := wemeetMeetingTemplateOidbClientProxy.EXPECT().UpdateTemplate(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	m.DoAndReturn(func(ctx context.Context, req interface{}, opts ...interface{}) (interface{}, error) {

		r, ok := req.(*pb.UpdateTemplateReq)
		if !ok {
			panic("invalid request")
		}

		rsp := &pb.UpdateTemplateRsp{}
		err := WemeetMeetingTemplateOidbService.UpdateTemplate(trpc.BackgroundContext(), r, rsp)
		return rsp, err
	})

	// 开始写单元测试逻辑
	req := &pb.UpdateTemplateReq{}

	rsp, err := wemeetMeetingTemplateOidbClientProxy.UpdateTemplate(trpc.BackgroundContext(), req)

	// 输出入参和返回 (检查t.Logf输出，运行 `go test -v`)
	t.Logf("WemeetMeetingTemplateOidb_UpdateTemplate req: %v", req)
	t.Logf("WemeetMeetingTemplateOidb_UpdateTemplate rsp: %v, err: %v", rsp, err)

	assert.Nil(t, err)
}

func Test_WemeetMeetingTemplateOidb_GetTemplateInfo(t *testing.T) {

	// 开始写mock逻辑
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wemeetMeetingTemplateOidbClientProxy := pb.NewMockWemeetMeetingTemplateOidbClientProxy(ctrl)

	// 预期行为
	m := wemeetMeetingTemplateOidbClientProxy.EXPECT().GetTemplateInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	m.DoAndReturn(func(ctx context.Context, req interface{}, opts ...interface{}) (interface{}, error) {

		r, ok := req.(*pb.GetTemplateInfoReq)
		if !ok {
			panic("invalid request")
		}

		rsp := &pb.GetTemplateInfoRsp{}
		err := WemeetMeetingTemplateOidbService.GetTemplateInfo(trpc.BackgroundContext(), r, rsp)
		return rsp, err
	})

	// 开始写单元测试逻辑
	req := &pb.GetTemplateInfoReq{}

	rsp, err := wemeetMeetingTemplateOidbClientProxy.GetTemplateInfo(trpc.BackgroundContext(), req)

	// 输出入参和返回 (检查t.Logf输出，运行 `go test -v`)
	t.Logf("WemeetMeetingTemplateOidb_GetTemplateInfo req: %v", req)
	t.Logf("WemeetMeetingTemplateOidb_GetTemplateInfo rsp: %v, err: %v", rsp, err)

	assert.Nil(t, err)
}

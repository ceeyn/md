package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"meeting_template/rpc"

	cachePb "git.code.oa.com/trpcprotocol/wemeet/common_meeting_cache"
	errpb "git.code.oa.com/trpcprotocol/wemeet/common_xcast_meeting_error_code"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestHandleUpdateMeetingShareItem(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UpdateMeetingShareItemReq
	}
	tests := []struct {
		name    string
		args    args
		wantRet int32
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "参数校验错误：FromUserInfo为空",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateMeetingShareItemReq{
					FromUserInfo: nil,
				},
			},
			wantRet: int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_INVALID_PARA),
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return true
			},
		},
		{
			name: "参数校验错误：GetMeetingId为0",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateMeetingShareItemReq{
					MeetingId:   proto.Uint64(0),
					ShareItemId: proto.String("2"),
					FromUserInfo: &pb.AppFromUserKey{
						AppUserId: proto.String("vicma"),
					},
				},
			},
			wantRet: int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_INVALID_PARA),
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return true
			},
		},
		{
			name: "获取会议信息错误",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateMeetingShareItemReq{
					MeetingId:   proto.Uint64(4794750316405402104),
					ShareItemId: proto.String("2"),
					FromUserInfo: &pb.AppFromUserKey{
						AppUserId:  proto.String("vicma"),
						TinyId:     proto.Uint64(144115355481710485),
						AppId:      proto.Uint32(1400143280),
						InstanceId: proto.Uint32(5),
					},
					ShowPassword: proto.Bool(true),
				},
			},
			wantRet: int32(errpb.ERROR_CODE_COMM_ERROR_CODE_COMM_MEETING_NOT_EXIST),
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return true
			},
		},
		{
			name: "设置会议信息错误",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateMeetingShareItemReq{
					MeetingId:   proto.Uint64(4794750316405402105),
					ShareItemId: proto.String("2"),
					FromUserInfo: &pb.AppFromUserKey{
						AppUserId:  proto.String("vicma"),
						TinyId:     proto.Uint64(144115355481710485),
						AppId:      proto.Uint32(1400143280),
						InstanceId: proto.Uint32(5),
					},
					ShowPassword: proto.Bool(true),
				},
			},
			wantRet: int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_SYSTEM_INNER_ERROR),
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return true
			},
		},
		{
			name: "设置showPassword错误",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateMeetingShareItemReq{
					MeetingId:   proto.Uint64(4794750316405402106),
					ShareItemId: proto.String("2"),
					FromUserInfo: &pb.AppFromUserKey{
						AppUserId:  proto.String("vicma"),
						TinyId:     proto.Uint64(144115355481710485),
						AppId:      proto.Uint32(1400143280),
						InstanceId: proto.Uint32(5),
					},
					ShowPassword: proto.Bool(true),
				},
			},
			wantRet: int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_SYSTEM_INNER_ERROR),
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return true
			},
		},
		{
			name: "成功",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateMeetingShareItemReq{
					MeetingId:   proto.Uint64(4794750316405402107),
					ShareItemId: proto.String("2"),
					FromUserInfo: &pb.AppFromUserKey{
						AppUserId:  proto.String("vicma"),
						TinyId:     proto.Uint64(144115355481710485),
						AppId:      proto.Uint32(1400143280),
						InstanceId: proto.Uint32(5),
					},
					ShowPassword: proto.Bool(true),
				},
			},
			wantRet: 0,
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return false
			},
		},
	}
	for _, tt := range tests {
		// mock rpc.GetMeetingInfo
		patchGetMeetingInfo := gomonkey.ApplyFunc(rpc.GetMeetingInfo,
			func(ctx context.Context, meetingId uint64) (*cachePb.MeetingInfo, bool, uint32, error) {
				// return error
				if meetingId == 4794750316405402104 {
					return nil, false, 0, errors.New("mock error")
				}
				return &cachePb.MeetingInfo{
					Uint64MeetingId:    proto.Uint64(meetingId),
					Uint32OrderEndTime: proto.Uint32(100),
				}, false, 0, nil
			})
		defer patchGetMeetingInfo.Reset()

		// mock rpc.SetMeetingInfo
		patchSetMeetingInfo := gomonkey.ApplyFunc(rpc.SetMeetingInfo,
			func(ctx context.Context, meetingInfo *cachePb.MeetingInfo, isOversea bool, cas uint32) error {
				// return error
				if meetingInfo.GetUint64MeetingId() == 4794750316405402105 {
					return errors.New("mock error")
				}
				return nil
			})
		defer patchSetMeetingInfo.Reset()

		// mock rpc.RDHSetShowPassword
		patchRDHSetShowPassword := gomonkey.ApplyFunc(rpc.RDHSetShowPassword,
			func(ctx context.Context, isOversea bool, meetingID uint64, orderEndTime uint32, showPassword bool) error {
				// return error
				if meetingID == 4794750316405402106 {
					return errors.New("mock error")
				}
				return nil
			})
		defer patchRDHSetShowPassword.Reset()

		t.Run(tt.name, func(t *testing.T) {
			gotRet, err := HandleUpdateMeetingShareItem(tt.args.ctx, tt.args.req)
			if !tt.wantErr(t, err, fmt.Sprintf("HandleUpdateMeetingShareItem(%v, %v)", tt.args.ctx, tt.args.req)) {
				return
			}
			assert.Equalf(t, tt.wantRet, gotRet, "HandleUpdateMeetingShareItem(%v, %v)", tt.args.ctx, tt.args.req)
		})
	}
}

package service

import (
	"context"
	"fmt"
	"testing"

	"meeting_template/config/config_rainbow"
	"meeting_template/rpc"
	"meeting_template/util"

	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestHandleAddScheduleData(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UpdateScheduleReq
		rsp *pb.UpdateScheduleRsp
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "req ScheduleList 为空检查",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateScheduleReq{},
				rsp: &pb.UpdateScheduleRsp{},
			},
			want: util.ERRLenIllegal,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
		{
			name: "req ScheduleList 长度不为1检查",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateScheduleReq{
					ScheduleList: []*pb.WebinarSchedule{
						&pb.WebinarSchedule{},
						&pb.WebinarSchedule{},
					},
				},
				rsp: &pb.UpdateScheduleRsp{},
			},
			want: util.ERRLenIllegal,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandleAddScheduleData(tt.args.ctx, tt.args.req, tt.args.rsp)
			if !tt.wantErr(t, err, fmt.Sprintf("HandleAddScheduleData(%v, %v, %v)", tt.args.ctx, tt.args.req, tt.args.rsp)) {
				return
			}
			assert.Equalf(t, tt.want, got, "HandleAddScheduleData(%v, %v, %v)", tt.args.ctx, tt.args.req, tt.args.rsp)
		})
	}
}

func TestHandleModifyScheduleData(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UpdateScheduleReq
		rsp *pb.UpdateScheduleRsp
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "req ScheduleList 为空检查",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateScheduleReq{},
				rsp: &pb.UpdateScheduleRsp{},
			},
			want: util.ERRLenIllegal,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
		{
			name: "req ScheduleList 长度不为1检查",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateScheduleReq{
					ScheduleList: []*pb.WebinarSchedule{
						&pb.WebinarSchedule{},
						&pb.WebinarSchedule{},
					},
				},
				rsp: &pb.UpdateScheduleRsp{},
			},
			want: util.ERRLenIllegal,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandleModifyScheduleData(tt.args.ctx, tt.args.req, tt.args.rsp)
			if !tt.wantErr(t, err, fmt.Sprintf("HandleModifyScheduleData(%v, %v, %v)", tt.args.ctx, tt.args.req, tt.args.rsp)) {
				return
			}
			assert.Equalf(t, tt.want, got, "HandleModifyScheduleData(%v, %v, %v)", tt.args.ctx, tt.args.req, tt.args.rsp)
		})
	}
}

func TestCheckScheduleParam(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UpdateScheduleReq
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "req ScheduleList 为空检查不通过",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateScheduleReq{},
			},
			want: util.ERRLenIllegal,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
		{
			name: "req ScheduleList 长度不为1检查不通过",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateScheduleReq{
					ScheduleList: []*pb.WebinarSchedule{
						&pb.WebinarSchedule{},
						&pb.WebinarSchedule{},
					},
				},
			},
			want: util.ERRLenIllegal,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
		{
			name: "req Uint64MeetingId 为0参数检查不通过",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateScheduleReq{
					ScheduleList: []*pb.WebinarSchedule{
						&pb.WebinarSchedule{},
					},
				},
			},
			want: util.ERRMeetingIdIllegal,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
		{
			name: "ScheduleName、ScheduleIndroduction 长度检查不通过",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateScheduleReq{
					ScheduleList: []*pb.WebinarSchedule{
						&pb.WebinarSchedule{
							ScheduleName:         proto.String("012345678901234567890"),
							ScheduleIndroduction: proto.String("123"),
						},
					},
					Uint64MeetingId: proto.Uint64(100),
				},
			},
			want: util.ERRLength,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
		{
			name: "ScheduleName 敏感词检查不通过",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateScheduleReq{
					ScheduleList: []*pb.WebinarSchedule{
						&pb.WebinarSchedule{
							ScheduleName:         proto.String("01234567890123456789"),
							ScheduleIndroduction: proto.String("123"),
						},
					},
					Uint64MeetingId: proto.Uint64(100),
					AppUserKey: &pb.AppFromUserKey{
						AppId:     proto.Uint32(0),
						AppUserId: proto.String("111"),
					},
				},
			},
			want: util.ERRNameSensitive,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
	}
	for _, tt := range tests {
		// mock rpc.RDHLenWebinarInfo
		patchGetMeetingInfo := gomonkey.ApplyFunc(rpc.RDHLenWebinarInfo,
			func(ctx context.Context, key string) (int64, error) {
				return 1, nil
			})
		defer patchGetMeetingInfo.Reset()

		// mock rpc.CheckHasSensitiveWords
		patchCheckHasSensitiveWords := gomonkey.ApplyFunc(rpc.CheckHasSensitiveWords,
			func(ctx context.Context, meetingId uint64, appId uint32, appUid string,
				word string, oldScenes string) (bool, error) {
				if oldScenes == SCScheduleName {
					return true, nil
				}
				return false, nil
			})
		defer patchCheckHasSensitiveWords.Reset()

		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckScheduleParam(tt.args.ctx, tt.args.req)
			if !tt.wantErr(t, err, fmt.Sprintf("CheckScheduleParam(%v, %v)", tt.args.ctx, tt.args.req)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CheckScheduleParam(%v, %v)", tt.args.ctx, tt.args.req)
		})
	}
}

func TestCheckScheduleParam1(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UpdateScheduleReq
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "schedule 数量校验不通过",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateScheduleReq{
					ScheduleList: []*pb.WebinarSchedule{
						&pb.WebinarSchedule{
							ScheduleName:         proto.String("01234567890123456789"),
							ScheduleIndroduction: proto.String("123"),
						},
					},
					Uint64MeetingId: proto.Uint64(100),
					AppUserKey: &pb.AppFromUserKey{
						AppId:     proto.Uint32(0),
						AppUserId: proto.String("111"),
					},
				},
			},
			want: util.ERRCount,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
	}
	for _, tt := range tests {
		// mock rpc.RDHLenWebinarInfo
		patchGetMeetingInfo := gomonkey.ApplyFunc(rpc.RDHLenWebinarInfo,
			func(ctx context.Context, key string) (int64, error) {
				return 20, nil
			})
		defer patchGetMeetingInfo.Reset()

		// mock config_rainbow.GetParticipantConfConfig
		patchGetParticipantConfConfig := gomonkey.ApplyFunc(config_rainbow.GetParticipantConfConfig,
			func() *config_rainbow.ParticipantConf {
				return &config_rainbow.ParticipantConf{
					ScheduleMaxCount: 10,
				}
			})
		defer patchGetParticipantConfConfig.Reset()

		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckScheduleParam(tt.args.ctx, tt.args.req)
			if !tt.wantErr(t, err, fmt.Sprintf("CheckScheduleParam(%v, %v)", tt.args.ctx, tt.args.req)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CheckScheduleParam(%v, %v)", tt.args.ctx, tt.args.req)
		})
	}
}

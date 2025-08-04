package service

import (
	"context"
	"fmt"
	"testing"

	"meeting_template/rpc"

	cachePb "git.code.oa.com/trpcprotocol/wemeet/common_meeting_cache"
	"github.com/agiledragon/gomonkey"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestGetTemplateInfoExpireTime(t *testing.T) {
	type args struct {
		ctx               context.Context
		meetingId         string
		isAfterMeetingEnd bool
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "",
			args: args{
				ctx:               context.Background(),
				meetingId:         "",
				isAfterMeetingEnd: false,
			},
			want: rpc.KeyTemplateExpireTime,
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return false
			},
		},
		{
			name: "",
			args: args{
				ctx:               context.Background(),
				meetingId:         "6812833127989168737",
				isAfterMeetingEnd: true,
			},
			want: rpc.KeyTemplateExpireTime,
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return false
			},
		},
	}
	for _, tt := range tests {
		// mock rpc.GetMeetingInfo
		patchGetMeetingInfo := gomonkey.ApplyFunc(rpc.GetMeetingInfo,
			func(ctx context.Context, meetingId uint64) (*cachePb.MeetingInfo, bool, uint32, error) {
				meetingInfo := &cachePb.MeetingInfo{
					Uint32OrderEndTime: proto.Uint32(1681110084), // 2023-04-10 15:01:24
				}
				return meetingInfo, false, 0, nil
			})
		defer patchGetMeetingInfo.Reset()

		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTemplateInfoExpireTime(tt.args.ctx, tt.args.meetingId, tt.args.isAfterMeetingEnd)
			fmt.Printf("got: %+v, err: %+v\n", got, err)
			if !tt.wantErr(t, err, fmt.Sprintf("GetTemplateInfoExpireTime(%v, %v, %v)", tt.args.ctx, tt.args.meetingId, tt.args.isAfterMeetingEnd)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetTemplateInfoExpireTime(%v, %v, %v)", tt.args.ctx, tt.args.meetingId, tt.args.isAfterMeetingEnd)
		})
	}
}

func TestCheckTemplateId(t *testing.T) {
	type args struct {
		templateId string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				templateId: "",
			},
			want: false,
		},
		{
			name: "",
			args: args{
				templateId: "tpl_1bc3a6d6-b45d-4bff-abf8-f3c84e6e0540",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				templateId: "t",
			},
			want: false,
		},
		{
			name: "",
			args: args{
				templateId: " tpl_1bc3a6d6-b45d-4bff-abf8-f3c84e6e0540",
			},
			want: false,
		},
		{
			name: "",
			args: args{
				templateId: "tpl_a0353d69-1743-4370-b160-5fbeefc9544e",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CheckTemplateId(tt.args.templateId), "CheckTemplateId(%v)", tt.args.templateId)
		})
	}
}

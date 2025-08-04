package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"meeting_template/config/config_rainbow"
	"meeting_template/material_control/background"
	"meeting_template/material_control/cache"
	"meeting_template/rpc"

	cachePb "git.code.oa.com/trpcprotocol/wemeet/common_meeting_cache"
	errpb "git.code.oa.com/trpcprotocol/wemeet/common_xcast_meeting_error_code"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/agiledragon/gomonkey"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func Test_handleQueryVirtualBackgroundList(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.QueryVirtualBackgroundListReq
		rsp *pb.QueryVirtualBackgroundListRsp
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "非cos安全校验请求，正常返回。" +
				"meetInfo.GetStrCreatorAppUid()==req.GetStrOperateAppuid() && " +
				"meetInfo.GetUint32CreatorSdkappid()==req.GetUint32OperateAppid()",
			args: args{
				ctx: context.Background(),
				req: &pb.QueryVirtualBackgroundListReq{
					Uint64MeetingId:    proto.Uint64(13242443280392425893),
					Uint32CarryCond:    proto.Uint32(0),
					StrOperateAppuid:   proto.String("100"),
					Uint32OperateAppid: proto.Uint32(200),
				},
				rsp: &pb.QueryVirtualBackgroundListRsp{
					ErrorCode:    proto.Int32(0),
					ErrorMessage: proto.String("ok"),
					MsgBackgroundInfo: []*pb.BackgroundInfo{
						&pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710139),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_1.png"),
							StrPicUrl:         proto.String("https://layout-1258344699.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						},
						&pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710140),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_2.png"),
							StrPicUrl:         proto.String("https://layout-61660710140.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						},
						&pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710141),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_2.png"),
							StrPicUrl:         proto.String("https://layout-61660710141.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						},
					},
				},
			},
			want: 0,
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return false
			},
		},
		{
			name: "非cos安全校验请求，正常返回。" +
				"meetInfo.GetStrCreatorAppUid()!=req.GetStrOperateAppuid() && hasPower== true",
			args: args{
				ctx: context.Background(),
				req: &pb.QueryVirtualBackgroundListReq{
					Uint64MeetingId:    proto.Uint64(13242443280392425893),
					Uint32CarryCond:    proto.Uint32(0),
					StrOperateAppuid:   proto.String("101"),
					Uint32OperateAppid: proto.Uint32(200),
				},
				rsp: &pb.QueryVirtualBackgroundListRsp{
					ErrorCode:    proto.Int32(0),
					ErrorMessage: proto.String("ok"),
					MsgBackgroundInfo: []*pb.BackgroundInfo{
						&pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710139),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_1.png"),
							StrPicUrl:         proto.String("https://layout-1258344699.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						},
						&pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710140),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_2.png"),
							StrPicUrl:         proto.String("https://layout-61660710140.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						},
						&pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710141),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_2.png"),
							StrPicUrl:         proto.String("https://layout-61660710141.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						},
					},
				},
			},
			want: 0,
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return false
			},
		},
		{
			name: "非cos安全校验请求，操作者权限校验不通过, " +
				"meetInfo.GetStrCreatorAppUid()!=req.GetStrOperateAppuid() 且 hasPower==false",
			args: args{
				ctx: context.Background(),
				req: &pb.QueryVirtualBackgroundListReq{
					Uint64MeetingId:    proto.Uint64(13242443280392425893),
					Uint32CarryCond:    proto.Uint32(0),
					StrOperateAppuid:   proto.String("102"),
					Uint32OperateAppid: proto.Uint32(200),
				},
				rsp: &pb.QueryVirtualBackgroundListRsp{},
			},
			want: int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_NO_PERMISSION),
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return true
			},
		},
		{
			name: "非cos安全校验请求，操作者权限校验不通过, " +
				"meetInfo.GetUint32CreatorSdkappid()!=req.GetUint32OperateAppid() 且 hasPower==false",
			args: args{
				ctx: context.Background(),
				req: &pb.QueryVirtualBackgroundListReq{
					Uint64MeetingId:    proto.Uint64(13242443280392425893),
					Uint32CarryCond:    proto.Uint32(0),
					StrOperateAppuid:   proto.String("102"),
					Uint32OperateAppid: proto.Uint32(201),
				},
				rsp: &pb.QueryVirtualBackgroundListRsp{},
			},
			want: int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_NO_PERMISSION),
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return true
			},
		},
	}
	for _, tt := range tests {
		// mock rpc.GetMeetingInfo
		patchGetMeetingInfo := gomonkey.ApplyFunc(rpc.GetMeetingInfo,
			func(ctx context.Context, meetingId uint64) (*cachePb.MeetingInfo, bool, uint32, error) {
				meetingInfo := &cachePb.MeetingInfo{
					Uint32OrderEndTime:    proto.Uint32(1681110084), // 2023-04-10 15:01:24
					StrCreatorAppUid:      proto.String("100"),
					Uint32CreatorSdkappid: proto.Uint32(200),
				}
				return meetingInfo, false, 0, nil
			})
		defer patchGetMeetingInfo.Reset()

		// mock rpc.CheckMeetingPermission
		patchCheckMeetingPermission := gomonkey.ApplyFunc(rpc.CheckMeetingPermission,
			func(ctx context.Context, appId uint32, appUid string, creatorAppId uint32) bool {
				if appUid == "101" {
					// appUid == 101，鉴权通过
					return true
				}
				return false
			})
		defer patchCheckMeetingPermission.Reset()

		// mock NewBackground
		patchNewBackground := gomonkey.ApplyFunc(background.NewBackground,
			func() cache.Backgrounds {
				return &background.BackgroundImp{}
			})
		defer patchNewBackground.Reset()

		// mock GetDefaultBackgroundSortSet
		var r *background.BackgroundImp
		patchGetDefaultBackgroundSortSet := gomonkey.ApplyMethod(reflect.TypeOf(r),
			"GetDefaultBackgroundSortSet",
			func(_ *background.BackgroundImp, ctx context.Context) (iDs []int64, err error) {
				iDs = append(iDs, 1)
				return iDs, nil
			})
		defer patchGetDefaultBackgroundSortSet.Reset()

		// mock GetMeetBackgroundSortSet
		patchGetMeetBackgroundSortSet := gomonkey.ApplyMethod(reflect.TypeOf(r),
			"GetMeetBackgroundSortSet",
			func(_ *background.BackgroundImp, ctx context.Context, meetingID uint64) (iDs []int64, err error) {
				iDs = append(iDs, 2)
				iDs = append(iDs, 3)
				return iDs, nil
			})
		defer patchGetMeetBackgroundSortSet.Reset()

		// mock GetBackground
		patchGetBackground := gomonkey.ApplyMethod(reflect.TypeOf(r), "GetBackground",
			func(_ *background.BackgroundImp, ctx context.Context,
				ids []int64) (background []*pb.BackgroundInfo, err error) {
				for _, id := range ids {
					if id == 1 {
						background = append(background, &pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710139),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_1.png"),
							StrPicUrl:         proto.String("https://layout-1258344699.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						})
					}
					if id == 2 {
						background = append(background, &pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710140),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_2.png"),
							StrPicUrl:         proto.String("https://layout-61660710140.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						})
					}
					if id == 3 {
						background = append(background, &pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710141),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_2.png"),
							StrPicUrl:         proto.String("https://layout-61660710141.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						})
					}
				}

				return background, nil
			})
		defer patchGetBackground.Reset()

		t.Run(tt.name, func(t *testing.T) {
			got, err := handleQueryVirtualBackgroundList(tt.args.ctx, tt.args.req, tt.args.rsp)
			if !tt.wantErr(t, err, fmt.Sprintf("handleQueryVirtualBackgroundList(%v, %v, %v)",
				tt.args.ctx, tt.args.req, tt.args.rsp)) {
				return
			}
			assert.Equalf(t, tt.want, got, "handleQueryVirtualBackgroundList(%v, %v, %v)",
				tt.args.ctx, tt.args.req, tt.args.rsp)
		})
	}
}

func Test_backgroundCosNotify(t *testing.T) {
	type args struct {
		ctx         context.Context
		meetingId   uint64
		backgrounds []*pb.BackgroundInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "开关关闭",
			args: args{
				ctx:         context.Background(),
				meetingId:   0,
				backgrounds: []*pb.BackgroundInfo{},
			},
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return false
			},
		},
	}
	for _, tt := range tests {
		// mock config_rainbow.GetBackgroundCosNotifySwitch
		patchGetBackgroundCosNotifySwitch := gomonkey.ApplyFunc(config_rainbow.GetBackgroundCosNotifySwitch,
			func() bool {
				return false
			})
		defer patchGetBackgroundCosNotifySwitch.Reset()

		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, backgroundCosNotify(tt.args.ctx, tt.args.meetingId, tt.args.backgrounds), fmt.Sprintf("backgroundCosNotify(%v, %v, %v)", tt.args.ctx, tt.args.meetingId, tt.args.backgrounds))
		})
	}
}

func Test_backgroundCosNotify1(t *testing.T) {
	type args struct {
		ctx         context.Context
		meetingId   uint64
		backgrounds []*pb.BackgroundInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "开关打开，rpc.GuestVbCosListNotify报错",
			args: args{
				ctx:         context.Background(),
				meetingId:   0,
				backgrounds: []*pb.BackgroundInfo{},
			},
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return true
			},
		},
	}
	for _, tt := range tests {
		// mock config_rainbow.GetBackgroundCosNotifySwitch
		patchGetBackgroundCosNotifySwitch := gomonkey.ApplyFunc(config_rainbow.GetBackgroundCosNotifySwitch,
			func() bool {
				return true
			})
		defer patchGetBackgroundCosNotifySwitch.Reset()

		// mock rpc.GuestVbCosListNotify
		patchGuestVbCosListNotify := gomonkey.ApplyFunc(rpc.GuestVbCosListNotify,
			func(ctx context.Context, meetingId uint64, backgrounds []*pb.BackgroundInfo) error {
				return errors.New("mock error")
			})
		defer patchGuestVbCosListNotify.Reset()

		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, backgroundCosNotify(tt.args.ctx, tt.args.meetingId, tt.args.backgrounds), fmt.Sprintf("backgroundCosNotify(%v, %v, %v)", tt.args.ctx, tt.args.meetingId, tt.args.backgrounds))
		})
	}
}

func Test_backgroundCosNotify2(t *testing.T) {
	type args struct {
		ctx         context.Context
		meetingId   uint64
		backgrounds []*pb.BackgroundInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "开关打开，rpc.GuestVbCosListNotify正常",
			args: args{
				ctx:         context.Background(),
				meetingId:   0,
				backgrounds: []*pb.BackgroundInfo{},
			},
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return true
			},
		},
	}
	for _, tt := range tests {
		// mock config_rainbow.GetBackgroundCosNotifySwitch
		patchGetBackgroundCosNotifySwitch := gomonkey.ApplyFunc(config_rainbow.GetBackgroundCosNotifySwitch,
			func() bool {
				return true
			})
		defer patchGetBackgroundCosNotifySwitch.Reset()

		// mock rpc.GuestVbCosListNotify
		patchGuestVbCosListNotify := gomonkey.ApplyFunc(rpc.GuestVbCosListNotify,
			func(ctx context.Context, meetingId uint64, backgrounds []*pb.BackgroundInfo) error {
				return nil
			})
		defer patchGuestVbCosListNotify.Reset()

		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, backgroundCosNotify(tt.args.ctx, tt.args.meetingId, tt.args.backgrounds), fmt.Sprintf("backgroundCosNotify(%v, %v, %v)", tt.args.ctx, tt.args.meetingId, tt.args.backgrounds))
		})
	}
}

func Test_handleGetVirtualBackgroundList(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.GetVirtualBackgroundListReq
		rsp *pb.GetVirtualBackgroundListRsp
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "其他服务请求，操作者权限校验不通过",
			args: args{
				ctx: context.Background(),
				req: &pb.GetVirtualBackgroundListReq{
					Uint64MeetingId: proto.Uint64(13242443280392425893),
					Uint32CarryCond: proto.Uint32(1), // 不查询默认背景图
				},
				rsp: &pb.GetVirtualBackgroundListRsp{},
			},
			want: int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_NO_PERMISSION),
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return true
			},
		},
		{
			name: "cos安全校验请求，正常返回",
			args: args{
				ctx: context.Background(),
				req: &pb.GetVirtualBackgroundListReq{
					Uint64MeetingId: proto.Uint64(13242443280392425893),
					Uint32CarryCond: proto.Uint32(1), // 不查询默认背景图
					StrAppFrom:      proto.String(AppFromWemeetLayoutCenter),
				},
				rsp: &pb.GetVirtualBackgroundListRsp{
					ErrorCode:    proto.Int32(0),
					ErrorMessage: proto.String("ok"),
					MsgBackgroundInfo: []*pb.BackgroundInfo{
						&pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710140),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_2.png"),
							StrPicUrl:         proto.String("https://layout-61660710140.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						},
						&pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710141),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_2.png"),
							StrPicUrl:         proto.String("https://layout-61660710141.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						},
					},
				},
			},
			want: 0,
			wantErr: func(assert.TestingT, error, ...interface{}) bool {
				return false
			},
		},
		{
			name: "cos安全校验请求，正常返回",
			args: args{
				ctx: context.Background(),
				req: &pb.GetVirtualBackgroundListReq{
					Uint64MeetingId: proto.Uint64(13242443280392425893),
					Uint32CarryCond: proto.Uint32(0), // 查询默认背景图
					StrAppFrom:      proto.String(AppFromWemeetLayoutCenter),
				},
				rsp: &pb.GetVirtualBackgroundListRsp{
					ErrorCode:    proto.Int32(0),
					ErrorMessage: proto.String("ok"),
					MsgBackgroundInfo: []*pb.BackgroundInfo{
						&pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710139),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_1.png"),
							StrPicUrl:         proto.String("https://layout-1258344699.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						},
						&pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710140),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_2.png"),
							StrPicUrl:         proto.String("https://layout-61660710140.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						},
						&pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710141),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_2.png"),
							StrPicUrl:         proto.String("https://layout-61660710141.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						},
					},
				},
			},
			want: 0,
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
					Uint32OrderEndTime:    proto.Uint32(1681110084), // 2023-04-10 15:01:24
					StrCreatorAppUid:      proto.String("100"),
					Uint32CreatorSdkappid: proto.Uint32(200),
				}
				return meetingInfo, false, 0, nil
			})
		defer patchGetMeetingInfo.Reset()

		// mock rpc.CheckMeetingPermission
		patchCheckMeetingPermission := gomonkey.ApplyFunc(rpc.CheckMeetingPermission,
			func(ctx context.Context, appId uint32, appUid string, creatorAppId uint32) bool {
				if appUid == "101" {
					// appUid == 101，鉴权通过
					return true
				}
				return false
			})
		defer patchCheckMeetingPermission.Reset()

		// mock NewBackground
		patchNewBackground := gomonkey.ApplyFunc(background.NewBackground,
			func() cache.Backgrounds {
				return &background.BackgroundImp{}
			})
		defer patchNewBackground.Reset()

		// mock GetDefaultBackgroundSortSet
		var r *background.BackgroundImp
		patchGetDefaultBackgroundSortSet := gomonkey.ApplyMethod(reflect.TypeOf(r),
			"GetDefaultBackgroundSortSet",
			func(_ *background.BackgroundImp, ctx context.Context) (iDs []int64, err error) {
				iDs = append(iDs, 1)
				return iDs, nil
			})
		defer patchGetDefaultBackgroundSortSet.Reset()

		// mock GetMeetBackgroundSortSet
		patchGetMeetBackgroundSortSet := gomonkey.ApplyMethod(reflect.TypeOf(r),
			"GetMeetBackgroundSortSet",
			func(_ *background.BackgroundImp, ctx context.Context, meetingID uint64) (iDs []int64, err error) {
				iDs = append(iDs, 2)
				iDs = append(iDs, 3)
				return iDs, nil
			})
		defer patchGetMeetBackgroundSortSet.Reset()

		// mock GetBackground
		patchGetBackground := gomonkey.ApplyMethod(reflect.TypeOf(r), "GetBackground",
			func(_ *background.BackgroundImp, ctx context.Context,
				ids []int64) (background []*pb.BackgroundInfo, err error) {
				for _, id := range ids {
					if id == 1 {
						background = append(background, &pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710139),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_1.png"),
							StrPicUrl:         proto.String("https://layout-1258344699.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						})
					}
					if id == 2 {
						background = append(background, &pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710140),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_2.png"),
							StrPicUrl:         proto.String("https://layout-61660710140.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						})
					}
					if id == 3 {
						background = append(background, &pb.BackgroundInfo{
							Int64BackgroundId: proto.Int64(61660710141),
							Uint32PicStatus:   proto.Uint32(4),
							StrPicId:          proto.String("background/default/default_2.png"),
							StrPicUrl:         proto.String("https://layout-61660710141.cos.accelerate.myqcloud.com/background/default/default_1.png?sign=707b8297677d9ae72f76a910de605064&t=1694589750"),
						})
					}
				}

				return background, nil
			})
		defer patchGetBackground.Reset()
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleGetVirtualBackgroundList(tt.args.ctx, tt.args.req, tt.args.rsp)
			if !tt.wantErr(t, err, fmt.Sprintf("handleGetVirtualBackgroundList(%v, %v, %v)", tt.args.ctx, tt.args.req, tt.args.rsp)) {
				return
			}
			assert.Equalf(t, tt.want, got, "handleGetVirtualBackgroundList(%v, %v, %v)", tt.args.ctx, tt.args.req, tt.args.rsp)
		})
	}
}

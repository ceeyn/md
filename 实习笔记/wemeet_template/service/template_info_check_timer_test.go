package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"meeting_template/model"
	"meeting_template/rpc"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestCheckTemplateInfoTimer(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "templateInfo.Description 为空时，rpc.RemoveNeedCheck 报错",
			args: args{
				ctx: context.Background(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
	}
	for _, tt := range tests {
		// mock rpc.GetNeedCheckIds
		patchGetMeetingInfo := gomonkey.ApplyFunc(rpc.GetNeedCheckIds,
			func(ctx context.Context, needCheckIds *[]string) error {
				*needCheckIds = append(*needCheckIds, "1")
				return nil
			})
		defer patchGetMeetingInfo.Reset()

		// mock GetTemplateInfoSingleFlight
		patchGetTemplateInfoSingleFlight := gomonkey.ApplyFunc(GetTemplateInfoSingleFlight,
			func(ctx context.Context, templateId string) (*model.TemplateInfo, error) {
				return &model.TemplateInfo{}, nil
			})
		defer patchGetTemplateInfoSingleFlight.Reset()

		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, CheckTemplateInfoTimer(tt.args.ctx), fmt.Sprintf("CheckTemplateInfoTimer(%v)", tt.args.ctx))
		})
	}
}

func TestCheckTemplateInfoTimer1(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "templateInfo.Description 为空时，rpc.RemoveNeedCheck 报错",
			args: args{
				ctx: context.Background(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
	}
	for _, tt := range tests {
		// mock rpc.GetNeedCheckIds
		patchGetMeetingInfo := gomonkey.ApplyFunc(rpc.GetNeedCheckIds,
			func(ctx context.Context, needCheckIds *[]string) error {
				*needCheckIds = append(*needCheckIds, "1")
				return nil
			})
		defer patchGetMeetingInfo.Reset()

		// mock GetTemplateInfoSingleFlight
		patchGetTemplateInfoSingleFlight := gomonkey.ApplyFunc(GetTemplateInfoSingleFlight,
			func(ctx context.Context, templateId string) (*model.TemplateInfo, error) {
				return &model.TemplateInfo{}, nil
			})
		defer patchGetTemplateInfoSingleFlight.Reset()

		// mock rpc.RemoveNeedCheck
		patchRemoveNeedCheck := gomonkey.ApplyFunc(rpc.RemoveNeedCheck,
			func(ctx context.Context, templateId string) error {
				return errors.New("test")
			})
		defer patchRemoveNeedCheck.Reset()
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, CheckTemplateInfoTimer(tt.args.ctx), fmt.Sprintf("CheckTemplateInfoTimer(%v)", tt.args.ctx))
		})
	}
}

func TestCheckTemplateInfoTimer2(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "templateInfo.Description base64 解码失败时 rpc.RemoveNeedCheck 报错",
			args: args{
				ctx: context.Background(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
	}
	for _, tt := range tests {
		// mock rpc.GetNeedCheckIds
		patchGetMeetingInfo := gomonkey.ApplyFunc(rpc.GetNeedCheckIds,
			func(ctx context.Context, needCheckIds *[]string) error {
				*needCheckIds = append(*needCheckIds, "1")
				return nil
			})
		defer patchGetMeetingInfo.Reset()

		// mock GetTemplateInfoSingleFlight
		patchGetTemplateInfoSingleFlight := gomonkey.ApplyFunc(GetTemplateInfoSingleFlight,
			func(ctx context.Context, templateId string) (*model.TemplateInfo, error) {
				return &model.TemplateInfo{
					Description: "123456",
				}, nil
			})
		defer patchGetTemplateInfoSingleFlight.Reset()

		// mock rpc.RemoveNeedCheck
		patchRemoveNeedCheck := gomonkey.ApplyFunc(rpc.RemoveNeedCheck,
			func(ctx context.Context, templateId string) error {
				return errors.New("test")
			})
		defer patchRemoveNeedCheck.Reset()
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, CheckTemplateInfoTimer(tt.args.ctx), fmt.Sprintf("CheckTemplateInfoTimer(%v)", tt.args.ctx))
		})
	}
}

func TestCheckTemplateInfoTimer3(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "rpc.ReplaceSensitiveData 报错, SetTemplateInfo 报错",
			args: args{
				ctx: context.Background(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
	}
	for _, tt := range tests {
		// mock rpc.GetNeedCheckIds
		patchGetMeetingInfo := gomonkey.ApplyFunc(rpc.GetNeedCheckIds,
			func(ctx context.Context, needCheckIds *[]string) error {
				*needCheckIds = append(*needCheckIds, "1")
				return nil
			})
		defer patchGetMeetingInfo.Reset()

		// mock GetTemplateInfoSingleFlight
		patchGetTemplateInfoSingleFlight := gomonkey.ApplyFunc(GetTemplateInfoSingleFlight,
			func(ctx context.Context, templateId string) (*model.TemplateInfo, error) {
				return &model.TemplateInfo{
					Description: "MTIzNDU=",
				}, nil
			})
		defer patchGetTemplateInfoSingleFlight.Reset()

		// mock rpc.ReplaceSensitiveData
		patchReplaceSensitiveData := gomonkey.ApplyFunc(rpc.ReplaceSensitiveData,
			func(ctx context.Context, data string) (string, error) {
				return "", errors.New("test")
			})
		defer patchReplaceSensitiveData.Reset()

		// mock SetTemplateInfo
		patchSetTemplateInfo := gomonkey.ApplyFunc(SetTemplateInfo,
			func(ctx context.Context, templateInfo *model.TemplateInfo) error {
				return nil
			})
		defer patchSetTemplateInfo.Reset()

		// mock rpc.RemoveNeedCheck
		patchRemoveNeedCheck := gomonkey.ApplyFunc(rpc.RemoveNeedCheck,
			func(ctx context.Context, templateId string) error {
				return errors.New("test")
			})
		defer patchRemoveNeedCheck.Reset()
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, CheckTemplateInfoTimer(tt.args.ctx), fmt.Sprintf("CheckTemplateInfoTimer(%v)", tt.args.ctx))
		})
	}
}

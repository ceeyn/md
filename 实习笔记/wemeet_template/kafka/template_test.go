package kafka

import (
	"context"
	"fmt"
	"git.code.oa.com/trpcprotocol/wemeet/wemee_db_proxy_center_join_wemeet_db_proxy_center_join"
	"github.com/golang/protobuf/proto"
	"meeting_template/model"
	"reflect"
	"testing"
)

func Test_getInsertTemplateInfoSql(t *testing.T) {
	type args struct {
		ctx          context.Context
		templateInfo model.TemplateInfo
	}
	tests := []struct {
		name  string
		args  args
		want  *string
		want1 []*wemee_db_proxy_center_join_wemeet_db_proxy_center_join.SqlTplParam
	}{
		{
			name: "",
			args: args{
				ctx: context.Background(),
				templateInfo: model.TemplateInfo{
					AppId:       "200000001",
					AppUid:      "144115217316437391",
					CoverList:   "",
					CoverName:   "",
					CoverUrl:    "",
					Description: "x�\u0001\u0000\u0000��\u0000\u0000\u0000\u0001",
					MeetingId:   "",
					Sponsor:     "",
					TemplateId:  "tpl_a0353d69-1743-4370-b160-5fbeefc9544e",
					WarmUpData:  "",
				},
			},
			want:  proto.String(""),
			want1: []*wemee_db_proxy_center_join_wemeet_db_proxy_center_join.SqlTplParam{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getInsertTemplateInfoSql(tt.args.ctx, tt.args.templateInfo)
			fmt.Printf("%+v\n", *got)
			fmt.Printf("%+v\n", got1)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getInsertTemplateInfoSql() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getInsertTemplateInfoSql() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getUpdateTemplateInfoSql(t *testing.T) {
	type args struct {
		ctx          context.Context
		templateInfo model.TemplateInfo
	}
	tests := []struct {
		name  string
		args  args
		want  *string
		want1 []*wemee_db_proxy_center_join_wemeet_db_proxy_center_join.SqlTplParam
	}{
		{
			name: "",
			args: args{
				ctx: context.Background(),
				templateInfo: model.TemplateInfo{
					AppId:       "200000001",
					AppUid:      "144115217316437391",
					CoverList:   "",
					CoverName:   "",
					CoverUrl:    "",
					Description: "x�\u0001\u0000\u0000��\u0000\u0000\u0000\u0001",
					MeetingId:   "",
					Sponsor:     "",
					TemplateId:  "tpl_a0353d69-1743-4370-b160-5fbeefc9544e",
					WarmUpData:  "",
				},
			},
			want:  proto.String(""),
			want1: []*wemee_db_proxy_center_join_wemeet_db_proxy_center_join.SqlTplParam{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getUpdateTemplateInfoSql(tt.args.ctx, tt.args.templateInfo)
			fmt.Println(*got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUpdateTemplateInfoSql() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getUpdateTemplateInfoSql() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

package util

import (
	"context"
	"testing"

	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/protobuf/proto"
)

func TestCheckCoverListFormat(t *testing.T) {
	type args struct {
		ctx        context.Context
		coverItems []*pb.CoverItem
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "RawUrl、CuttedUrl 为nil",
			args: args{
				ctx: context.Background(),
				coverItems: []*pb.CoverItem{
					&pb.CoverItem{},
				},
			},
			wantErr: true,
		},
		{
			name: "RawUrl、CuttedUrl 为空字符串",
			args: args{
				ctx: context.Background(),
				coverItems: []*pb.CoverItem{
					&pb.CoverItem{
						RawUrl:    proto.String(""),
						CuttedUrl: proto.String(""),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		// mock IsValidCosId
		patchIsValidCosId := gomonkey.ApplyFunc(IsValidCosId,
			func(ctx context.Context, cosId string) bool {
				return false
			})
		defer patchIsValidCosId.Reset()

		t.Run(tt.name, func(t *testing.T) {
			if err := CheckCoverListFormat(tt.args.ctx, tt.args.coverItems); (err != nil) != tt.wantErr {
				t.Errorf("CheckCoverListFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

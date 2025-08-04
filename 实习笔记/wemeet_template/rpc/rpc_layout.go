package rpc

import (
	"context"
	"fmt"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	"git.code.oa.com/trpcprotocol/wemeet/layout_center"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/golang/protobuf/proto"
)

// GuestVbCosListNotify 用于虚拟背景图片cos安全校验通知
// 注意：虚拟背景图片cos资源清理不依赖于该通知调用成功率 @thunderge
func GuestVbCosListNotify(ctx context.Context, meetingId uint64, backgrounds []*pb.BackgroundInfo) error {
	funcName := "GuestVbCosListNotify"
	baseLog := fmt.Sprintf("%+v(meetingId:%+v, backgrounds:%+v)", funcName, meetingId, backgrounds)

	req := &layout_center.GuestVbCosListNotifyReq{}
	strMeetingId := fmt.Sprintf("%v", meetingId)
	for _, background := range backgrounds {
		guestVbCosList := &layout_center.GuestVbCosList{
			StrMeetingId: proto.String(strMeetingId),
			StrCosId:     proto.String(background.GetStrPicId()),
			StrImageUrl:  proto.String(background.GetStrPicUrl()),
		}
		req.GuestVirtualBackgroundCosList = append(req.GuestVirtualBackgroundCosList, guestVbCosList)
	}

	proxy := layout_center.NewOidbLayoutCenterClientProxy()
	metrics.IncrCounter("rpc.guestVbCosListNotify.total", 1)
	rsp, err := proxy.GuestVbCosListNotify(ctx, req)
	log.DebugContextf(ctx, "%+v req:%+v, rsp:%+v, err:%+v", baseLog, req, rsp, err)
	if err != nil || rsp.GetErrorCode() != 0 {
		metrics.IncrCounter("rpc.guestVbCosListNotify.failed", 1)
		log.ErrorContextf(ctx, "%+v req:%+v, rsp:%+v, err:%+v", baseLog, req, rsp, err)
	}
	return err
}

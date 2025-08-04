package rpc

import (
	"context"
	"errors"
	"fmt"
	"git.code.oa.com/going/attr"
	"git.code.oa.com/meettrpc/meet_util"
	"git.code.oa.com/trpc-go/trpc-go/log"
	cachePb "git.code.oa.com/trpcprotocol/wemeet/common_meeting_cache"
	"github.com/golang/protobuf/proto"
)

var (
	NoExistError = errors.New("key not exist ") //grocery 数据不存在
)



//返回meetingInfo和国内还是海外
func gyGetMeetingInfoAllArea(ctx context.Context, meetingID uint64) (*cachePb.MeetingInfo, bool, uint32, error) {

	isOversea := false
	key := fmt.Sprintf("%v", meetingID)
	buf, cas, err := meet_util.GyGetValueCas(ctx, key, 0)
	if err == meet_util.NoExistError {
		//如果国内查不到，去海外查询
		buf, cas, err = meet_util.GyGetValueCas(ctx, key, 1)
		if err == meet_util.NoExistError {
			attr.AttrAPI(35928343, 1)//[template]获取海外数据为空
			return nil, isOversea, cas, err
		}
		isOversea = true
		//todo attr
	}
	if err != nil {
		//todo attr
		attr.AttrAPI(35928344, 1)//[template]获取国内数据出错
		log.ErrorContextf(ctx, "GyGetMeetingInfo from grocery fail, err: %v", err)
		return nil, isOversea, cas, err
	}

	mtInfo := &cachePb.MeetingInfo{}
	err = proto.Unmarshal(buf, mtInfo)
	if err != nil {
		attr.AttrAPI(35928345, 1)//[template]反序列化数据出错
		log.ErrorContextf(ctx, "GyGetMeetingInfo parse fail, err: %v", err)
		return nil, isOversea, cas, fmt.Errorf("parse fail key:%v err:%v ", key, err.Error())
	}

	attr.AttrAPI(35928346, 1)//[template]获取数据成功
	return mtInfo, isOversea, cas, nil
}


// GetMeetingInfo 获取会议信息
func GetMeetingInfo(ctx context.Context, meetingId uint64) (*cachePb.MeetingInfo, bool, uint32, error) {

	info, isOversea, cas, err := gyGetMeetingInfoAllArea(ctx, meetingId)
	if err != nil {
		attr.AttrAPI(35928347, 1)//[template]获取会议信息失败
		log.ErrorContextf(ctx, "meetingID=%v error: %v", meetingId, err.Error())
		return nil, isOversea, cas, err
	}
	attr.AttrAPI(35928348, 1)//[template]获取会议信息成功

	log.InfoContextf(ctx, "meetingID=%v district=%v, cas=%v, MeetingInfo=%+v", meetingId, isOversea, cas, info)
	return info, isOversea, cas, nil

}


// SetMeetingInfo 设置会议信息，目前只支持国内会议
func SetMeetingInfo(ctx context.Context, meetingInfo *cachePb.MeetingInfo, isOversea bool, cas uint32) error {

	if isOversea == false {
		key := fmt.Sprintf("%v", meetingInfo.GetUint64MeetingId())
		err := meet_util.GySetMeetingPbCas(ctx, key, meetingInfo, cas, 0)
		if err != nil {
			attr.AttrAPI(35928349, 1)//[template]设置会议信息失败
			log.ErrorContextf(ctx, "GySetMeetingPbCas from grocery fail, err: %v", err)
			return err
		}
		attr.AttrAPI(35928350, 1)//[template]设置会议信息成功
	}
	return nil
}

package service

import (
	"context"
	"encoding/base64"
	"errors"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	errpb "git.code.oa.com/trpcprotocol/wemeet/common_xcast_meeting_error_code"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	nb "meeting_template/material_control/namebadge"
	mcrpc "meeting_template/material_control/rpc"
	"meeting_template/rpc"
	"meeting_template/util"
	"path"
	"time"
)

const StorePath = "/tmp/"

//QueryNameBadgeInfoList 查询名牌样式列表
func (s *WemeetMeetingTemplateOidbServiceImpl) QueryNameBadgeInfoList(ctx context.Context,
	req *pb.QueryNamebadgeInfoListReq, rsp *pb.QueryNamebadgeInfoListRsp) error {
	metrics.IncrCounter("QueryNameBadgeInfoList Total", 1)
	start := time.Now()

	rst, err := handleQueryNameBadgeInfoList(ctx, req, rsp)
	if err != nil {
		metrics.IncrCounter("QueryNameBadgeInfoList fail", 1)
		rsp.ErrorCode = proto.Int32(rst)
		rsp.ErrorMessage = proto.String(err.Error())
	} else {
		metrics.IncrCounter("QueryNameBadgeInfoList Succ", 1)
		rsp.ErrorCode = proto.Int32(0)
		rsp.ErrorMessage = proto.String("ok")
	}
	log.InfoContextf(ctx, "QueryNameBadgeInfoList, cost:%v, rst:%v, err:%+v, req:%+v, rsp:%+v",
		time.Since(start), rst, err, req, rsp)
	return nil
}

//handleQueryNameBadgeInfoList ..
func handleQueryNameBadgeInfoList(ctx context.Context, req *pb.QueryNamebadgeInfoListReq,
	rsp *pb.QueryNamebadgeInfoListRsp) (int32, error) {

	//校验参数
	if req.GetUint64MeetingId() == 0 {
		metrics.IncrCounter("QueryNameBadgeInfoList invalid param", 1)
		return util.InvalidParam, errors.New("invalid param")
	}
	//操作者权限校验
	meetInfo, _, _, err := rpc.GetMeetingInfo(ctx, req.GetUint64MeetingId())
	if err != nil {
		metrics.IncrCounter("rpc.GetMeetingInfo fail", 1)
		log.ErrorContextf(ctx, "rpc.GetMeetingInfo fail, meetID:%v, err:%+v", req.GetUint64MeetingId(), err)
		return util.ERRQueryMeetInfo, err
	}
	//判断操作者身份
	if meetInfo.GetStrCreatorAppUid() != req.GetStrOperateAppuid() ||
		meetInfo.GetUint32CreatorSdkappid() != req.GetUint32OperateAppid() {
		hasPower := rpc.CheckMeetingPermission(ctx, req.GetUint32OperateAppid(), req.GetStrOperateAppuid(),
			meetInfo.GetUint32CreatorSdkappid())
		if !hasPower {
			metrics.IncrCounter("QueryNameBadgeInfoList Not Permission", 1)
			return int32(errpb.ERROR_CODE_MEETING_LOGIC_WEBINAR_ERROR_CODE_MEETING_LOGIC_WEBINAR_NO_PERMISSION),
				errors.New("not permission")
		}
	}

	nameBadge := nb.NewNameBadge()
	nameBadgeList, err := nameBadge.QueryDefNameBadgeInfoList(ctx)
	if err != nil {
		metrics.IncrCounter("QueryDefNameBadgeInfoList fail", 1)
		log.ErrorContextf(ctx, "QueryDefNameBadgeInfoList fail, meetID:%v, err:%+v",
			req.GetUint64MeetingId(), err)
		return util.ErrQueryCacheData, err
	}

	rsp.NamebadgeInfo = nameBadgeList
	return 0, nil
}

//SetNameBadgeInfoList 设置名牌信息
func (s *WemeetMeetingTemplateOidbServiceImpl) SetNameBadgeInfoList(ctx context.Context,
	req *pb.SetNameBadgeInfoListReq, rsp *pb.SetNameBadgeInfoListRsp) error {
	start := time.Now()

	err := handleNameBadgeConfConfig(ctx, req)
	if err != nil {
		metrics.IncrCounter("SetNameBadgeInfoList fail", 1)
		rsp.ErrorCode = proto.Int32(1)
		rsp.ErrorMessage = proto.String(err.Error())
	} else {
		metrics.IncrCounter("SetNameBadgeInfoList Succ", 1)
		rsp.ErrorCode = proto.Int32(0)
		rsp.ErrorMessage = proto.String("")
	}

	log.InfoContextf(ctx, "SetNameBadgeInfoList, cost:%v, err:%+v, req:%+v, rsp:%+v",
		time.Since(start), err, req, rsp)
	return nil
}

func handleNameBadgeConfConfig(ctx context.Context, req *pb.SetNameBadgeInfoListReq) error {

	//推送到素材系统
	//todo:异步
	tasks := make([]func() error, 0, len(req.GetMsgNamebadgeList()))
	for _, val := range req.GetMsgNamebadgeList() {
		nameBadgeInfo := val
		tasks = append(tasks, func() error {
			//推素材系统
			rst, err := mcrpc.PushNameBadgeZIPToProPlatformFrontEnd(ctx, nameBadgeInfo.GetStrNamebadgeId(),
				nameBadgeInfo.GetInt64TaskContentId(), nameBadgeInfo.GetUint64NamebadgeSize(),
				nameBadgeInfo.GetStrNamebadgeMd5(), nameBadgeInfo.GetStrNamebadgeUrl())
			if rst != 0 {
				log.ErrorContextf(ctx, "PushNameBadgeZIPToFrontEnd fail, rst:%v, err:%+v, nameBadgeInfo:%+v",
					rst, err, nameBadgeInfo)
				rst, err = mcrpc.PushNameBadgeZIPToProPlatformFrontEnd(ctx, nameBadgeInfo.GetStrNamebadgeId(),
					nameBadgeInfo.GetInt64TaskContentId(), nameBadgeInfo.GetUint64NamebadgeSize(),
					nameBadgeInfo.GetStrNamebadgeMd5(), nameBadgeInfo.GetStrNamebadgeUrl())
				if rst != 0 {
					metrics.IncrCounter("PushNameBadgeZIPToFrontEnd again fail", 1)
					log.ErrorContextf(ctx, "PushNameBadgeZIPToFrontEnd again fail, rst:%v, err:%+v, "+
						"nameBadgeInfo:%+v", rst, err, nameBadgeInfo)
					return err
				}
			}
			data, err := nameBadgeZipDeal(ctx, nameBadgeInfo)
			if err != nil {
				metrics.IncrCounter("nameBadgeZipDeal fail", 1)
				log.ErrorContextf(ctx, "nameBadgeZipDeal fail, nameBadgeInfo:%+v, err:%+v", nameBadgeInfo, err)
				return err
			}
			//将样式信息写入到Redis中
			nbInit := nb.NewNameBadge()
			infoList := []*pb.NameBadgeInfo{data}
			err = nbInit.SetDefNameBadgeList(ctx, infoList)
			if err != nil {
				metrics.IncrCounter("SetDefNameBadgeList fail", 1)
				log.ErrorContextf(ctx, "SetDefNameBadgeList fail, err:%+v", err)
				return err
			}
			ids := []string{nameBadgeInfo.GetStrNamebadgeId()}
			err = nbInit.SetDefNameBadgeSortSet(ctx, ids)
			if err != nil {
				metrics.IncrCounter("SetDefNameBadgeSortSet fail", 1)
				log.ErrorContextf(ctx, "SetDefNameBadgeSortSet fail, ids:%+v,err:%+v", ids, err)
				return err
			}

			return nil
		})
	}
	err := trpc.GoAndWait(tasks...)
	if err != nil {
		metrics.IncrCounter("handleNameBadgeConfConfig fail", 1)
		log.ErrorContextf(ctx, "handleNameBadgeConfConfig, err:%+v", err)
		return err
	}
	return nil
}

//nameBadgeZipDeal ..
func nameBadgeZipDeal(ctx context.Context, nameBadgeInfo *pb.SetNameBadgeInfo) (*pb.NameBadgeInfo, error) {
	//1.获取压缩包
	filePath, err := mcrpc.DownLoadingFileForHttpUrl(ctx, nameBadgeInfo.GetStrNamebadgeName(),
		path.Join(StorePath, nameBadgeInfo.GetStrNamebadgeName()), nameBadgeInfo.GetStrNamebadgeUrl())
	if err != nil {
		metrics.IncrCounter("mcrpc.DownLoadingFileForHttpUrl fail", 1)
		log.ErrorContextf(ctx, "mcrpc.DownLoadingFileForHttpUrl fail, nameBadgeInfo:%+v, err:%+v",
			nameBadgeInfo, err)
		return nil, err
	}
	//2.解析zip包
	err = mcrpc.Decompress(filePath, StorePath)
	if err != nil {
		metrics.IncrCounter("mcrpc.Decompress fail", 1)
		log.ErrorContextf(ctx, "mcrpc.Decompress fail, err :%+v", err)
		return nil, err
	}
	//判断是否有图片及获取图片名称
	ImageList := make([]*pb.NameBadgeInfo_NameBadgeImage, 0)
	isExist, _ := nb.PathExists(path.Join(StorePath, nameBadgeInfo.GetStrNamebadgeName(), "/images"))
	if isExist {
		fileNames := nb.GetDirFileName(path.Join(StorePath, nameBadgeInfo.GetStrNamebadgeName(), "/images"))
		for _, fileName := range fileNames {
			//将图片存储到cos上
			fileUrl, err := mcrpc.UploadingFile(ctx, fileName,
				path.Join(StorePath, nameBadgeInfo.GetStrNamebadgeName(), "/images/"))
			if err != nil {
				metrics.IncrCounter("mcrpc.UploadingFile fail", 1)
				log.ErrorContextf(ctx, "mcrpc.UploadingFile fail, nameBadgeInfo:%+v, err:%+v",
					nameBadgeInfo, err)
				return nil, err
			}
			imageInfo := &pb.NameBadgeInfo_NameBadgeImage{
				StrImageName: proto.String(fileName),
				StrImageUrl:  proto.String(fileUrl),
			}
			ImageList = append(ImageList, imageInfo)
		}
	}
	//获取样式json数据
	content, err := ioutil.ReadFile(path.Join(StorePath, nameBadgeInfo.GetStrNamebadgeName(), "style.json"))
	if err != nil || len(content) == 0 {
		metrics.IncrCounter("ioutil.ReadFile style.json fail", 1)
		log.ErrorContextf(ctx, " ioutil.ReadFile style.json fail, err:%+v, len(content):%v",
			err, len(content))
		return nil, err
	}

	rst := &pb.NameBadgeInfo{
		StrNamebadgeId:      proto.String(nameBadgeInfo.GetStrNamebadgeId()),
		StrNamebadgeContent: proto.String(base64.StdEncoding.EncodeToString(content)),
		NamebadgeImage:      ImageList,
	}

	return rst, nil
}

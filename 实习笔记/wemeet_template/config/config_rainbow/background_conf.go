package config_rainbow

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"time"

	"git.code.oa.com/going/attr"
	"git.code.oa.com/meettrpc/meet_util"
	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/golang/protobuf/proto"
	bg "meeting_template/material_control/background"
)

// BackgroundConf ...
type BackgroundConf struct {
	BackgroundId int64  `json:"background_id"`
	PicId        string `json:"pic_id"`
	PicUrl       string `json:"pic_url"`
	PicDesc      string `json:"pic_desc"`
}

//BackgroundConfList ..
type BackgroundConfList struct {
	BackgroundSize            int64             `json:"background_list_size"`
	BackgroundAuditTime       int64             `json:"background_audit_time"`
	BackgroundCosNotifySwitch bool              `json:"background_cos_notify_switch"`
	BackgroundInfo            []*BackgroundConf `json:"background_list"`
}

var (
	backgroundCfgValue atomic.Value
)

func init() {
	backgroundCfgValue.Store(&BackgroundConfList{
		BackgroundSize:            10,
		BackgroundAuditTime:       300,
		BackgroundCosNotifySwitch: false,
	})
}

//GetBackgroundConfSize ..
func GetBackgroundConfSize() int64 {
	cfg := backgroundCfgValue.Load().(*BackgroundConfList)
	log.Infof("GetBackgroundConfSize:%+v", cfg.BackgroundSize)
	return cfg.BackgroundSize
}

//GetBackgroundConfAuditTime ..
func GetBackgroundConfAuditTime() int64 {
	cfg := backgroundCfgValue.Load().(*BackgroundConfList)
	log.Infof("GetBackgroundConfSize:%+v", cfg.BackgroundAuditTime)
	return cfg.BackgroundAuditTime
}

//GetBackgroundCosNotifySwitch 获取虚拟背景图片上传，cos信息通知layout开关
func GetBackgroundCosNotifySwitch() bool {
	cfg := backgroundCfgValue.Load().(*BackgroundConfList)
	log.Infof("BackgroundCosNotifySwitch:%+v", cfg.BackgroundCosNotifySwitch)
	return cfg.BackgroundCosNotifySwitch
}

//CompareIDSliceEqual 判断两个ID切片是否相同
func CompareIDSliceEqual(src, dst []int64) bool {
	if len(src) != len(dst) {
		return false
	}
	//长度相同，只需要判断src在dst中是否存在
	srcMap := make(map[int64]struct{}, 0)
	for _, id := range src {
		srcMap[id] = struct{}{}
	}
	for _, id := range dst {
		if _, ok := srcMap[id]; !ok {
			return false
		}
	}
	return true
}

//HandleBackgroundConfConfig 获取七彩石配置
func HandleBackgroundConfConfig(data string) error {

	cfg := &BackgroundConfList{}
	err := json.Unmarshal([]byte(data), cfg)
	if err != nil {
		attr.AttrAPI(35927486, 1)
		log.Errorf("HandleBackgroundConfConfig json Unmarshal error, err:%v", err)
		return err
	}
	log.Infof("HandleBackgroundConfConfig, cfg:%+v", cfg)
	backgroundCfgValue.Store(cfg)

	go func() {
		defer meet_util.DefPanicFun()
		time.Sleep(time.Second * 180)
		handleBackgroundConfConfig(cfg)
	}()
	return nil
}

//HandleBackgroundConfConfig ..
func handleBackgroundConfConfig(cfg *BackgroundConfList) error {

	rainbowCfgIdList := make([]int64, 0)
	backgroundList := make([]*pb.BackgroundInfo, 0, len(cfg.BackgroundInfo))
	for _, val := range cfg.BackgroundInfo {
		backgroundInfo := &pb.BackgroundInfo{
			Int64BackgroundId: proto.Int64(val.BackgroundId),
			Uint32PicStatus:   proto.Uint32(bg.PicStatusDefault), //图片状态为默认，
			StrPicId:          proto.String(val.PicId),           //图片的cosId
			StrPicUrl:         proto.String(val.PicUrl),
			StrPicDesc:        proto.String(val.PicDesc),
			Uint64IncrTime:    proto.Uint64(uint64(time.Now().Unix())),
		}
		backgroundList = append(backgroundList, backgroundInfo)
		rainbowCfgIdList = append(rainbowCfgIdList, val.BackgroundId)
	}

	ctx := context.Background()
	bgInit := bg.NewBackground()

	CurIds, err := bgInit.GetDefaultBackgroundSortSet(ctx)
	if err == nil {
		isEqual := CompareIDSliceEqual(rainbowCfgIdList, CurIds)
		if isEqual {
			return nil
		}
	}

	err = bgInit.SetDefBackgroundListInfo(ctx, backgroundList)
	if err != nil {
		log.Errorf("SetDefBackgroundListInfo fail, err:%+v", err)
	}

	return nil
}

package config_rainbow

import (
	"encoding/json"
	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"sync/atomic"
)

// CallBackConf ...
type CallBackConf struct {
	GateWayImgBackgroundActon      bool `json:"gate_way_img_background_acton"`
	GateWayImgScenesAction         bool `json:"gate_way_img_scenes_action"`
	SCParticipantAvatarAction      bool `json:"sc_participant_avatar_action"`
	WarmUpShowVideoAction         bool `json:"warm_up_show_video_action"`
}

var (
	callBackConf atomic.Value
)

func init() {
	callBackConf.Store(&CallBackConf{
		GateWayImgBackgroundActon:false,
		GateWayImgScenesAction: false,
		SCParticipantAvatarAction:false,
		WarmUpShowVideoAction:false,
	})
}

//GetCallBackConf ..
func GetCallBackConf() *CallBackConf {
	attr.AttrAPI(36337052, 1)
	cfg := callBackConf.Load().(*CallBackConf)
	log.Infof("GetCallBackConf:%+v", cfg)
	return cfg
}


//HandleCallBackConf ...
func HandleCallBackConf(data string) error {
	cfg := &CallBackConf{}
	err := json.Unmarshal([]byte(data), cfg)
	if err != nil {
		attr.AttrAPI(36337053,1)
		log.Errorf("HandleCallBackConf json Unmarshal error, err:%v", err)
		return err
	}
	log.Infof("HandleCallBackConf, cfg:%+v", cfg)
	callBackConf.Store(cfg)
	return nil
}

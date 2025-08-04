package config_rainbow // NOCA:golint/package(设计如此)

import (
	"encoding/json"
	"sync/atomic"

	"git.code.oa.com/trpc-go/trpc-go/log"
)

// CommonConf 通用配置
type CommonConf struct {
	AdminAppIdWhiteList []string `json:"admin_app_id_white_list"` // 需要判断超管权限的白名单企业
}

var (
	commonCfgValue atomic.Value
)

//HandleCommonConfig 获取七彩石配置
func HandleCommonConfig(data string) error {

	cfg := &CommonConf{}
	err := json.Unmarshal([]byte(data), cfg)
	if err != nil {
		log.Errorf("HandleCommonConfig json Unmarshal error, err:%v", err)
		return err
	}
	log.Infof("HandleCommonConfig, cfg:%+v", cfg)
	commonCfgValue.Store(cfg)

	return nil
}

//GetCommonConfig ..
func GetCommonConfig() *CommonConf {
	cfg := commonCfgValue.Load().(*CommonConf)
	log.Infof("GetCommonConfig:%+v", cfg)
	return cfg
}

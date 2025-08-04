package config_rainbow

import (
	"encoding/json"
	"sync/atomic"

	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

// ESConf ...
type ESConf struct {
	EsHost   string `json:"es_host"`
	EsUser   string `json:"es_user"`
	EsPwd    string `json:"es_pwd"`
	EsSwitch string `json:"es_switch"`
}

var (
	esConfCfgValue atomic.Value
)

func init() {
	esConfCfgValue.Store(&ESConf{
		EsHost:   "http://9.144.32.31:9200", //默认
		EsUser:   "elastic",                 //默认
		EsPwd:    "wemeet2020!",
		EsSwitch: "open",
	})
}

//GetEsConfConfig ..
func GetEsConfConfig() *ESConf {
	attr.AttrAPI(36337052, 1)
	cfg := esConfCfgValue.Load().(*ESConf)
	log.Infof("GetEsConfConfig:%+v", cfg)
	return cfg
}

//HandleEsConfig ...
func HandleEsConfig(data string) error {
	cfg := &ESConf{}
	err := json.Unmarshal([]byte(data), cfg)
	if err != nil {
		attr.AttrAPI(36337053, 1)
		log.Errorf("HandleEsConfig json Unmarshal error, err:%v", err)
		return err
	}
	log.Infof("HandleEsConfig, cfg:%+v", cfg)
	esConfCfgValue.Store(cfg)
	return nil
}

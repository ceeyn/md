package config_rainbow // NOCA:golint/package(设计如此)

import (
	"encoding/json"
	"sync/atomic"

	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
)

// WarmMusicConf ...
type WarmMusicConf struct {
	SlowMusic       string `json:"slow_music"`
	FastMusic       string `json:"fast_music"`
	WarmPowerSwitch string `json:"warm_power_switch"`
	DelaySpan       int64  `json:"delay_span"`    // tdmq的间隔时长
	MaxRetryCnt     int64  `json:"max_retry_cnt"` // tdmq的最大重试次数
}

var (
	musicConfCfgValue atomic.Value
)

func init() {
	musicConfCfgValue.Store(&WarmMusicConf{
		SlowMusic:       "/wemeet_webinar/mp3/408235998464462745/408235998464528281_slow.mp3", //默认
		FastMusic:       "",                                                                   //默认
		WarmPowerSwitch: "close",                                                              //开关默认关
		DelaySpan:       60,                                                                   //tdmq的间隔时长60s
		MaxRetryCnt:     25,                                                                   //tdmq的最大重试次数25次
	})
}

//GetMusicConfConfig ..
func GetMusicConfConfig() *WarmMusicConf {
	attr.AttrAPI(35927485, 1)
	cfg := musicConfCfgValue.Load().(*WarmMusicConf)
	log.Infof("GetMusicConfConfig:%+v", cfg)
	return cfg
}

//HandleWarmMusicConfConfig ...
func HandleWarmMusicConfConfig(data string) error {
	cfg := &WarmMusicConf{}
	err := json.Unmarshal([]byte(data), cfg)
	if err != nil {
		attr.AttrAPI(35927486, 1)
		log.Errorf("HandleWarmMusicConfConfig json Unmarshal error, err:%v", err)
		return err
	}
	log.Infof("HandleWarmMusicConfConfig, cfg:%+v", cfg)
	musicConfCfgValue.Store(cfg)
	return nil
}

//GetWarmPowerSwitchConfConfig ..
func GetWarmPowerSwitchConfConfig() string {
	metrics.IncrCounter("GetWarmPowerSwitchConfConfig_Cnt", 1)
	cfg := musicConfCfgValue.Load().(*WarmMusicConf)
	log.Infof("GetWarmPowerSwitchConfConfig:%+v", cfg.WarmPowerSwitch)
	return cfg.WarmPowerSwitch
}

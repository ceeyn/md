package config_rainbow // NOCA:golint/package(设计如此)

import (
	"context"
	"git.code.oa.com/meettrpc/meet_util"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/config"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

var (
	configName  = "rainbow"
	participant = "participant"
	music       = "music"
	background  = "background"
	nameBadge   = "nameBadge"
	esConfig    = "esConfig"
	common      = "common"
	callBackKey = "callBackKey"
)

// InitRainV2 初始化
func InitRainV2() error {
	watchRainbow(configName, participant, HandleParticipantConfConfig)
	watchRainbow(configName, music, HandleWarmMusicConfConfig)
	watchRainbow(configName, background, HandleBackgroundConfConfig)
	watchRainbow(configName, esConfig, HandleEsConfig)
	watchRainbow(configName, common, HandleCommonConfig)
	watchRainbow(configName, callBackKey, HandleCallBackConf)
	return nil
}

//watchRainbow ..
func watchRainbow(name string, key string, handle func(string) error) {
	c, err := config.Get(name).Get(context.TODO(), key)
	if err != nil {
		log.Errorf("get config failed: %s", err.Error())
		panic(err)
	} else {
		log.Infof("config: %s", c.Value())
		err = handle(c.Value())
		if err != nil {
			log.Errorf("get meeting overload config failed, error :%s", err.Error())
			panic(err)
		}
	}
	go func() {
		defer meet_util.DefPanicFun()
		ch, _ := config.Get(name).Watch(trpc.BackgroundContext(), key)
		for r := range ch {
			log.Infof("watch wechat_config :%+v, MetaData:%+v", r.Value(), r.MetaData())
			err = handle(r.Value())
			if err != nil {
				log.Errorf("handle watch config event failed,key:%s error :%s", key, err.Error())
			}
		}
	}()
}

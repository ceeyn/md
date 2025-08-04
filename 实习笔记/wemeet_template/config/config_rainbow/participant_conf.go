package config_rainbow // NOCA:golint/package(设计如此)

import (
	"encoding/json"
	"sync/atomic"

	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

// ParticipantConf ...
type ParticipantConf struct {
	ParticipantMaxCount uint64 `json:"participant_max_count"`
	ScheduleMaxCount    uint64 `json:"schedule_max_count"`
}

var (
	participantConfCfgValue atomic.Value
)

func init() {
	participantConfCfgValue.Store(&ParticipantConf{
		ParticipantMaxCount: 200, //默认200
		ScheduleMaxCount:    50,  //默认50
	})
}

//GetParticipantConfConfig ..
func GetParticipantConfConfig() *ParticipantConf {
	attr.AttrAPI(35919087, 1)
	cfg := participantConfCfgValue.Load().(*ParticipantConf)
	log.Infof("GetParticipantConfConfig:%+v", cfg)
	return cfg
}

//HandleParticipantConfConfig ...
func HandleParticipantConfConfig(data string) error {
	cfg := &ParticipantConf{}
	err := json.Unmarshal([]byte(data), cfg)
	if err != nil {
		attr.AttrAPI(35919088, 1)
		log.Errorf("HandleParticipantConfConfig json Unmarshal error, err:%v", err)
		return err
	}
	log.Infof("HandleParticipantConfConfig, cfg:%+v", cfg)
	participantConfCfgValue.Store(cfg)
	return nil
}

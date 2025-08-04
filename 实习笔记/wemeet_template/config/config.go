package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	"git.code.oa.com/going/config"
	trpc "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

var (
	buildstamp = ""
	gitversion = ""
	gitbranch  = ""
	goversion  = ""
	sysuname   = ""
	ConfigPath = "../conf/" //配置文件路径
	TdmqConfig TDMQConfig
)

// TDMQConfig TDMQ配置
type TDMQConfig struct {
	TdmqClusterUrl string `yaml:"TdmqClusterUrl"`
	TdmqTopic      string `yaml:"TdmqTopic"`
	TdmqRoleToken  string `yaml:"TdmqRoleToken"`
}

var Conf = struct {
	//msgBox配置
	MsgBoxConfig struct {
		Secret string `json:"Secret"`
	}
}{}

//信息
func PrintVersion(c bool) {
	f := func(c bool, format string, a ...interface{}) {
		if c {
			log.Infof(format, a)
		} else {
			fmt.Printf(format, a)
		}
	}
	f(c, "\tGitVersion: %v \n", gitversion)
	f(c, "\tGitBranch: %v \n", gitbranch)
	f(c, "\tBuildStamp: %v \n", buildstamp)
	f(c, "\tGoVersion: %v \n", goversion)
	f(c, "\tSysUname: %v \n", sysuname)

}

//更改 增加 可执行程序版本信息
func InitFlag() {
	verOpt := flag.Bool("v", false, "Print application version")
	ConfigPathOpt := flag.String("c", "", "config file path")
	flag.Parse()
	if *verOpt {
		PrintVersion(false)
		os.Exit(0) //输出版本信息,退出
	}
	//更改配置文件
	if *ConfigPathOpt != "" {
		ConfigPath = *ConfigPathOpt
		trpc.ServerConfigPath = ConfigPath + "trpc_go.yaml"
		log.Infof("change ServerConfigPath : %v \n", trpc.ServerConfigPath)
	}
	PrintVersion(true)
}

// InitImConfig .
func InitMsgBoxConfig() {
	config.ConfPath = "./conf/config.toml"
	err := config.Parse(&Conf)
	if err != nil {
		log.Errorf(" get config.toml fail %v", err.Error())
		panic(err)
	}
}

// InitConfig 初始化配置
func InitConfig(fileName string) error {
	//加载TDMQ配置文件
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Error("read config file fail:", err)
		return err
	}
	err = yaml.Unmarshal(content, &TdmqConfig)
	if err != nil {
		log.Error("invalid yaml config:", err)
		return err
	}
	log.Infof("customer.yaml :%+v", TdmqConfig)
	return nil
}

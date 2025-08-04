package main

import (
	"meeting_template/config"
	"meeting_template/config/config_rainbow"
	"meeting_template/dao"
	"meeting_template/es"
	"meeting_template/kafka"
	"meeting_template/rpc"
	"meeting_template/service"
	"meeting_template/service_kafka"

	_ "git.code.oa.com/cm-metrics/trpc-metrics-tccm-go"
	_ "git.code.oa.com/meettrpc/meet_logic"
	_ "git.code.oa.com/meettrpc/meetlog"
	_ "git.code.oa.com/phoenixs/degraded"
	_ "git.code.oa.com/phoenixs/polaris_alert"
	_ "git.code.oa.com/phoenixs/trpc-log-meet/log" // 匿名引入插件配置
	_ "git.code.oa.com/trpc-go/trpc-codec/oidb"
	_ "git.code.oa.com/trpc-go/trpc-config-rainbow"
	_ "git.code.oa.com/trpc-go/trpc-config-tconf"
	_ "git.code.oa.com/trpc-go/trpc-filter/recovery"
	_ "git.code.oa.com/trpc-go/trpc-log-atta"
	_ "git.code.oa.com/trpc-go/trpc-metrics-attr"
	_ "git.code.oa.com/trpc-go/trpc-metrics-m007"
	_ "git.code.oa.com/trpc-go/trpc-metrics-runtime"
	_ "git.code.oa.com/trpc-go/trpc-naming-polaris"
	_ "git.code.oa.com/trpc-go/trpc-opentracing-tjg"
	_ "git.woa.com/opentelemetry/opentelemetry-go-ecosystem/instrumentation/oteltrpc"
	_ "git.woa.com/tencent-meeting/wemeet_event_access"
	_ "go.uber.org/automaxprocs"

	skm "git.code.oa.com/phoenixs/secret-key-sdk-go"
	td_kafka "git.code.oa.com/trpc-go/trpc-database/kafka"
	"git.code.oa.com/trpc-go/trpc-database/tdmq"
	trpc "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	auth "git.woa.com/phoenixs/interface_auth" //鉴权模块。其他包不列举，请自行引入
	mh "git.woa.com/phoenixs/meet-http"
)

func main() {
	mh.InitMeetHttpTransport()
	//注意这里的配置 路径, 配置文件在当前项目的上一级, 配合织云打包
	trpc.ServerConfigPath = "./conf/trpc_go.yaml"
	//初始化命令行参数,里面可以强行更改 trpc.ServerConfigPath的目录
	config.InitFlag()
	config.InitMsgBoxConfig()
	if err := config.InitConfig("./conf/customer.yaml"); err != nil {
		metrics.IncrCounter("InitConfigFailed", 1)
		log.Fatal(err)
	}
	s := skm.NewTrpcServer()    //切换成SDK提供启动方法
	config_rainbow.InitRainV2() //先初始化七彩石

	if err := auth.InitByRainbow("interface_auth", "interface_auth"); err != nil {
		log.Fatal(err)
		panic(err)
	}

	dao.InitDb()
	es.Init()
	kafka.InitProducer()
	rpc.InitTdmqProducer()
	tdmq.RegisterConsumerService(
		s.Service("trpc.wemeet.wemeet_template.tdmq_consumer"),
		&service.Consumer{},
	)
	pb.RegisterWemeetMeetingTemplateOidbService(s, &service.WemeetMeetingTemplateOidbServiceImpl{})
	pb.RegisterWemeetMeetingTemplateHttpService(s, &service.WemeetMeetingTemplateHttpServiceImpl{})

	td_kafka.RegisterConsumerService(s.Service("trpc.kafka.webinar_template.service"), &service_kafka.KafkaConsumer{})

	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}

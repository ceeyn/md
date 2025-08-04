module meeting_template

go 1.12

require (
	git.code.oa.com/cm-metrics/trpc-metrics-tccm-go v1.1.0
	git.code.oa.com/going/attr v0.3.1
	git.code.oa.com/going/config v0.1.0
	git.code.oa.com/meettrpc/meet_logic v1.0.19
	git.code.oa.com/meettrpc/meet_util v1.0.33
	git.code.oa.com/meettrpc/meetlog v1.0.7
	git.code.oa.com/phoenixs/degraded v0.0.3
	git.code.oa.com/phoenixs/polaris_alert v0.0.9
	git.code.oa.com/phoenixs/secret-key-sdk-go v0.1.26
	git.code.oa.com/phoenixs/trpc-log-meet v0.1.27
	git.code.oa.com/rainbow/golang-sdk v0.4.12 // indirect
	git.code.oa.com/trpc-go/trpc-codec/oidb v0.2.16
	git.code.oa.com/trpc-go/trpc-config-rainbow v0.1.18
	git.code.oa.com/trpc-go/trpc-config-tconf v0.1.8
	git.code.oa.com/trpc-go/trpc-database/gorm v0.2.2
	git.code.oa.com/trpc-go/trpc-database/kafka v0.2.9
	git.code.oa.com/trpc-go/trpc-database/redis v0.1.8
	git.code.oa.com/trpc-go/trpc-database/tdmq v0.2.10
	git.code.oa.com/trpc-go/trpc-filter/recovery v0.1.2
	git.code.oa.com/trpc-go/trpc-go v0.9.5
	git.code.oa.com/trpc-go/trpc-log-atta v0.1.13
	git.code.oa.com/trpc-go/trpc-metrics-attr v0.2.2
	git.code.oa.com/trpc-go/trpc-metrics-m007 v0.4.2
	git.code.oa.com/trpc-go/trpc-metrics-runtime v0.3.3
	git.code.oa.com/trpc-go/trpc-naming-polaris v0.3.5-patch
	git.code.oa.com/trpc-go/trpc-opentracing-tjg v0.1.8
	git.code.oa.com/trpc-go/trpc-selector-cl5 v0.2.0
	git.code.oa.com/trpcprotocol/tencent_meeting/common_service_role_auth v1.4.68
	git.code.oa.com/trpcprotocol/wemeet/common_meeting_cache v1.4.15
	git.code.oa.com/trpcprotocol/wemeet/common_msgbox v1.1.19
	git.code.oa.com/trpcprotocol/wemeet/common_upload v1.2.1
	git.code.oa.com/trpcprotocol/wemeet/common_xcast_im_conversion_logic v1.11.72
	git.code.oa.com/trpcprotocol/wemeet/common_xcast_meeting_error_code v1.2.51
	git.code.oa.com/trpcprotocol/wemeet/layout_center v1.3.20
	git.code.oa.com/trpcprotocol/wemeet/meeting_template v1.2.13
	git.code.oa.com/trpcprotocol/wemeet/wemee_db_proxy_center_join_wemeet_db_proxy_center_join v1.1.12
	git.code.oa.com/trpcprotocol/wemeet/wemeet_kafka_message v1.4.5
	git.code.oa.com/trpcprotocol/wemeet/wemeet_meet_sensitive v1.1.61
	git.code.oa.com/trpcprotocol/wemeet/wemeet_safe_gate v1.1.52
	git.code.oa.com/trpcprotocol/wemeet/wemeet_safe_gateway v1.2.48
	git.code.oa.com/trpcprotocol/wemeet/wemeet_user_notify v1.1.41
	git.code.oa.com/wesee_ugc/go.uuid v1.2.0
	git.woa.com/tencent-meeting/wemeet_event_access v1.0.2
	git.woa.com/opentelemetry/opentelemetry-go-ecosystem/instrumentation/oteltrpc v0.3.1
	git.woa.com/phoenixs/interface_auth v0.0.7
	git.woa.com/phoenixs/meet-http v1.0.4
	git.woa.com/wemeet-public/sdk-meet/safe/evt_ugc_config_join_welcome v1.0.1-beta
	github.com/agiledragon/gomonkey v2.0.2+incompatible
	github.com/agiledragon/gomonkey/v2 v2.2.0
	github.com/apache/pulsar-client-go v0.10.0
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.6.0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/olivere/elastic v6.2.37+incompatible
	github.com/pkg/errors v0.9.1
	github.com/sony/sonyflake v1.0.0
	github.com/stretchr/testify v1.8.2
	github.com/tencentyun/cos-go-sdk-v5 v0.7.50
	go.uber.org/automaxprocs v1.4.0
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	google.golang.org/protobuf v1.30.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/gorm v1.23.6
)

replace (
	git.code.oa.com/cm-metrics/trpc-metrics-tccm-go v1.2.2 => git.code.oa.com/phoenixs/wemeet-cm v1.0.3
	git.code.oa.com/going/going v0.4.5 => git.code.oa.com/meettrpc/going v1.0.4
	//git.code.oa.com/trpc-go/trpc-go => git.code.oa.com/trpc-go/trpc-go v0.7.3
	git.code.oa.com/trpcprotocol/wemeet/common_xcast_im_comm v1.1.15 => git.woa.com/trpcprotocol/wemeet/common_xcast_im_comm v1.1.15
)

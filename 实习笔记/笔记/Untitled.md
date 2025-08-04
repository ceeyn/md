问题一：为什么没有打印到同一个地方

```
    log.InfoContextf(c.Ctx, "EncryptData-EncryptField [uid]:%+v,enField:%+v,GCbcKey:%+v",
       tmpField, string(enField), enKey)
} else if field == "meeting_code" {
    log.InfoContextf(c.Ctx, "EncryptData-EncryptField [meeting_code]:%+v,enField:%+v,GCbcKey:%+v",
       tmpField, string(enField), enKey)
```

问题二：period cache check and get failed IsSuccHandleService() failed!





### 一。客户端私参加密前

### 二。客户端私参加密后

### 三、非腾讯公参加密前

这段日志记录了在调用 `EncryptBackendPublicParaData` 方法之前的系统状态和参数信息。以下是对这段日志的详细解释：

### 日志解析

#### 基本信息

1. **Topic**: `wemeet-backend_report_kafka-cicd-ap-guangzhou`
   - Kafka 主题名称，表示日志消息所属的主题，用于消息分类和订阅。
2. **Biz_seq**: `1aab6f4f-828c-4d10-8526-6e1cf528d0af`
   - 业务序列号，用于唯一标识某个业务流程或交易。
3. **TraceID**: `70436ba1e7e04b751863d5a18d510bd9`
   - 跟踪 ID，用于关联分布式系统中同一请求的多个日志条目，方便跟踪和调试请求流程。
4. **Message**: `Before EncryptBackendPublicParaData (isTxCrop:false):map[...]`
   - 日志消息的具体内容，表示在调用 `EncryptBackendPublicParaData` 方法之前的参数状态，其中 `isTxCrop:false` 表示当前未启用 Tx 企业加密。

#### 服务器相关信息

1. **Server**: `log_svr`
   - 记录日志的服务器名称。
2. **Offset**: `1918047398`
   - Kafka 消息的偏移量，用于标识消息在 Kafka 分区中的位置。
3. **Env**: `test`
   - 表示当前的运行环境是测试环境。

#### 时间戳和代码位置信息

1. **Timestamp**: `2024-07-24 10:36:34.054`
   - 日志记录的时间戳，表示日志产生的具体时间。
2. **Line**: `mtx/cmd_atta_comm.go:239`
   - 代码文件及其行号，表示日志是在 `mtx/cmd_atta_comm.go` 文件的第 239 行产生的。

#### 追踪和分区信息

1. **SpanID**: `b14885fd73de1078`
   - 一个跨度 ID，通常用于分布式追踪，表示一个特定的操作或事务。
2. **Partition**: `1`
   - Kafka 主题的分区编号，表示消息存储在主题的哪个分区中。

#### 日志级别和采样

1. **TraceId**: `523883659112002920`
   - 另一个跟踪 ID，用于追踪请求的流程。
2. **Topic**: `center_report_test`
   - 另一个 Kafka 主题的名称，表示日志消息可能也被发送到这个主题。
3. **Level**: `INFO`
   - 日志级别，表示这是一个信息日志。
4. **Sampled**: `false`
   - 表示这个日志条目是否被采样，这里表示未被采样。

#### 时间和文件路径

1. **Time**: `2024-07-24 10:36:34.066`
   - 另一个时间戳，可能表示日志被处理或传输的时间。
2. **Filename**: `/logdata/backend_report_kafka/trpc.log`
   - 日志文件的路径，表示日志被存储在服务器的哪个文件中。

#### 日志来源信息

1. **Source**: `30.188.181.104`
   - 生成日志的服务器的 IP 地址。
2. **Hostname**: `backend-report-kafka-test-0`
   - 生成日志的服务器的主机名。
3. **Indexname**: `backend_report_kafka`
   - 日志索引名称，通常用于日志管理系统（如 ElasticSearch）中。
4. **CICD**: `056b5e23`
   - 可能是一个与 CI/CD（持续集成/持续交付）相关的标识符，用于标识某个构建或部署过程。

### 具体消息内容解析

**日志消息**:

```
plaintext
复制代码
Before EncryptBackendPublicParaData (isTxCrop:false):map[city: client_ip: client_ip_port:0 corp_id:200000001 country: current_account_type: device_id: error_code:0 error_msg:LCBwZXJpb2QgY2FjaGUgY2hlY2sgYW5kIGdldCBmYWlsZWQgSXNTdWNjSGFuZGxlU2VydmljZSgpIGZhaWxlZCE= event_code:e#instruction_error_code event_time:1721788593 event_value:operate_role%3D3%26str_to_app_uid%3D%26uint32_to_sdkappid%3D0%26uint32_from_sdkappid%3D200000001%26creator_tiny_id%3D0%26uint32_from_user_access_type%3DUSER_ACCESS_VOIP%26creator_corp_id%3D0%26nick_name%3DEWM2GCNo1K6mXP7ZDJJjAg%3D%3D%26str_from_app_uid%3Ddefault%26call_source_type%3D0%26creator_uid%3D%26cmd_error_code%3D19801%26creator_instance_id%3D0%26str_from_phone_number%3D%26uint32_from_device_type%3D10000%26uint32_to_device_type%3D10000%26uint32_to_instance_id%3D0%26enter_env_type%3D1%26uint32_from_instance_id%3D0%26uint32_service_main_cmd%3D1445%26str_to_phone_number%3D instance_id:0 is_free_meeting:1 is_multclient:0 is_overseas:0 language: meeting_code: meeting_id:12087449934123025173 meeting_type:0 operator: pro_level:0 product:meeting province: qimei_36: result_code:1 result_msg:8ac3c2db-4965-11ef-acfa-525400ca2f23 sdk_id:0 server_ip:30.188.181.104 server_time:1721788594 source_id:3 sub_meeting_id: tiny_id:0 uid:default unique_report_id:a0686174-918e-4f04-8b1c-a51f35136bf4 version:]
```

这段信息是在调用 `EncryptBackendPublicParaData` 方法之前记录的，展示了当前状态下的请求参数。以下是对关键字段的解释：

- isTxCrop

  ：表示当前未启用 Tx 企业加密。

- **city**：城市名称。

- **client_ip**：客户端 IP 地址。

- **client_ip_port**：客户端 IP 端口号，当前值为 0。

- **corp_id**：企业 ID，当前值为 `200000001`。

- **country**：国家名称。

- **current_account_type**：当前账户类型。

- **device_id**：设备 ID。

- **error_code**：错误代码，当前值为 `0`。

- **error_msg**：错误消息，当前值为 `LCBwZXJpb2QgY2FjaGUgY2hlY2sgYW5kIGdldCBmYWlsZWQgSXNTdWNjSGFuZGxlU2VydmljZSgpIGZhaWxlZCE=`（这是一个 Base64 编码的字符串，解码后为 "period cache check and get failed IsSuccHandleService() failed!"）。

- **event_code**：事件代码，当前值为 `e#instruction_error_code`。

- **event_time**：事件时间戳，当前值为 `1721788593`。

- event_value

  ：事件值，是一个编码后的查询字符串，包含以下信息：

  - `operate_role=3`
  - `str_to_app_uid=`
  - `uint32_to_sdkappid=0`
  - `uint32_from_sdkappid=200000001`
  - `creator_tiny_id=0`
  - `uint32_from_user_access_type=USER_ACCESS_VOIP`
  - `creator_corp_id=0`
  - `nick_name=EWM2GCNo1K6mXP7ZDJJjAg==`（Base64 编码的昵称）
  - `str_from_app_uid=default`
  - `call_source_type=0`
  - `creator_uid=`
  - `cmd_error_code=19801`
  - `creator_instance_id=0`
  - `str_from_phone_number=`
  - `uint32_from_device_type=10000`
  - `uint32_to_device_type=10000`
  - `uint32_to_instance_id=0`
  - `enter_env_type=1`
  - `uint32_from_instance_id=0`
  - `uint32_service_main_cmd=1445`
  - `str_to_phone_number=`

- **instance_id**：实例 ID，当前值为 `0`。

- **is_free_meeting**：是否是免费会议，当前值为 `1`。

- **is_multclient**：是否是多客户端，当前值为 `0`。

- **is_overseas**：是否是海外用户，当前值为 `0`。

- **language**：语言设置。

- **meeting_code**：会议代码。

- **meeting_id**：会议 ID，当前值为 `12087449934123025173`。

- **meeting_type**：会议类型，当前值为 `0`。

- **operator**：运营商。

- **pro_level**：专业级别，当前值为 `0`。

- **product**：产品，当前值为 `meeting`。

- **province**：省份。

- **qimei_36**：设备唯一标识符。

- **result_code**：结果代码，当前值为 `1`。

- **result_msg**：结果消息，当前值为 `8ac3c2db-4965-11ef-acfa-525400ca2f23`。

- **sdk_id**：SDK ID，当前值为 `0`。

- **server_ip**：服务器 IP 地址，当前值为 `30.188.181.104`。

- **server_time**：服务器时间戳，当前值为 `1721788594`。

- **source_id**：来源 ID，当前值为 `3`。

- **sub_meeting_id**：子会议 ID。

- **tiny_id**：Tiny ID，当前值为 `0`。

- **uid**：用户 ID，当前值为 `default`。

- **unique_report_id**：唯一报告 ID，当前值为 `a0686174-918e-4f04-8b1c-a51f35136bf4`。

- **version**：版本信息。

### 总结

这段日志提供了在调用 `EncryptBackendPublicParaData` 方法之前的详细快照，包括大量的参数和上下文信息。这些信息对于调试和分析系统行为非常有用。通过这些日志，可以追踪特定请求的处理流程，了解在特定时间点系统的状态和所处理的数据

### 四、非腾讯公参加密后

```
// 记录加密后的公参数据日志
	log.InfoContextf(c.Ctx, "After EncryptBackendPublicParaData (isTxCrop:%v):%+v", c.IsEncrypt, reqMap)
```

这段日志消息记录了在调用 `EncryptBackendPublicParaData` 方法后的系统状态，包括各种参数和上下文信息。以下是对这段日志消息的详细解释：

### 日志解析

#### 基本信息

1. **Topic**: `wemeet-backend_report_kafka_other-cicd-ap-guangzhou-2552`
   - Kafka 主题名称，表示日志消息所属的主题，用于消息分类和订阅。
2. **Biz_seq**: `1e12caac-3cbe-46c9-b6e1-e104acaed1a9`
   - 业务序列号，用于唯一标识某个业务流程或交易。
3. **TraceID**: `ee883fb758eb2108cf10f804d2852034`
   - 跟踪 ID，用于关联分布式系统中同一请求的多个日志条目，方便跟踪和调试请求流程。
4. **Message**: `After EncryptBackendPublicParaData (isTxCrop:false):map[...]`
   - 日志消息的具体内容，表示在调用 `EncryptBackendPublicParaData` 方法后的参数状态，其中 `isTxCrop:false` 表示当前未启用 Tx 企业加密。

#### 服务器相关信息

1. **Server**: `log_svr`
   - 记录日志的服务器名称。
2. **Offset**: `5067633872`
   - Kafka 消息的偏移量，用于标识消息在 Kafka 分区中的位置。
3. **Env**: `test`
   - 表示当前的运行环境是测试环境。

#### 时间戳和代码位置信息

1. **Timestamp**: `2024-07-24 10:09:27.080`
   - 日志记录的时间戳，表示日志产生的具体时间。
2. **Line**: `mtx/cmd_atta_comm.go:245`
   - 代码文件及其行号，表示日志是在 `mtx/cmd_atta_comm.go` 文件的第 245 行产生的。

#### 追踪和分区信息

1. **SpanID**: `8870e7bfcc34b566`
   - 一个跨度 ID，通常用于分布式追踪，表示一个特定的操作或事务。
2. **Partition**: `0`
   - Kafka 主题的分区编号，表示消息存储在主题的哪个分区中。

#### 日志级别和采样

1. **TraceId**: `523880929492791579`
   - 另一个跟踪 ID，用于追踪请求的流程。
2. **Topic**: `other_report_test`
   - 另一个 Kafka 主题的名称，表示日志消息可能也被发送到这个主题。
3. **Level**: `INFO`
   - 日志级别，表示这是一个信息日志。
4. **Sampled**: `false`
   - 表示这个日志条目是否被采样，这里表示未被采样。

#### 时间和文件路径

1. **Time**: `2024-07-24 10:09:27.110`
   - 另一个时间戳，可能表示日志被处理或传输的时间。
2. **Filename**: `/logdata/backend_report_kafka_other/trpc.log`
   - 日志文件的路径，表示日志被存储在服务器的哪个文件中。

#### 日志来源信息

1. **Source**: `11.151.241.27`
   - 生成日志的服务器的 IP 地址。
2. **Hostname**: `backend-report-kafka-other-test-1`
   - 生成日志的服务器的主机名。
3. **Indexname**: `wemeet_backend_report_kafka_other_cicd_ap_guangzhou_2552`
   - 日志索引名称，通常用于日志管理系统（如 ElasticSearch）中。

### 具体消息内容解析

**日志消息**:

```
plaintext
复制代码
After EncryptBackendPublicParaData (isTxCrop:false):map[event_code:e#tab_exp_expose event_value:user_type%3D%26exp_method%3Dpassive%26upload_time%3D2024-07-24+09%3A39%3A29%26scene_id%3D-1%26sdk_version%3DGO_SDK_1.5.9%26business_code%3D1251%26env%3Dpro%26exp_percentage%3D%26exp_params%3D experiment_group_id:11933137 product:meeting server_ip:11.151.241.27 server_time:1721786967 source_id:2 uid:16801441153947357592726 unique_report_id:9a1c94d0-9c61-4971-a270-8ae0a1fad0b4 upload_time:2024-07-24 09:39:29]
```

这段信息是在调用 `EncryptBackendPublicParaData` 方法之后记录的，展示了当前状态下的请求参数。以下是对关键字段的解释：

- isTxCrop

  ：表示当前未启用 Tx 企业加密。

- **event_code**：事件代码，当前值为 `e#tab_exp_expose`。

- event_value

  ：事件值，是一个编码后的查询字符串，包含以下信息：

  - `user_type=`：用户类型（未填值）。
  - `exp_method=passive`：曝光方法，被动。
  - `upload_time=2024-07-24 09:39:29`：上传时间。
  - `scene_id=-1`：场景 ID。
  - `sdk_version=GO_SDK_1.5.9`：SDK 版本。
  - `business_code=1251`：业务代码。
  - `env=pro`：环境，生产环境。
  - `exp_percentage=`：曝光百分比（未填值）。
  - `exp_params=`：曝光参数（未填值）。

- **experiment_group_id**：实验组 ID，当前值为 `11933137`。

- **product**：产品，当前值为 `meeting`。

- **server_ip**：服务器 IP 地址，当前值为 `11.151.241.27`。

- **server_time**：服务器时间戳，当前值为 `1721786967`。

- **source_id**：来源 ID，当前值为 `2`。

- **uid**：用户 ID，当前值为 `16801441153947357592726`。

- **unique_report_id**：唯一报告 ID，当前值为 `9a1c94d0-9c61-4971-a270-8ae0a1fad0b4`。

- **upload_time**：上传时间，当前值为 `2024-07-24 09:39:29`。



### 五、非腾讯私参加密前后





### 消息结构

#### `GetHealthReq` 和 `GetHealthRsp`

- **GetHealthReq**: 空的消息，用于健康检查请求。
- **GetHealthRsp**: 包含一个可选的字符串字段 `status`，用于返回健康检查的状态。

#### `EventParams`

- EventParams

  : 用于表示事件的参数。

  - `optional string event_code`: 事件代码。
  - `map<string, string> event_value`: 事件的键值对。
  - `optional uint64 event_time`: 事件发生的时间（毫秒值）。

#### `MessCommParams`

- MessCommParams

  : 用于公参结构体。

  - `map<string, string> comm_params`: 公参的键值对。

#### 枚举类型

##### `REPORT_FLAG`

- REPORT_FLAG

  : 定义不同的上报标志。

  - `REPORT_FLAG_DEFAULT`: 默认不带公参的上报。
  - `REPORT_FLAG_INIT`: 需要重新全量公参的情况。
  - `REPORT_FLAG_UPDATE`: 更新公参。
  - `REPORT_FLAG_TRANS`: 透传，不需要读取公参。

##### `REQ_TYPE`

- REQ_TYPE

  : 用于区分请求结构体的类型。

  - `REQ_TYPE_DEFAULT`: 默认类型。
  - `REQ_TYPE_CLIENT`: 客户端上报。
  - `REQ_TYPE_SECOND`: 秒控上报。
  - `REQ_TYPE_BACKEND_MAP`: 服务端 map 格式上报。
  - `REQ_TYPE_BACKEND_STRING`: 服务端 string 格式上报。
  - `REQ_TYPE_BACKEND_STANDARD`: 后端上报的标准格式。

##### `CONN_TYPE`

- CONN_TYPE

  : 定义连接类型。

  - `DEFAULT`: 默认类型。
  - `TDBANK`: TDBANK 类型。
  - `ATTA`: ATTA 类型。
  - `KAFKA`: KAFKA 类型。
  - `TDBANK_COMM_MONITOR`: 发送到监控类的 TDBANK 表。
  - `TDBANK_PUBLIC`: 根据加载的配置表生成公参。

##### `RETURN_CODE`

- RETURN_CODE

  : 定义返回代码。

  - `SUCCESS`: 成功。
  - `ERR_SYSTEM`: 系统错误。
  - `ERR_PARAMETER`: 错误的上报参数。
  - `ERR_NO_MATCH`: 没有匹配到规则。
  - `ERR_SESSION_ID`: 错误的 session ID。
  - `ERR_STORE_SESSION`: session ID 没保存。
  - `ERR_UPDATE_COMM`: 更新公参失败。
  - `ERR_SEESION_SEQ`: 错误序列号。

##### `CLIENT_ACTION`

- CLIENT_ACTION

  : 定义客户端动作。

  - `ACTION_DEFAULT`: 默认无意义。
  - `UPDATE_COMM`: 强制更新公参。
  - `MUST_TRANS`: 强制透传消息。

##### `PRODUCT_TYPE`

- PRODUCT_TYPE

  : 定义产品类型。

  - `PRODUCT_DEFAULT`: 无意义。
  - `RRODUCT_MEETING`: 国内会议。
  - `RRODUCT_VOOV`: 海外。
  - `RRODUCT_CALENDAR`: 日历。

##### `MEETING_TYPE`

- MEETING_TYPE

  : 定义会议类型。

  - `MEETING_TYPE_DEFAULT`: 无意义。
  - `MEETING_TYPE_NORMAL`: 普通会议。
  - `MEETING_TYPE_CYCLE`: 周期性会议。
  - `MEETING_TYPE_PMI`: PMI 会议。

##### `SCENE_TYPE`

- SCENE_TYPE

  : 定义会议场景类型。

  - `SCENE_TYPE_DEFAULT`: 无意义。
  - `SCENE_TYPE_NORMAL`: 普通场景会议。
  - `SCENE_TYPE_WEBNIAR`: Webinar 会议。
  - `SCENE_TYPE_GROUP`: 分组会议。

### 消息结构（续）

#### `BasicMessageReq`

- BasicMessageReq

  : 客户端的连接请求。

  - `optional string report_cmd`: 上报类型。
  - `optional uint64 request_time`: 发送时间（毫秒值）。
  - `optional string client_ip`: 客户端 IP。
  - `map<string, string> comm_params`: 公参列表。
  - `repeated EventParams event_params`: 事件参数。
  - `optional int32 report_flag`: 上报标志。
  - `optional int32 debug`: debug 字段。
  - `optional int32 direct`: 是否直达消息。
  - `optional string source_id`: 消息来源。

#### `BasicMessageRsp`

- BasicMessageRsp

  : 基本消息的响应。

  - `optional int32 ret_code`: 返回代码。
  - `optional string ret_msg`: 返回消息。
  - `optional int32 action`: 动作。

#### `ServerReq`

- ServerReq

  : 秒级监控上报协议。

  - 包含多个字段，例如 `tdbank_bid`、`tdbank_tid`、`sys_id`、`intf_id` 等，用于描述秒级监控的详细信息。

#### `ServerRsp`

- ServerRsp

  : 服务器响应。

  - `optional int32 ret_code`: 返回代码。
  - `optional string ret_msg`: 返回消息。

#### `BackendReportMapReq`

- BackendReportMapReq

  : 通过 map 的形式请求到后端。

  - `optional uint64 request_time`: 发送时间（毫秒值）。
  - `optional string client_ip`: 客户端 IP。
  - `map<string, string> comm_params`: 公参列表。
  - `repeated EventParams event_params`: 事件参数。

#### `BackendReportMapRsp`

- BackendReportMapRsp

  : 后端 map 请求的响应。

  - `optional int32 ret_code`: 返回代码。
  - `optional string ret_msg`: 返回消息。

#### `BackendReportReq`

- BackendReportReq

  : 后台标准形式的上报。

  - `optional uint64 request_time`: 发送时间（毫秒值）。
  - `optional string event_code`: 事件代码。
  - `map<string, string> comm_params`: 公参列表。
  - `map<string, string> event_value`: 业务参数。
  - `optional string client_ip`: 上报的机器 IP。
  - `map<string, string> client_params`: 客户端透传。
  - `map<string, string> server_params`: 服务端透传。
  - `optional string version`: 版本号。

#### `BackendReportRsp`

- BackendReportRsp

  : 后端标准形式上报的响应。

  - `optional int32 ret_code`: 返回代码。
  - `optional string ret_msg`: 返回消息。

#### `BackendReportStringReq`

- BackendReportStringReq

  : 通过 string 的形式请求到后端。

  - `optional uint64 request_time`: 发送时间（毫秒值）。
  - `optional string client_ip`: 客户端 IP。
  - `optional string event_code`: 事件代码。
  - `optional string event_value`: 事件值。

#### `BackendReportStringRsp`

- BackendReportStringRsp

  : 后端 string 请求的响应。

  - `optional int32 ret_code`: 返回代码。
  - `optional string ret_msg`: 返回消息。

#### `Events`

- Events

  : 灯塔协议的私参。

  - `optional string eventCode`: 事件代码。
  - `optional string eventTime`: 事件发生时间。
  - `map<string, string> mapValue`: 键值对。

#### `BeaconReq`

- BeaconReq

  : 用于兼容灯塔协议上报的能力。

  - `optional string appVersion`: 应用版本。
  - `optional string sdkId`: SDK ID。
  - `optional string sdkVersion`: SDK 版本。
  - `optional string mainAppKey`: 主应用密钥。
  - `map<string, string> common`: 公参。
  - `repeated Events events`: 业务私参。

#### `BeaconRsp`

- BeaconRsp

  : 兼容灯塔协议上报的能力返回码。

  - `optional int32 result`: 返回代码。
  - `optional string srcGatewayIp`: 来源网关 IP。
  - `optional string serverTime`: 服务器时间。
  - `optional string msg`: 消息。

### 服务定义

#### `log_svr_http`

- log_svr_http

  : 定义了一些 HTTP 服务接口。

  - `rpc GetHealth(GetHealthReq) returns (GetHealthRsp)`: 健康检查。
  - `rpc DataReport(BasicMessageReq) returns (BasicMessageRsp)`: 客户端上报。
  - `rpc ReportMonitor(ServerReq) returns (ServerRsp)`: 秒控上报。
  - `rpc DataReportBeacon(BeaconReq) returns (BeaconRsp)`

#### `log_svr_trpc`

- log_svr_trpc

  : 定义了一些 tRPC 服务接口。

  - `rpc DataReport(BasicMessageReq) returns (BasicMessageRsp)`: 客户端上报。
  - `rpc ReportMonitor(ServerReq) returns (ServerRsp)`: 秒控上报。
  - `rpc ServerReportMap(BackendReportMapReq) returns (BackendReportMapRsp)`: 服务端上报（map 格式）。
  - `rpc ServerReportString(BackendReportStringReq) returns (BackendReportStringRsp)`: 服务端上报（string 格式）。
  - `rpc ServerReport(BackendReportReq) returns (BackendReportRsp)`: 服务端上报（标准格式）。

#### `log_svr_oidb`

- log_svr_oidb

  : 定义了一些 oidb 服务接口。

  - `rpc ReportMonitor(ServerReq) returns (ServerRsp)`: 秒控上报。

#### `backend_report_oidb`

- backend_report_oidb

  : 定义了一些 oidb 后端上报接口。

  - `rpc ServerReportMap(BackendReportMapReq) returns (BackendReportMapRsp)`: 服务端上报（map 格式）。
  - `rpc ServerReportString(BackendReportStringReq) returns (BackendReportStringRsp)`: 服务端上报（string 格式）。
  - `rpc ServerReport(BackendReportReq) returns (BackendReportRsp)`: 服务端上报（标准格式）。

#### `backend_report_http`

- backend_report_http

  : 定义了一些 HTTP 后端上报接口。

  - `rpc ServerReportMap(BackendReportMapReq) returns (BackendReportMapRsp)`: 服务端上报（map 格式）。
  - `rpc ServerReportString(BackendReportStringReq) returns (BackendReportStringRsp)`: 服务端上报（string 格式）。
  - `rpc ServerReport(BackendReportReq) returns (BackendReportRsp)`: 服务端上报（标准格式）。

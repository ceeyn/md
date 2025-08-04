## wiki
### 需求单: http://tapd.oa.com/Tencent_Meeting2/prong/stories/view/1020375092862437167
### 方案iwiki地址: https://iwiki.woa.com/pages/viewpage.action?pageId=680766484
### 协议rick地址: http://trpc.rick.oa.com/rick/pb/detail?id=21337

## 监控（ monitor , uls , 007, 天机阁 ） 
### monitor 监控配置

      属性ID及名称
      ID: 35620175    名称：template_GetHealthReq
      ID: 35620176    名称：template_CreateTemplateReq
      ID: 35620177    名称：template_CreateTemplateFail
      ID: 35620178    名称：template_GetTemplateInfoReq
      ID: 35620179    名称：template_GetTemplateInfoFail
      ID: 35620180    名称：template_UpdateTemplateReq
      ID: 35620181    名称：template_UpdateTemplateFail
      ID: 35620182    名称：template_RedisGetTemplateInfoReq
      ID: 35620183    名称：template_RedisGetTemplateInfoFail
      ID: 35620184    名称：template_RedisCreateSetTemplateInfo
      ID: 35620185    名称：template_RedisUpdateSetTemplateInfo
      ID: 35620186    名称：template_RedisSetTemplateInfoFail
      ID: 35620187    名称：template_RedisGetNeedCheckIdsFail
      ID: 35620188    名称：template_RedisSetNeedCheckIdFail
      ID: 35620189    名称：template_RedisSetExpire
      ID: 35620190    名称：template_RedisSetExpireFail
      ID: 35620191    名称：template_BatchQueryTempUrl
      ID: 35620192    名称：template_BatchQueryTempUrlFail
      ID: 35620193    名称：template_GetTextModerationReq
      ID: 35620194    名称：template_GetTextModerationFail
      ID: 35620195    名称：template_QueryMeetingSensitiveWords
      ID: 35620196    名称：template_QueryMeetingSensitiveWordsFail
      ID: 35620197    名称：template_TransCosIdToUrlFail
      ID: 35620198    名称：template_Base64DecodeDescriptionFail
      ID: 35620199    名称：template_Base64DecodeCoverUrlFail
      ID: 35620200    名称：template_HtmlParseDescriptionFail
      ID: 35620201    名称：template_ReplaceDiscriptionImgFail
      ID: 35620202    名称：template_ReplaceDiscriptionSensitiveFail
      ID: 35620203    名称：template_ReplaceDiscriptionTextFail
      ID: 35620204    名称：template_SponsorCoverNameHitSensitive

## 发布：编译地址(Qci), 现网发布地址(织云,stke ), 测试环境发布地址( 织云，stke ) 
之间一直使用STKE平台发布，后续应用管理部署编排完善后，需要迁到应用管理发布

## 服务接口介绍
### OIDB接口
#### 创建会议模板信息
    CreateTemplate(CreateTemplateReq) returns (CreateTemplateRsp); //@alias=/0x199000/1
处理逻辑：根据请求传过来的信息生成`template_id`作为key，并将信息存储到redis。

这里最早只用于webinar会议，在预订webinar会议时，先创建template_id,
创建成功后，center侧再预订会议并将template_id填充到meeting_info。创建template_id时，一般带过来的信息比较少，所以创建并存储redis的耗时
比较少。 后续普通云会议也增加了模版信息，在开启报名时会调用该接口。

暖场视频信息是后续新增的，所有关于视频审核状态的处理只放在Update接口。

#### 更新会议模板信息
    rpc UpdateTemplate(UpdateTemplateReq) returns (UpdateTemplateRsp); //@alias=/0x199000/2
处理逻辑： 更新模版信息。
* 请求里数据信息转换为redis存储结构
* 一般信息合法性检测处理
  * base64编解码处理
  * cosId格式检测
  * html文本解析，字段提取处理
  * json数据格式检测
  * 文本信息敏感词检测【同步】
  * 长文本解压缩处理
* 暖场视频相关信息检测处理
  * 暖场图片视频cosId格式检测
  * 暖场视频合法检测【异步】
  * 确认信息变更后，调用user_notify发送暖场物料变更通知
* 更新redis存储信息

这里最初主要处理的时详情页的富文本， 后续新增了暖场信息的处理，二者目前放在同一个key，处理代码相对独立，分布在不同的文件

#### 获取会议模板信息
    rpc GetTemplateInfo(GetTemplateInfoReq) returns (GetTemplateInfoRsp); //@alias=/0x199000/3
处理逻辑：根据请求里的`template_id`从redis读取存储的信息，并做相应转换组装回包。
* 读redis
* 长文本解压缩
* 提取cosIds
* 并行调用common_upload获取下载链接
* 组装数据回包

#### 获取暖场物料信息
    rpc GetWarmUpData(GetWarmUpDataReq) returns (GetWarmUpDataRsp);  //@alias=/0x199000/4
处理逻辑：跟上面GetTemplateInfo接口类似，只是处理的数据更少。
#### 获取品牌完整信息
    rpc GetWholeBrandInfo(GetWholeBrandInfoReq) returns (GetWholeBrandInfoRsp); //@alias=/0x199000/5
处理逻辑：根据请求里的`meeting_id`从redis读取meeting_info, 并通过meeting_info里的`template_id`信息获取templateInfo,之后组装数据回包。
* 读redis获取meeting_info
* 获取TemplateInfo(逻辑同GetTemplateInfo接口)
* 组装数据回包

#### 视频审核结果回调，用于获取视频审核结果
    rpc VideoCensorCallback (wemeet.wemeet_safe_gateway.GetVideoCallbackReq)   returns (wemeet.wemeet_safe_gateway.GetVideoCallbackRsp);  //@alias=/0x199000/6
处理逻辑：异步处理暖场视频审核结果，并将结果更新到redis存储。

一个暖场视频在调用updateTemplate首次上传会将状态设未"未审核",并调用送审接口，后续查询时发现状态为"未审核"时也会调用送审接口，
只在收到回调结果后才会修改审核状态。

### HTTP接口
#### 健康检测接口
    rpc GetHealth(GetHealthReq) returns (GetHealthRsp); //@alias=/health-check
该接口与业务无关。
#### 获取暖场物料信息
    rpc GetWarmUpData(GetWarmUpDataReq) returns (GetWarmUpDataRsp);  //@alias=/wemeet-template/get-warmup-data
处理逻辑：与同名OIDB接口完全相同，提供该接口方便web侧调用,避免cgi中间层做不必要的开发。


## 服务调用上下游
### 上游模块
* wemeet_center[包括拆分后的模块]
* cgi
* web
### 下游模块
* redis/cache_read 读取meeting_info, 读写templateInfo等信息
* common_upload 通过cos_id获取图片/视频下载链接
* wemeet_meet_sensitive  文本过信安检测
* wemeet_user_notify   暖场物料变更时，调用该模块推送消息
* wemeet_msgbox   暖场视频审核失败时，调用该模块推送消息
* wemeet_safe_gateway 视频过信安检测

## Q&A

## 资源: 有申请特殊资源的，需要备注资源的申请地址,以及接口人 

## 服务文件夹规则  
    目录层级关系如下：
    conf: 配置文件存放文件 
        目的：存放服务配置文件 
        规则: trpc_go.yaml的配置规则，请参考配置文件 
        备注：pri：私有化，pre预发布，prod：现网，test:测试，oversea_prod:海外现网  
    service|service_oidb|service_http: 业务处理目录 
        备注: 复杂的业务可以是 service/service_功能1， service/service_功能2， 
        目的：存放各个业务的处理文件 
        来源：对应trpc_go.yaml文件中的 server的目录配置 
        规则: 各个业务之间如果有通用的业务函数，可以放在 service里面 , 也可以根据功能来划分  
        命名：cmd_协议名_业务函数，如：cmd_oidb_get_user_info.go 
        备注：定时器必须统一使用 service_timer, 对应配置trpc_go.yaml文件 service.timer 
    rpc : 所有Rpc请求目录 
        目的：存放各个Rpc请求处理文件 
        来源：对应trpc_go.yaml的 client目录配置 
        规则: 包括 grocery，mysql, oidb调用，http调用，七彩石 等等 网络rpc调用 
        命名：rpc_业务服务，如 rpc_grocery.go 
        备注：单元测试中，需要对所有rpc目录文件，ApplyFunc该文件进行单元测试 
    util: 插件或者第三方库 
        目的：存放各种通用的工具函数或者第三方库
        来源：通用的代码  
        备注：util中不涉及第三方调用，一般要求有单元测试用例 
    config: 配置读取,
        目的：存储配置读取，
        规则: 包括 yaml文件，json文件,七彩石读取等 配置文件读取    
        命名：config_xxx  
        备注：common.go 增加了 全局变量和常量的定义  
    其他：
        main.go: 请使用 rick自动生成 
        go.mod:  请求确保引入了需要的组件( cmlb,cl5,北极星)
        README.md: 说明文件，必须包含：
            wiki： 需求单，方案iwiki地址，rick地址 
            监控; monitor , uls , 007, 天机阁
            发布：编译地址(Qci), 现网发布地址(织云,stke ), 测试环境发布地址( 织云，stke ) 
            文件说明：核心业务服务，需要增加一些文件说明 
            资源: 有申请特殊资源的，需要备注资源的申请地址,以及接口人 


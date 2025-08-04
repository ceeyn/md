package model

const (
	TemplateInfoDBName    = "wemeet_info"
	TemplateInfoTableName = "wemeet_meeting_template_info"
)

// TemplateInfo 模板信息
// 注意此结构体字段名和redis(rpc/rpc_redis_access.go)、dbProxy(kafka/template.go)保持一致，若修改了字段名，需要同步
type TemplateInfo struct {
	Id          			uint32 `json:"id" gorm:"primaryKey;autoIncrement"`
	TemplateId  			string `json:"template_id"`                   // 模板资源id
	Sponsor     			string `json:"sponsor"`                       // 主办方
	CoverName   			string `json:"cover_name"`                    // 封面名称
	CoverUrl    			string `json:"cover_url"   gorm:"type:text"`  // 封面路径
	Description 			string `json:"description" gorm:"type:text"`  // webinar会议详情
	CoverList   			string `json:"cover_list"   gorm:"type:text"` // 封面列表，序列化后的json结构
	WarmUpData 	 			string `json:"warm_up_data" gorm:"type:text"` // 暖场物料信息，序列化后的json结构
	MeetingId   			string `json:"meeting_id"`                    // 会议id
	AppId       			string `json:"app_id"`                        // 创建或修改者appId
	AppUid      			string `json:"app_uid"`                       // 创建或修改者appUid
}

// TableName 返回对应的mysql表名
func (TemplateInfo) TableName() string {
	return TemplateInfoTableName
}

// TemplateMqMsg ... tdmq延时消息
type TemplateMqMsg struct {
	MeetingId    	uint64   // 会议Id
	TemplateId   	string   // 模版ID
	CosId           string   // CosId
	MsgType     	int32    // 预留字段
	TryCount       	uint32   // 补偿消息重试次数
	MsgTranceId 	string   // 消息追踪ID
}
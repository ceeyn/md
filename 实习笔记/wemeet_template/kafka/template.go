package kafka

import (
	"context"
	"encoding/base64"
	"fmt"

	"meeting_template/model"
	"meeting_template/util"

	"git.code.oa.com/trpc-go/trpc-go/log"
	dbProxy "git.code.oa.com/trpcprotocol/wemeet/wemee_db_proxy_center_join_wemeet_db_proxy_center_join"
	"git.code.oa.com/trpcprotocol/wemeet/wemeet_kafka_message"
	uuid "git.code.oa.com/wesee_ugc/go.uuid"
	"github.com/golang/protobuf/proto"
)

// InsertTemplateInfoDbProxy 通过DbProxy插入TemplateInfo记录
func InsertTemplateInfoDbProxy(ctx context.Context, templateInfo model.TemplateInfo) error {
	util.ReportOne(util.DBProxyInsertTemplateInfoReq)
	uniqueId := uuid.NewV4().String()
	sqlStatement, parameters := getInsertTemplateInfoSql(ctx, templateInfo)
	sqlTplReq := &dbProxy.SqlTplReq{
		Parameters:   parameters,                                                      // sql的填充参数
		Ts:           proto.Int64(util.Now()),                                         // 时间戳，单位ms
		RequestId:    proto.String(uniqueId),                                          // 请求id，可以用uuid，确保唯一
		SqlStatement: sqlStatement,                                                    // sql语句
		ExecMode:     proto.Int32(int32(dbProxy.SqlExecMode_SQL_EXEC_MODE_STATEMENT)), // 采用原始的SQL语句执行
		UniqueId:     proto.String(uniqueId),                                          // 上报的唯一id
	}

	meetingKafkaSqlTplMsg := &wemeet_kafka_message.MeetingKafkaSqlTplMsg{
		SqlTplMsg: &dbProxy.SqlTplRequest{
			Reqs: []*dbProxy.SqlTplReq{sqlTplReq},
		},
	}
	log.InfoContextf(ctx, "InsertTemplateInfo meetingKafkaSqlTplMsg:%+v ", meetingKafkaSqlTplMsg)
	message, err := proto.Marshal(meetingKafkaSqlTplMsg)
	if err != nil {
		util.ReportOne(util.DBProxyInsertTemplateInfoFail)
		log.ErrorContextf(ctx, "InsertTemplateInfo meetingKafkaSqlTplMsg marshal err. "+
			"meetingKafkaSqlTplMsg:%+v err:%+v,", meetingKafkaSqlTplMsg, err)
		return err
	}
	err = producerProxyInstance.TemplateDBProxyProducer.Produce(ctx, []byte(templateInfo.TemplateId), message)
	if err != nil {
		util.ReportOne(util.DBProxyInsertTemplateInfoFail)
		log.ErrorContextf(ctx, "InsertTemplateInfo produce err. "+
			"meetingKafkaSqlTplMsg:%+v err:%+v,", meetingKafkaSqlTplMsg, err)
		return err
	}
	return nil
}

// UpdateTemplateInfoDbProxy 通过DbProxy修改TemplateInfo记录
func UpdateTemplateInfoDbProxy(ctx context.Context, templateInfo model.TemplateInfo) error {
	util.ReportOne(util.DBProxyUpdateTemplateInfoReq)
	uniqueId := uuid.NewV4().String()

	sqlStatement, parameters := getUpdateTemplateInfoSql(ctx, templateInfo)
	sqlTplReq := &dbProxy.SqlTplReq{
		Parameters:   parameters,                                                      // sql的填充参数
		Ts:           proto.Int64(util.Now()),                                         // 时间戳，单位ms
		RequestId:    proto.String(uniqueId),                                          // 请求id，可以用uuid，确保唯一
		SqlStatement: sqlStatement,                                                    // sql语句
		ExecMode:     proto.Int32(int32(dbProxy.SqlExecMode_SQL_EXEC_MODE_STATEMENT)), // 采用原始的SQL语句执行
		UniqueId:     proto.String(uniqueId),                                          // 上报的唯一id
	}
	meetingKafkaSqlTplMsg := &wemeet_kafka_message.MeetingKafkaSqlTplMsg{
		SqlTplMsg: &dbProxy.SqlTplRequest{
			Reqs: []*dbProxy.SqlTplReq{sqlTplReq},
		},
	}
	log.InfoContextf(ctx, "UpdateTemplateInfo meetingKafkaSqlTplMsg:%+v ", meetingKafkaSqlTplMsg)
	message, err := proto.Marshal(meetingKafkaSqlTplMsg)
	if err != nil {
		util.ReportOne(util.DBProxyUpdateTemplateInfoFail)
		log.ErrorContextf(ctx, "UpdateTemplateInfo meetingKafkaSqlTplMsg marshal err. "+
			"meetingKafkaSqlTplMsg:%+v err:%+v,", meetingKafkaSqlTplMsg, err)
		return err
	}
	err = producerProxyInstance.TemplateDBProxyProducer.Produce(ctx, []byte(templateInfo.TemplateId), message)
	if err != nil {
		util.ReportOne(util.DBProxyUpdateTemplateInfoFail)
		log.ErrorContextf(ctx, "UpdateTemplateInfo produce err. "+
			"meetingKafkaSqlTplMsg:%+v err:%+v,", meetingKafkaSqlTplMsg, err)
		return err
	}
	return nil
}

// getInsertTemplateInfoSql 根据templateInfo获取DbProxy插入所需sql及参数列表
func getInsertTemplateInfoSql(ctx context.Context, templateInfo model.TemplateInfo) (*string,
	[]*dbProxy.SqlTplParam) {

	sqlTplParams := getTemplateInfoSqlTplParams(ctx, templateInfo)
	sqlStatement := fmt.Sprintf("INSERT INTO %s.%s (", model.TemplateInfoDBName, model.TemplateInfoTableName)
	for i, sqlTplParam := range sqlTplParams {
		sqlStatement += sqlTplParam.GetName()
		if i != len(sqlTplParams)-1 {
			sqlStatement += ","
		}
	}
	sqlStatement += ") VALUES ("
	for i, _ := range sqlTplParams {
		sqlStatement += "?"
		if i != len(sqlTplParams)-1 {
			sqlStatement += ","
		}
	}
	sqlStatement += ")"
	return proto.String(sqlStatement), sqlTplParams
}

// getUpdateTemplateInfoSql 根据templateInfo获取DbProxy修改所需sql及参数列表
func getUpdateTemplateInfoSql(ctx context.Context, templateInfo model.TemplateInfo) (*string,
	[]*dbProxy.SqlTplParam) {

	sqlTplParams := getTemplateInfoSqlTplParams(ctx, templateInfo)
	sqlStatement := fmt.Sprintf("UPDATE %s.%s SET ", model.TemplateInfoDBName, model.TemplateInfoTableName)
	for i, sqlTplParam := range sqlTplParams {
		if sqlTplParam.GetName() != "template_id" {
			sqlStatement += fmt.Sprintf("%s=?", sqlTplParam.GetName())
			if i != len(sqlTplParams)-2 {
				sqlStatement += ","
			}
		}
	}
	sqlStatement += " WHERE template_id=?"
	return proto.String(sqlStatement), sqlTplParams
}

// getTemplateInfoSqlTplParams 根据templateInfo获取DbProxy修改所需参数列表
func getTemplateInfoSqlTplParams(ctx context.Context,
	templateInfo model.TemplateInfo) []*dbProxy.SqlTplParam {
	// description 为 zlib 加密后字符串，直接存mysql可能导致乱码，用base64编码后再存mysql
	base64Description := base64.StdEncoding.EncodeToString([]byte(templateInfo.Description))
	sqlTplParams := []*dbProxy.SqlTplParam{
		// 注意字段名和model.TemplateInfo保持一致，如果修改需同步
		&dbProxy.SqlTplParam{Name: proto.String("sponsor"), Value: proto.String(templateInfo.Sponsor)},
		&dbProxy.SqlTplParam{Name: proto.String("cover_name"), Value: proto.String(templateInfo.CoverName)},
		&dbProxy.SqlTplParam{Name: proto.String("cover_url"), Value: proto.String(templateInfo.CoverUrl)},
		&dbProxy.SqlTplParam{Name: proto.String("description"), Value: proto.String(base64Description)},
		&dbProxy.SqlTplParam{Name: proto.String("cover_list"), Value: proto.String(templateInfo.CoverList)},
		&dbProxy.SqlTplParam{Name: proto.String("warm_up_data"), Value: proto.String(templateInfo.WarmUpData)},
		&dbProxy.SqlTplParam{Name: proto.String("meeting_id"), Value: proto.String(templateInfo.MeetingId)},
		&dbProxy.SqlTplParam{Name: proto.String("app_id"), Value: proto.String(templateInfo.AppId)},
		&dbProxy.SqlTplParam{Name: proto.String("app_uid"), Value: proto.String(templateInfo.AppUid)},
		&dbProxy.SqlTplParam{Name: proto.String("template_id"), Value: proto.String(templateInfo.TemplateId)},
	}
	return sqlTplParams
}

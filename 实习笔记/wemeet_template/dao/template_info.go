package dao

import (
	"context"
	"encoding/base64"
	"fmt"

	"meeting_template/model"
	"meeting_template/util"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// GetTemplateInfoByTemplateId 基于 templateId 获取 TemplateInfo 记录
func GetTemplateInfoByTemplateId(ctx context.Context, templateId string, db *gorm.DB) (*model.TemplateInfo, error) {
	funcName := "GetTemplateInfoByTemplateId"
	baseLog := fmt.Sprintf("%+v(templateId:%+v)", funcName, templateId)
	util.ReportOne(util.MysqlGetTemplateInfoByTemplateIdReq) //[GetTemplateInfoMySQL]请求
	log.InfoContextf(ctx, "%+v,", baseLog)

	r := &model.TemplateInfo{TemplateId: templateId}
	if db == nil { // 非事务场景，上层传入nil
		db = DbInstance.TemplateReadonlyProxy
	}
	result := db.Debug().WithContext(ctx).Where("template_id=?", r.TemplateId).Take(r)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			util.ReportOne(util.MysqlGetTemplateInfoByTemplateIdNotFound) //[GetTemplateInfoMySQL]NotFound
		} else {
			util.ReportOne(util.MysqlGetTemplateInfoByTemplateIdFail) //[GetTemplateInfoMySQL]失败
		}
		log.ErrorContextf(ctx, "%+v error, record:%+v, err:%+v", baseLog, r, result.Error)
		return nil, result.Error
	}
	// description 为 zlib 加密后字符串，直接存mysql可能导致乱码，用base64编码后再存mysql
	description, err := base64.StdEncoding.DecodeString(r.Description)
	if err != nil {
		//[GetTemplateInfoByTemplateId]description base64解码失败
		util.ReportOne(util.MysqlGetTemplateInfoByTemplateIdDescriptionBase64DecodeFail)
		log.ErrorContextf(ctx, "%+v description base64 decode error, record:%+v, err:%+v", baseLog, r, err)
		return nil, errors.Errorf("%+v, err:[%+v]", baseLog, result.Error)
	}
	r.Description = string(description)
	util.ReportOne(util.MysqlGetTemplateInfoByTemplateIdSucc) //[GetTemplateInfoMySQL]成功
	log.InfoContextf(ctx, "%+v success, record:%+v", baseLog, r)
	return r, nil
}

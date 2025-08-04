package exception

import (
	"net/http"

	"git.code.oa.com/trpc-go/trpc-go/errs"
)

const (
	HttpResponseOk = http.StatusOK
	// 参数错误
	HttpResParamError = 20000 // 参数错误

	// 内部服务错误
	HttpResServerError       = 21001 // 内部服务错误
	HttpDecodeError          = 21004 // base64解析失败
	HttpQueryIsecResultError = 21005 // base64解析失败

	// cos错误
	HttpResUploadCosError  = 22001 // 上传数据到cos错误
	HttpResCosAuthError    = 22002 // cos下载失败失败
	HttpResCosTempUrlError = 20003 // 获取cos临时链接失败
	OidbRes

	// mysql调用错误
	HttpResSqlError = 23001 // sql调用失败

	// 鉴黄调用错误
	HttpResIsecCheckError    = 24001 // 图片鉴定调用失败
	HttpResIsecCallBackError = 24002 // 图片鉴定调用失败

	// 敏感词校验
	HtmlCheckSensitive = 25001 // 校验敏感词失败

	// 付费企业调用错误
	HtmlQueryCorpFree = 26001 // 付费企业查询失败

	// redis失败
	HttpResRedisServerError = 27001
)

var (
	ErrRequestParamsNotExists = errs.New(HttpResParamError, "request params not exists")
	ErrRequestCosIdNotExists  = errs.New(HttpResParamError, "request cos_id not exists")
	ErrRequestAppIdNotExists  = errs.New(HttpResParamError, "request app_id not exists")
	ErrRequestAppUidNotExists = errs.New(HttpResParamError, "request app_uid not exists")
)

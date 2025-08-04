package util

import (
	"git.code.oa.com/going/attr"
)

var GetHealthReq int = 35620175

var CreateTemplateReq int = 35620176
var CreateTemplateFail int = 35620177

var GetTemplateInfoReq int = 35620178
var GetTemplateInfoFail int = 35620179

var UpdateTemplateReq int = 35620180
var UpdateTemplateFail int = 35620181

var RedisGetTemplateInfoReq int = 35620182
var RedisGetTemplateInfoFail int = 35620183
var RedisSetTemplateInfoReq int = 35620185
var RedisSetTemplateInfoFail int = 35620186
var RedisGetNeedCheckIdsFail int = 35620187
var RedisSetNeedCheckIdFail int = 35620188
var RedisSetExpireReq int = 35620189
var RedisSetExpireFail int = 35620190

var BatchQueryTempUrl int = 35620191
var BatchQueryTempUrlFail int = 35620192

var GetTextModerationReq int = 35620193
var GetTextModerationFail int = 35620194
var QueryMeetingSensitiveWords int = 35620195
var QueryMeetingSensitiveWordsFail int = 35620196

var TransCosIdToUrlFail int = 35620197
var Base64DecodeDescriptionFail int = 35620198
var Base64DecodeCoverUrlFail int = 35620199
var HtmlParseDescriptionFail int = 35620200
var ReplaceDiscriptionImgFail int = 35620201
var ReplaceDiscriptionSensitiveFail int = 35620202
var ReplaceDiscriptionTextFail int = 35620203
var SponsorCoverNameHitSensitive int = 35620204
var CheckTextEmpty int = 35663188
var CoverUrlCosIdInvalid int = 35676663
var DescriptionCosIdInvalid int = 35676664
var CompressDescriptionFail int = 35684348
var DeCompressDescriptionFail int = 35684349
var CoverListRawCuttedNumDiff int = 35702559
var GetWarmUpDataReq int = 35738263
var GetWarmUpDataFail int = 35738264
var JsonUnmarshalFail int = 35738265
var JsonMarshalFail int = 35738266
var PaseUintMeetingIdFail int = 35738267
var PaseUintAppIdFail int = 35738268
var DescriptionHitSensitive int = 35738269
var VideoCensorCallback int = 35738270
var VideoCensorCallbackFail int = 35738271
var VideoCensorCallbakTraceIdEmpty = 35738272
var VideoCensorCallbakTraceIdFormatError int = 35738273
var VideoCensorCallbakWarmupDataEmpty int = 35738274
var WarmUpVideoDataFormatError int = 35738275
var CensorCallbackCosIdExpired int = 35738276
var CensorCallbackAuditPass int = 35738277
var CensorCallbackAuditFail int = 35738278
var GetVidoReq int = 35738279
var GetVidoReqFail int = 35738280
var UpdateWarmUpUserNotify int = 35738281
var UpdateWarmUpUserNotifyFail int = 35738282
var GetWholeBrandInfoReqAttr = 35737413
var GetWholeBrandInfoFailAttr = 35737414
var GetWholeBrandInfoSuccAttr = 35737415
var HandleGetWholeBrandInfoGetMeetingInfoFailAttr = 35737416
var HandleGetWholeBrandInfoGetTemplateFailAttr = 35737417
var GetTemplateIdEmpty int = 35741712
var UpdateTemplateIdEmpty int = 35741713

var MysqlGetTemplateInfoByTemplateIdReq = 35810010
var MysqlGetTemplateInfoByTemplateIdSucc = 35810020
var MysqlGetTemplateInfoByTemplateIdFail = 35810030
var MysqlGetTemplateInfoByTemplateIdNotFound = 35810040
var MysqlGetTemplateInfoByTemplateIdDescriptionBase64DecodeFail = 35810050

var GetTemplateInfoSingleFlightCall = 35820010
var GetTemplateInfoSingleFlightCallShared = 35820020
var GetTemplateInfoTemplateIdIllegal = 35820030
var GetTemplateInfoFromRedisSucc = 35820040
var GetTemplateInfoFromRedisFail = 35820050
var GetTemplateInfoFromMySQLSucc = 35820060
var GetTemplateInfoFromMySQLNotFound = 35820070
var GetTemplateInfoFromMySQLFail = 35820080

var CreateTemplateInfoReq = 35830010
var UpdateTemplateInfoReq = 35830020
var UpdateTemplateTemplateIdIllegal = 35830030

var DBProxyInsertTemplateInfoReq int = 35840010
var DBProxyInsertTemplateInfoFail int = 35840020
var DBProxyUpdateTemplateInfoReq int = 35840030
var DBProxyUpdateTemplateInfoFail int = 35840040

var UpdateMeetingShareItemReq int = 35928392
var UpdateMeetingShareItemFail int = 35928393
var UpdateMeetingShareItemInValidReq = 35928394
var UpdateMeetingShareItemGetMeetingInfoFail = 35928395
var UpdateMeetingShareItemAuthFail = 35928396
var UpdateMeetingShareItemSetMeetingInfoFail = 35928397

// monitor累计量上报
func ReportOne(monitorId int) {
	if monitorId > 0 {
		attr.AttrAPI(monitorId, 1)
	}
}

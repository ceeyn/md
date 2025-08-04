package es

import (
	"context"
	"encoding/base64"
	"errors"
	"reflect"
	"time"

	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go/log"
	cachePb "git.code.oa.com/trpcprotocol/wemeet/common_meeting_cache"
	"github.com/olivere/elastic"
)

const (
	ESDocType                     = "doc" //es6.X 保存es文档的时候需要指定Type
	Itinerary_NAME                = 1
	Itinerary_INTRODUCTION        = 2
	MEETING_ORDER_TIME_TYPE       = 1
	MEETING_ORDER_START_TIME_TYPE = 2
	MEETING_ORDER_END_TIME_TYPE   = 3
)

// getEsClient 获取链接ES的客户端
func getEsClient(ctx context.Context) (*elastic.Client, error) {
	esClient := EsClient
	if esClient == nil && ReConnectES(ctx) == false {
		log.ErrorContext(ctx, "[elasticsearch] reconnect fail")
		return nil, errors.New("[elasticsearch] reconnect fail")
	}

	if EsClient == nil {
		log.ErrorContext(ctx, "[elasticsearch] client is nil")
		return nil, errors.New("[elasticsearch] client is nil")
	}

	return EsClient, nil
}

// SaveItineraryToES 单条保存
func SaveItineraryToES(ctx context.Context, info *Itinerary, docId string, meetingInfo *cachePb.MeetingInfo) error {

	if info == nil {
		return errors.New("Itinerary nil")
	}
	esClient, err := getEsClient(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "getEsClient fail,error:%+v", err)
		return err
	}
	now := time.Now().Format("2006-01-02 15:04:05")

	info.CreateTime = now
	info.UpdateTime = now
	// 会议预定时间
	orderTm := time.Unix(int64(meetingInfo.GetUint32OrderTime()), 0)
	strOrderTime := orderTm.Format("2006-01-02 15:04:05")
	// 会议预定开始时间
	orderStartTm := time.Unix(int64(meetingInfo.GetUint32OrderStartTime()), 0)
	strOrderStartTime := orderStartTm.Format("2006-01-02 15:04:05")
	// 会议预定结束时间
	orderEndTm := time.Unix(int64(meetingInfo.GetUint32OrderEndTime()), 0)
	strOrderEndTime := orderEndTm.Format("2006-01-02 15:04:05")
	//会议主题
	meetingSubject, _ := base64.StdEncoding.DecodeString(string(meetingInfo.GetBytesMeetingSubject()))
	info.MeetingOrderTime = strOrderTime
	info.MeetingOrderStartTime = strOrderStartTime
	info.MeetingOrderEndTime = strOrderEndTime
	info.MeetingSubject = string(meetingSubject)
	esRes, err := esClient.Index().Index(ESItineraryIndex).Type(ESDocType).Id(docId).BodyJson(info).Do(ctx)
	if err != nil {
		attr.AttrAPI(36336994, 1)
		log.ErrorContextf(ctx, "SaveItineraryIntroductionToES failed, meetingId:%+v, err:%+v",
			info.MeetingId, err)
		return err
	}
	log.InfoContextf(ctx, "[elasticsearch] insert data success,index:%v, id:%v, value:%+v, res:%+v",
		ESItineraryIndex, info.MeetingId, info, esRes)
	return nil
}

//ModifyItineraryToES ...
func ModifyItineraryToES(ctx context.Context, info *Itinerary, docId string) error {
	if info == nil {
		return errors.New("Itinerary nil")
	}
	esClient, err := getEsClient(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "getEsClient fail,error:%+v", err)
		return err
	}
	esRes, err := esClient.Update().Index(ESItineraryIndex).Type(ESDocType).Id(docId).Doc(info).Do(ctx)
	if err != nil {
		attr.AttrAPI(36336995, 1)
		log.ErrorContextf(ctx, "ModifyItineraryToES failed, meetingId:%+v, err:%+v", info.MeetingId, err)
		return err
	}
	log.InfoContextf(ctx, "[elasticsearch] modify data success,index:%v,id:%v,value:%+v,res:%+v",
		ESItineraryIndex, info.MeetingId, info, esRes)
	return nil
}

// DelItineraryToES ...
func DelItineraryToES(ctx context.Context, docId string) error {
	if docId == "" {
		return errors.New("docId empty")
	}
	esClient, err := getEsClient(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "getEsClient fail, error:%+v", err)
		return err
	}
	_, err = esClient.Delete().Index(ESItineraryIndex).Type(ESDocType).Id(docId).Refresh("true").Do(ctx)
	if err != nil {
		attr.AttrAPI(36336996, 1)
		log.ErrorContextf(ctx, "DelItineraryToES failed, err: %v", err)
		return err
	}

	log.InfoContextf(ctx, "[elasticsearch] DelItineraryToES success,docId:%+v, index:%v",
		docId, ESItineraryIndex)
	return nil
}

// UpsertIntroductionToES ... 修改文档（不存在则插入）
func UpsertIntroductionToES(ctx context.Context, docId string, introduction *Introduction) error {
	if docId == "" {
		return errors.New("docId empty")
	}
	esClient, err := getEsClient(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "getEsClient fail, error:%+v", err)
		return err
	}
	esRes, err := esClient.Update().Index(ESIntroductionIndex).Type(ESDocType).
		Id(docId).Doc(introduction).Upsert(introduction).Do(ctx)
	if err != nil {
		attr.AttrAPI(36336997, 1)
		log.ErrorContextf(ctx, "UpsertIntroductionToES failed, meetingId:%+v, err:%+v", docId, err)
		return err
	}
	log.InfoContextf(ctx, "[elasticsearch] UpsertIntroductionToES data success, "+
		"index:%v, id:%v, value:%+v, res:%+v", ESIntroductionIndex, docId, introduction, esRes)
	return nil
}

// DelIntroductionToES ...
func DelIntroductionToES(ctx context.Context, docId string) error {
	if docId == "" {
		return errors.New("docId empty")
	}
	esClient, err := getEsClient(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "getEsClient fail, error:%+v", err)
		return err
	}
	_, err = esClient.Delete().Index(ESIntroductionIndex).Type(ESDocType).Id(docId).Refresh("true").Do(ctx)
	if err != nil {
		attr.AttrAPI(36336998, 1)
		log.ErrorContextf(ctx, "DelIntroductionToES failed, err: %v", err)
		return err
	}

	log.InfoContextf(ctx, "[elasticsearch] DelIntroductionToES success,docId:%+v, index:%v",
		docId, ESIntroductionIndex)
	return nil
}

// FuzzyQueryItineraryInfoFromES 根据关键字模糊查询日程
func FuzzyQueryItineraryInfoFromES(ctx context.Context, searchKey string,
	pageNum, pageSize int, searchType uint32, searchTimeType uint32,
	startTime, endTime string) (uint32, []*Itinerary, error) {

	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 50
	}

	log.InfoContextf(ctx, "FuzzyQueryItineraryNameFromES, searchKey:%+v, pageSize:%+v,"+
		"page:%+v", searchKey, pageSize, pageNum)

	var name string
	if searchType == Itinerary_NAME {
		name = "itinerary_name"
	} else if searchType == Itinerary_INTRODUCTION {
		name = "itinerary_introduction"
	}
	// 时间检索类型
	rangeType := "meeting_order_time"
	if searchTimeType == MEETING_ORDER_START_TIME_TYPE { //会议预定开始时间
		rangeType = "meeting_order_start_time"
	} else if searchTimeType == MEETING_ORDER_END_TIME_TYPE { //会议预定结束时间
		rangeType = "meeting_order_end_time"
	}

	boolQuery := elastic.NewBoolQuery()
	nameMatchPhrasePrefixQuery := elastic.NewMatchPhrasePrefixQuery(name, searchKey)
	nameWildCardQuery := elastic.NewWildcardQuery(name, "*"+searchKey+"*")
	shouldQuery := elastic.NewBoolQuery().Should(nameMatchPhrasePrefixQuery, nameWildCardQuery)
	//时间范围range
	rangeQuery := elastic.NewRangeQuery(rangeType).
		Gte(startTime). //起始时间
		Lte(endTime) //终止时间
	boolQuery.Must(shouldQuery, rangeQuery)

	searchResult, total, err := QueryItineraryFromEsByCondition(ctx, boolQuery, "meeting_order_start_time",
		(pageNum-1)*pageSize, pageSize)
	if err != nil {
		attr.AttrAPI(36336999, 1)
		log.ErrorContextf(ctx, "FuzzyQueryItineraryNameFromES QueryFromEsByCondition fail,err:%+v", err)
		return 0, nil, err
	}

	itineraryList := []*Itinerary{}
	for _, v := range searchResult.Each(reflect.TypeOf(&Itinerary{})) {
		vs := v.(*Itinerary)
		itineraryList = append(itineraryList, vs)
	}
	log.InfoContextf(ctx, "FuzzyQueryItineraryNameFromES get result itineraryList:%+v, searchKey:%+v",
		itineraryList, searchKey)
	// 返回折叠后的total
	return total, itineraryList, nil
}

// QueryItineraryFromEsByCondition 从es中查询日程数据
func QueryItineraryFromEsByCondition(ctx context.Context, query elastic.Query, sortBy string, from,
	size int) (*elastic.SearchResult, uint32, error) {
	esClient, err := getEsClient(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "getEsClient fail,err:%+v", err)
		return nil, 0, err
	}
	log.InfoContextf(ctx, "QueryFromEsByCondition index:%s", ESItineraryIndex)

	//设置统计字段
	aggs := elastic.NewCardinalityAggregation().Field("meeting_id")
	//注意一下折叠功能必须是keyword关键字
	collapseData := elastic.NewCollapseBuilder("meeting_id")

	esRes, err := esClient.
		Search(ESItineraryIndex). // 日程的index
		Type(ESDocType). // es6注意Type
		Query(query).
		Collapse(collapseData). // 折叠
		Aggregation("total", aggs). //去重total统计
		Sort(sortBy, true).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		attr.AttrAPI(36337000, 1)
		log.ErrorContextf(ctx, "QueryFromEsByCondition from es failed, error: %+v", err)
		return nil, 0, err
	}

	// 使用Cardinality函数和前面定义的聚合条件名称，查询结果
	var total uint32
	agg, found := esRes.Aggregations.Cardinality("total")
	if found {
		total = uint32(*agg.Value)
		log.InfoContextf(ctx, "QueryItineraryFromEsByCondition get total:%+v, agg.Value:%+v", total, *agg.Value)
	}

	return esRes, total, nil
}

// FuzzyQueryIntroductionFromES 根据关键字 模糊查询 会议介绍
func FuzzyQueryIntroductionFromES(ctx context.Context, searchKey string,
	pageNum, pageSize int, searchTimeType uint32, startTime, endTime string) (uint32, []*Introduction, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 50
	}

	log.InfoContextf(ctx, "FuzzyQueryIntroductionFromES, searchKey:%+v, pageSize:%+v,"+
		"page:%+v", searchKey, pageSize, pageNum)

	var name = "meeting_introduction"
	// 时间检索类型
	rangeType := "meeting_order_time"
	if searchTimeType == MEETING_ORDER_START_TIME_TYPE { //会议预定开始时间
		rangeType = "meeting_order_start_time"
	} else if searchTimeType == MEETING_ORDER_END_TIME_TYPE { //会议预定结束时间
		rangeType = "meeting_order_end_time"
	}

	boolQuery := elastic.NewBoolQuery()
	nameMatchPhrasePrefixQuery := elastic.NewMatchPhrasePrefixQuery(name, searchKey)
	nameWildCardQuery := elastic.NewWildcardQuery(name, "*"+searchKey+"*")
	shouldQuery := elastic.NewBoolQuery().Should(nameMatchPhrasePrefixQuery, nameWildCardQuery)
	//时间范围range
	rangeQuery := elastic.NewRangeQuery(rangeType).
		Gte(startTime). //起始时间
		Lte(endTime) //终止时间
	boolQuery.Must(shouldQuery, rangeQuery)

	searchResult, err := QueryIntroductionFromEsByCondition(ctx, boolQuery, "meeting_order_start_time",
		(pageNum-1)*pageSize, pageSize)

	if err != nil {
		attr.AttrAPI(36337001, 1)
		log.ErrorContextf(ctx, "FuzzyQueryIntroductionFromES QueryIntroductionFromEsByCondition fail,err:%+v",
			err)
		return 0, nil, err
	}

	introductionList := []*Introduction{}
	for _, v := range searchResult.Each(reflect.TypeOf(&Introduction{})) {
		vs := v.(*Introduction)
		introductionList = append(introductionList, vs)
	}
	log.InfoContextf(ctx, "FuzzyQueryIntroductionFromES get result introductionList:%+v, searchKey:%+v",
		introductionList, searchKey)

	return uint32(searchResult.Hits.TotalHits), introductionList, nil
}

// QueryIntroductionFromEsByCondition 从es中查询 会议介绍数据
func QueryIntroductionFromEsByCondition(ctx context.Context, query elastic.Query, sortBy string, from,
	size int) (*elastic.SearchResult, error) {
	esClient, err := getEsClient(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "getEsClient fail,err:%+v", err)
		return nil, err
	}
	log.InfoContextf(ctx, "QueryIntroductionFromEsByCondition index:%s", ESIntroductionIndex)

	esRes, err := esClient.
		Search(ESIntroductionIndex). // 会议介绍的Index
		Type(ESDocType). // es6注意Type
		Query(query).
		Sort(sortBy, true).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		attr.AttrAPI(36337002, 1)
		log.ErrorContextf(ctx, "QueryIntroductionFromEsByCondition from es failed, error: %+v", err)
		return nil, err
	}

	return esRes, nil
}

//AccurateSearchMeetingIntroduction ...精确搜索 根据会议id搜会议介绍
func AccurateSearchMeetingIntroduction(ctx context.Context, meetingId string) ([]*Introduction, error) {
	esClient, err := getEsClient(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "getEsClient fail,err:%+v", err)
		return nil, err
	}
	introductionList := []*Introduction{}

	// 创建 term 查询条件，用于精确查询
	termQuery := elastic.NewTermQuery("meeting_id", meetingId)
	searchResult, err := esClient.Search().
		Index(ESIntroductionIndex). // 设置索引名 会议介绍的索引
		Type(ESDocType). // docType
		Query(termQuery). // 设置查询条件
		Sort("meeting_order_start_time", true). // 设置排序字段，根据 meeting_order_start_time 字段升序排序
		Do(ctx) // 执行请求

	if err != nil {
		attr.AttrAPI(36337003, 1)
		log.ErrorContextf(ctx, "AccurateSearchMeetingIntroduction TermQuery ES failed. meetingId:%+v, err:%+v",
			meetingId, err)
		return introductionList, err
	}

	log.InfoContextf(ctx, "AccurateSearchMeetingIntroduction get %+v results, results:%+v",
		searchResult.Hits.TotalHits, searchResult.Each(reflect.TypeOf(&Introduction{})))

	for _, v := range searchResult.Each(reflect.TypeOf(&Introduction{})) {
		vs := v.(*Introduction)
		introductionList = append(introductionList, vs)
	}
	return introductionList, nil
}

//AccurateSearchMeetingAllItinerary ... 精确搜索 根据会议id搜这场会议所有日程
func AccurateSearchMeetingAllItinerary(ctx context.Context, meetingId string) ([]*Itinerary, error) {
	esClient, err := getEsClient(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "getEsClient fail,err:%+v", err)
		return nil, err
	}
	itineraryList := []*Itinerary{}

	// 创建 term 查询条件，用于精确查询
	termQuery := elastic.NewTermQuery("meeting_id", meetingId)
	searchResult, err := esClient.Search().
		Index(ESItineraryIndex). // 设置索引名 日程的索引
		Type(ESDocType). // docType
		Query(termQuery). // 设置查询条件
		Sort("meeting_order_start_time", true). // 设置排序字段，根据 meeting_order_start_time 字段升序排序
		Do(ctx) // 执行请求

	if err != nil {
		attr.AttrAPI(36337004, 1)
		log.ErrorContextf(ctx, "AccurateSearchMeetingAllItinerary TermQuery ES failed. meetingId:%+v, err:%+v",
			meetingId, err)
		return itineraryList, err
	}

	log.InfoContextf(ctx, "AccurateSearchMeetingAllItinerary get total results:%+v , results:%+v",
		searchResult.Hits.TotalHits, searchResult.Each(reflect.TypeOf(&Itinerary{})))

	for _, v := range searchResult.Each(reflect.TypeOf(&Itinerary{})) {
		vs := v.(*Itinerary)
		itineraryList = append(itineraryList, vs)
	}
	return itineraryList, nil
}

//BatchSearchMeetingIntroduction ... 批量搜索会议介绍
func BatchSearchMeetingIntroduction(ctx context.Context, strIds []string) ([]*Introduction, error) {
	esClient, err := getEsClient(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "getEsClient fail,err:%+v", err)
		return nil, err
	}
	introductionList := []*Introduction{}
	searchResult, err := esClient.Search(ESIntroductionIndex).Type(ESDocType).Query(
		elastic.NewIdsQuery().Ids(strIds...)).Size(len(strIds)).Do(ctx)

	if err != nil {
		attr.AttrAPI(36337005, 1)
		return introductionList, err
	}

	if searchResult.TotalHits() == 0 {
		return introductionList, nil
	}
	for _, v := range searchResult.Each(reflect.TypeOf(&Introduction{})) {
		vs := v.(*Introduction)
		introductionList = append(introductionList, vs)
	}
	return introductionList, nil
}

//BatchSearchMeetingItineraryList ...批量搜索会议日程
func BatchSearchMeetingItineraryList(ctx context.Context, strIds []string) ([]*Itinerary, error) {
	esClient, err := getEsClient(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "getEsClient fail,err:%+v", err)
		return nil, err
	}
	itineraryList := []*Itinerary{}

	interfaceIds := []interface{}{}
	for _, strId := range strIds {
		interfaceIds = append(interfaceIds, strId)
	}
	// 创建 term 查询条件，用于精确查询
	termQuery := elastic.NewTermsQuery("meeting_id", interfaceIds...)
	searchResult, err := esClient.Search().
		Index(ESItineraryIndex). // 设置索引名
		Type(ESDocType). // Type
		Query(termQuery). // 设置查询条件
		Sort("meeting_order_start_time", true). // 设置排序字段，根据 meeting_order_start_time 字段升序排序
		Do(ctx) // 执行请求

	if err != nil {
		attr.AttrAPI(36337006, 1)
		return itineraryList, err
	}

	log.InfoContextf(ctx, "BatchSearchMeetingItineraryList  get %+v results, results:%+v",
		searchResult.Hits.TotalHits, searchResult.Each(reflect.TypeOf(&Itinerary{})))

	for _, v := range searchResult.Each(reflect.TypeOf(&Itinerary{})) {
		vs := v.(*Itinerary)
		itineraryList = append(itineraryList, vs)
	}
	return itineraryList, nil
}
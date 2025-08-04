package es

import (
	"context"
	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-go"
	"meeting_template/config/config_rainbow"
	"net/http"
	"sync"
	"time"


	"git.code.oa.com/trpc-go/trpc-go/log"

	"github.com/olivere/elastic"
)

// ES 连接 httpClient 使用参数
const (
	//MaxIdleConns 最大节点连接数
	MaxIdleConns = 2000

	//MaxIdleConnsPerHost 最大空闲连接
	MaxIdleConnsPerHost = 2000

	//IdleConnTimeout 节点连接超时时间
	IdleConnTimeout = 90 * time.Second

	//MaxBulkSize BulkSize 5M
	MaxBulkSize = 5 * 1024 * 1024

	//MaxReadBufferSize BulkSize 5M
	MaxReadBufferSize = 10 * 1024 * 1024
)

// ES的client
var (
	EsClient *elastic.Client
	// ESCltLock elasticsearch客户端重连互斥锁
	ESCltLock *sync.Mutex
)

var (
	ESItineraryIndex string
	ESIntroductionIndex string
)

// Init ES初始化
func Init() {
	ESItineraryIndex = WebinarItineraryIndex + trpc.GlobalConfig().Global.Namespace           //会议日程index
	ESIntroductionIndex = WebinarIntroductionIndex + trpc.GlobalConfig().Global.Namespace     //会议介绍index
	EsClient = NewEsClient(context.Background())
	initIndex(context.Background())
}

// initIndex ...
func initIndex(ctx context.Context) {
	// 会议日程的索引
	itineraryExist, err := EsClient.IndexExists(ESItineraryIndex).Do(ctx)
	if err != nil {
		attr.AttrAPI(36337007,1)
		log.Errorf("itineraryIndexExists fail, err=%v", err)
	}
	if !itineraryExist {
		_, err := EsClient.CreateIndex(ESItineraryIndex).Body(esItineraryMapping).Do(ctx)
		if err != nil {
			attr.AttrAPI(36337008,1)
			log.Errorf("itinerary es CreateIndex :%s fail, err:%v", ESItineraryIndex, err)
		}
	}
	// 会议介绍的索引
	introductionExist, err := EsClient.IndexExists(ESIntroductionIndex).Do(ctx)
	if err != nil {
		attr.AttrAPI(36337009,1)
		log.Errorf("introductionIndexExists fail, err=%v", err)
	}
	if !introductionExist {
		_, err := EsClient.CreateIndex(ESIntroductionIndex).Body(esIntroductionMapping).Do(ctx)
		if err != nil {
			attr.AttrAPI(36337010,1)
			log.Errorf("introduction es CreateIndex :%s fail, err:%v", ESIntroductionIndex, err)
		}
	}
}

// NewEsClient 创建连接ES客户端
func NewEsClient(ctx context.Context) *elastic.Client {
	esConfig := config_rainbow.GetEsConfConfig()
	client, err := elastic.NewClient(
		elastic.SetURL(esConfig.EsHost),
		elastic.SetBasicAuth(esConfig.EsUser, esConfig.EsPwd),
		elastic.SetSniff(false),
		elastic.SetHttpClient(&http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        MaxIdleConns,        // 最大连接数,默认0无穷大
				MaxIdleConnsPerHost: MaxIdleConnsPerHost, // 对每个host的最大连接数量(MaxIdleConnsPerHost<=MaxIdleConns)
				IdleConnTimeout:     IdleConnTimeout,     // 多长时间未使用自动关闭连接
				ReadBufferSize:      MaxReadBufferSize,
				WriteBufferSize:     MaxBulkSize,
			},
		}),
	)

	if err != nil {
		attr.AttrAPI(36337011,1)
		log.Errorf("load es fail, err=%v", err)
	}
	attr.AttrAPI(36337012,1)
	log.InfoContextf(ctx, "elasticsearch new client success, ESScheduleAndIntroIndex:%+v",
		client, ESItineraryIndex)

	return client
}

// ReConnectES elasticsearch客户端重连单例模式
func ReConnectES(ctx context.Context) bool {
	ESCltLock.Lock()
	defer ESCltLock.Unlock()
	if EsClient == nil {
		EsClient = NewEsClient(ctx)
		if EsClient == nil {
			return false
		}
	}
	return true
}

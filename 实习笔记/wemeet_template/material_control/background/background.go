package background

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"meeting_template/material_control/cache"
	"meeting_template/material_control/rpc"
	"meeting_template/util"

	"git.code.oa.com/going/attr"
	"git.code.oa.com/trpc-go/trpc-database/redis"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/golang/protobuf/proto"
)

const (
	PicStatusNormal       = 0
	PicStatusAuditing     = 1 //待审核
	PicStatusAuditFail    = 2
	PicStatusAuditTimeOut = 3
	PicStatusDefault      = 4
)

//BackgroundImp ..
type BackgroundImp struct {
	redisProxy redis.Client
}

//NewBackground ..
func NewBackground() cache.Backgrounds {
	imp := &BackgroundImp{
		redisProxy: redis.NewClientProxy("trpc.wemeet.wemeet_template.redis_storage"),
	}
	_, err := imp.redisProxy.Do(context.Background(), "ping")
	if err != nil {
		panic(err)
	}

	return imp
}

//GetBackgroundIndex 生成索引
func (b *BackgroundImp) GetBackgroundIndex(ctx context.Context) (index int64) {

	reply, err := redis.Int64(b.redisProxy.Do(ctx, "INCR", MakeBackgroundIndexKey()))
	if err != nil {
		log.ErrorContextf(ctx, "BackgroundIndex Incr err:%+v", err)
		reply, err = redis.Int64(b.redisProxy.Do(ctx, "INCR", MakeBackgroundIndexKey()))
		if err != nil {
			metrics.IncrCounter("BackgroundIndex Incr again", 1)
			log.ErrorContextf(ctx, "BackgroundIndex Incr again err:%+v", err)
			//失败使用纳秒数据
			index = util.Now()
		}
	} else {
		curTime := time.Now().Unix()
		index = int64(curTime<<30 | reply)
	}

	log.InfoContextf(ctx, "GetBackgroundIndex:%v", index)
	return index
}

//SetBackgroundList 设置虚拟背景
func (b *BackgroundImp) SetBackgroundList(ctx context.Context, background []*pb.BackgroundInfo) error {
	dataMap := make(map[string]string, 0)
	for _, val := range background {
		log.InfoContextf(ctx, "background:%+v", background)
		buf, err := proto.Marshal(val)
		if err != nil {
			metrics.IncrCounter("BackgroundInfo proto.Marshal", 1)
			log.ErrorContextf(ctx, "BackgroundInfo proto.Marshal val:%+v, err:%v", val, err)
			continue
		}

		key := MakeBackgroundInfoKey(val.GetInt64BackgroundId())
		dataMap[key] = string(buf)
	}

	_, err := b.redisProxy.Do(ctx, "MSET", redis.Args{}.AddFlat(dataMap)...)
	if err != nil {
		metrics.IncrCounter("SetBackgroundList fail", 1)
		log.ErrorContextf(ctx, "SetBackgroundList fail cmd=MSET, err:%+v", err)
		return err
	}
	metrics.IncrCounter("SetBackgroundList Succ", 1)
	return nil
}

//SetCache180DaysExpireTimeDuration 设置过期时间
func (b *BackgroundImp) SetCache180DaysExpireTimeDuration(ctx context.Context, key string, startTime uint32) error {

	if startTime == 0 {
		metrics.IncrCounter("SetExpireTime Invalid", 1)
		log.ErrorContextf(ctx, "SetExpireTime Invalid, key:%v, startTime:%v", key, startTime)
		return nil
	}

	_, err := b.redisProxy.Do(ctx, "EXPIRE", key, util.Get180DaysExpireTimeDuration(startTime))
	if err != nil {
		attr.AttrAPI(35434438, 1) //设置缓存过期时间失败
		log.ErrorContextf(ctx, "SetCache180DaysExpireTimeDuration error, cmd = EXPIRE, key:[%v],  err = %v",
			key, err)
		return err
	}
	metrics.IncrCounter("SetCache180DaysExpireTimeDuration Succ", 1)
	return nil
}

//GetBackground ..
func (b *BackgroundImp) GetBackground(ctx context.Context, ids []int64) (background []*pb.BackgroundInfo, err error) {

	strKeys := make([]string, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		key := MakeBackgroundInfoKey(id)
		strKeys = append(strKeys, key)
	}
	if len(strKeys) == 0 {
		log.InfoContextf(ctx, "GetBackground len(strKeys)==0")
		return nil, nil
	}

	reply, err := redis.Strings(b.redisProxy.Do(ctx, "MGET", redis.Args{}.AddFlat(strKeys)...))
	if err != nil {
		metrics.IncrCounter("GetBackground fail", 1)
		log.ErrorContextf(ctx, "GetBackground fail, cmd=MGET, ids:%+v, err:%+v", ids, err)
		return nil, err
	}

	for _, val := range reply {
		data := &pb.BackgroundInfo{}
		err := proto.Unmarshal([]byte(val), data)
		if err != nil {
			metrics.IncrCounter("BackgroundInfo proto.Unmarshal err", 1)
			log.ErrorContextf(ctx, "BackgroundInfo proto.Unmarshal err:%+v", err)
			continue
		}

		//兼容历史默认背景图七彩石没有配cosId问题
		if data.GetStrPicId() == "" && strings.HasSuffix(data.GetStrPicUrl(), "default_1.png") {
			data.StrPicId = proto.String("background/default/default_1.png")
			defBgList := []*pb.BackgroundInfo{data}
			b.SetBackgroundList(ctx, defBgList)
		}

		background = append(background, data)
	}

	//todo:图片根据有效期获取,后期可优化
	err = b.GetTempImageUrlByID(ctx, background)
	if err != nil {
		metrics.IncrCounter("GetTempImageUrlByID fail", 1)
		log.ErrorContextf(ctx, "b.GetTempImageUrlByID fail, err:%+v", err)
		return nil, err
	}

	return background, nil
}

//SetMeetBackgroundSortSet ..
func (b *BackgroundImp) SetMeetBackgroundSortSet(ctx context.Context, meetingID uint64, backgroundID []int64) error {

	key := MakeMeetBackgroundKey(meetingID)
	IDMap := make(map[int64]string, 0)
	for i := 0; i < len(backgroundID); i++ {
		strID := strconv.FormatInt(backgroundID[i], 10)
		index := time.Now().Unix()
		IDMap[index+int64(i)] = strID
	}

	_, err := b.redisProxy.Do(ctx, "ZADD", redis.Args{}.Add(key).AddFlat(IDMap)...)
	if err != nil {
		metrics.IncrCounter("SetMeetBackgroundSortSet fail", 1)
		log.ErrorContextf(ctx, "SetMeetBackgroundSortSet fail cmd=ZAdd, err:%+v", err)
		return err
	}

	metrics.IncrCounter("SetMeetBackgroundSortSet Succ", 1)
	return nil
}

//GetMeetBackgroundSortSet ..
func (b *BackgroundImp) GetMeetBackgroundSortSet(ctx context.Context, meetingID uint64) (iDs []int64, err error) {

	key := MakeMeetBackgroundKey(meetingID)

	reply, err := redis.Strings(b.redisProxy.Do(ctx, "ZRANGEBYSCORE", key, "-inf", "+inf"))
	if err != nil {
		metrics.IncrCounter("GetMeetBackgroundSortSet fail", 1)
		log.ErrorContextf(ctx, "GetMeetBackgroundSortSet fail meetID:%v, err:%+v", meetingID, err)
		return nil, err
	}

	if len(reply) == 0 {
		metrics.IncrCounter("GetMeetBackgroundSortSet empty", 1)
		log.InfoContextf(ctx, "GetMeetBackgroundSortSet empty, meetID:%v", meetingID)
		return nil, nil
	}

	for _, s := range reply {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			metrics.IncrCounter("strID strconv.ParseInt fail", 1)
			log.ErrorContextf(ctx, "strID strconv.ParseInt fail, id:%v, err:%+v", id, err)
			continue
		}
		iDs = append(iDs, id)
	}

	metrics.IncrCounter("GetMeetBackgroundSortSet Succ", 1)
	return iDs, nil
}

//SetDefBackgroundSortSet 设置默认背景图zset
func (b *BackgroundImp) SetDefBackgroundSortSet(ctx context.Context, backgroundID []int64) error {

	if len(backgroundID) == 0 {
		log.InfoContextf(ctx, "SetDefBackgroundSortSet len(backgroundID)==0")
		return nil
	}
	key := makeDefBackgroundKey()
	IDMap := make(map[int64]string, 0)
	for i := 0; i < len(backgroundID); i++ {
		strID := strconv.FormatInt(backgroundID[i], 10)
		index := time.Now().Unix()
		IDMap[index+int64(i)] = strID
	}
	_, err := b.redisProxy.Do(ctx, "ZADD", redis.Args{}.Add(key).AddFlat(IDMap)...)
	if err != nil {
		metrics.IncrCounter("SetMeetBackgroundSortSet fail", 1)
		log.ErrorContextf(ctx, "SetMeetBackgroundSortSet fail cmd=ZAdd, err:%+v", err)
		return err
	}

	metrics.IncrCounter("SetMeetBackgroundSortSet Succ", 1)
	return nil
}

//GetDefaultBackgroundSortSet 获取默认背景图zset
func (b *BackgroundImp) GetDefaultBackgroundSortSet(ctx context.Context) (iDs []int64, err error) {
	key := makeDefBackgroundKey()
	reply, err := redis.Strings(b.redisProxy.Do(ctx, "ZRANGEBYSCORE", key, "-inf", "+inf"))
	if err != nil {
		metrics.IncrCounter("DefBackground Zset fail", 1)
		log.ErrorContextf(ctx, "GetDefaultBackgroundSortSet fail err:%+v", err)
		return nil, err
	}

	if len(reply) == 0 {
		metrics.IncrCounter("GetMeetBackgroundSortSet empty", 1)
		log.InfoContextf(ctx, "GetMeetBackgroundSortSet empty")
		return nil, nil
	}

	for _, filed := range reply {
		id, err := strconv.ParseInt(filed, 10, 64)
		if err != nil {
			metrics.IncrCounter("strID strconv.ParseInt fail", 1)
			log.ErrorContextf(ctx, "strID strconv.ParseInt fail, id:%v, err:%+v", id, err)
			continue
		}
		iDs = append(iDs, id)
	}
	metrics.IncrCounter("GetDefaultBackgroundSortSet", 1)
	return iDs, nil
}

//DeleteMeetBackground 删除会议背景图
func (b *BackgroundImp) DeleteMeetBackground(ctx context.Context, meetingID uint64, backgroundID []int64) error {
	key := MakeMeetBackgroundKey(meetingID)

	if len(backgroundID) == 0 {
		log.InfoContextf(ctx, "DeleteMeetBackground len(ids)==0")
		return nil
	}
	strIDs := make([]string, 0, len(backgroundID))
	for _, id := range backgroundID {
		strIDs = append(strIDs, strconv.FormatInt(id, 10))
	}
	//清除sortSet中的数据
	reply, err := redis.Int(b.redisProxy.Do(ctx, "ZREM", redis.Args{}.Add(key).AddFlat(strIDs)...))
	if err != nil {
		log.ErrorContextf(ctx, "DeleteMeetBackground fail, CMD=ZREM, meetID:%v, err:%+v", meetingID, err)
		return err
	}

	if reply != len(backgroundID) {
		log.ErrorContextf(ctx, "DeleteMeetBackground ZREM reply:[%v] != len(backgroundID):[%v]",
			reply, len(backgroundID))
	}

	//删除背景图信息
	strKeyList := make([]string, 0, len(backgroundID))
	for _, id := range backgroundID {
		strKey := MakeBackgroundInfoKey(id)
		strKeyList = append(strKeyList, strKey)
	}
	_, err = b.redisProxy.Do(ctx, "UNLINK", redis.Args{}.AddFlat(strKeyList)...)
	if err != nil {
		metrics.IncrCounter("DeleteMeetBackground UNLINK fail", 1)
		log.ErrorContextf(ctx, "DeleteMeetBackground UNLINK fail, meetID:%v, err:%+v, backgroundID:%+v",
			meetingID, err, backgroundID)
		return err
	}

	return nil
}

//DelDefBackgroundAllSortSet 删除默认背景图所有ID信息
func (b *BackgroundImp) DelDefBackgroundAllSortSet(ctx context.Context) error {

	key := makeDefBackgroundKey()
	_, err := b.redisProxy.Do(ctx, "UNLINK", key)
	if err != nil {
		metrics.IncrCounter("DelDefBackgroundAllSortSet unlink fail", 1)
		log.ErrorContextf(ctx, "DelDefBackgroundAllSortSet unlink fail, err:%+v", err)
		_, err = b.redisProxy.Do(ctx, "UNLINK", key)
		if err != nil {
			metrics.IncrCounter("DelDefBackgroundAllSortSet unlink again fail", 1)
			log.ErrorContextf(ctx, "DelDefBackgroundAllSortSet unlink again fail, err:%+v", err)
			return err
		}
	}

	return nil
}

//DelDefBackgroundIDListSortSet 删除默认背景图zset
func (b *BackgroundImp) DelDefBackgroundIDListSortSet(ctx context.Context, ids []int64) error {

	if len(ids) == 0 {
		log.InfoContextf(ctx, "DeleteMeetBackground len(ids)==0")
		return nil
	}

	key := makeDefBackgroundKey()
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, strconv.FormatInt(id, 10))
	}
	_, err := b.redisProxy.Do(ctx, "ZREM", redis.Args{}.Add(key).AddFlat(strIDs)...)
	if err != nil {
		metrics.IncrCounter("DelDefBackgroundIDListSortSet ZREM fail", 1)
		log.ErrorContextf(ctx, "DelDefBackgroundIDListSortSet ZREM fail, err :%+v", err)
		_, err = b.redisProxy.Do(ctx, "ZREM", redis.Args{}.Add(key).AddFlat(strIDs)...)
		if err != nil {
			metrics.IncrCounter("DelDefBackgroundIDListSortSet ZREM again fail", 1)
			log.ErrorContextf(ctx, "DelDefBackgroundIDListSortSet ZREM again fail, err :%+v", err)
			return err
		}
	}
	return nil
}

//DelBackgroundInfo 删除背景图信息
func (b *BackgroundImp) DelBackgroundInfo(ctx context.Context, ids []int64) error {

	if len(ids) == 0 {
		log.InfoContextf(ctx, "DeleteMeetBackground len(ids)==0")
		return nil
	}

	strKeys := make([]string, 0, len(ids))
	for _, id := range ids {
		key := MakeBackgroundInfoKey(id)
		strKeys = append(strKeys, key)
	}
	_, err := b.redisProxy.Do(ctx, "UNLINK", redis.Args{}.AddFlat(strKeys)...)
	if err != nil {
		metrics.IncrCounter("DelBackgroundInfo UNLINK fail", 1)
		log.ErrorContextf(ctx, "DelBackgroundInfo UNLINK fail, err:%+v", err)
		_, err = b.redisProxy.Do(ctx, "UNLINK", redis.Args{}.AddFlat(strKeys)...)
		if err != nil {
			metrics.IncrCounter("DelBackgroundInfo UNLINK again fail", 1)
			log.ErrorContextf(ctx, "DelBackgroundInfo UNLINK again fail, err:%+v", err)
			return err
		}
	}
	return nil
}

//SetDefBackgroundListInfo ..
func (b *BackgroundImp) SetDefBackgroundListInfo(ctx context.Context, backgroundList []*pb.BackgroundInfo) error {

	//获取默认背景图的ID，区分哪些是新增或者删除
	ids, err := b.GetDefaultBackgroundSortSet(ctx)
	if err != nil {
		metrics.IncrCounter("GetDefaultBackgroundSortSet fail", 1)
		log.ErrorContextf(ctx, "GetDefaultBackgroundSortSet fail, err :%+v", err)
		ids, err = b.GetDefaultBackgroundSortSet(ctx)
		if err != nil {
			metrics.IncrCounter("GetDefaultBackgroundSortSet fail again", 1)
			log.ErrorContextf(ctx, "GetDefaultBackgroundSortSet fail again, err :%+v", err)
		}
		//失败，清空后重新进行排序
		b.DelDefBackgroundAllSortSet(ctx)

	}
	log.InfoContextf(ctx, "SetDefBackgroundListInfo backgroundList:%+v", backgroundList)

	if len(backgroundList) == 0 {
		b.DelDefBackgroundAllSortSet(ctx)
		b.DelBackgroundInfo(ctx, ids)
		return nil
	}

	reqIds := make([]int64, 0, len(backgroundList))
	for _, val := range backgroundList {
		reqIds = append(reqIds, val.GetInt64BackgroundId())
	}

	if len(ids) == 0 { //为0的话设置新的背景信息
		b.SetDefBackgroundSortSet(ctx, reqIds)
		b.SetBackgroundList(ctx, backgroundList)
	} else { //对ID进行检敛
		reqIDMap := make(map[int64]struct{}, 0)
		IdMap := make(map[int64]struct{}, 0)
		for _, reqID := range reqIds {
			reqIDMap[reqID] = struct{}{}
		}
		for _, id := range ids {
			IdMap[id] = struct{}{}
		}
		//判断是否有新增的
		addBackgroundIDList := make([]int64, 0)
		for _, reqId := range reqIds {
			if _, ok := IdMap[reqId]; !ok {
				addBackgroundIDList = append(addBackgroundIDList, reqId)
			}
		}
		//判断是否有新增的
		delBackgroundIDList := make([]int64, 0)
		for _, ID := range ids {
			if _, ok := reqIDMap[ID]; !ok {
				delBackgroundIDList = append(delBackgroundIDList, ID)
			}
		}

		//删除的，操作KV删除
		b.DelBackgroundInfo(ctx, delBackgroundIDList)
		b.DelDefBackgroundIDListSortSet(ctx, ids)
		//新增的，操作KV新增
		if len(backgroundList) == 0 {
			return nil
		}
		b.SetBackgroundList(ctx, backgroundList)
		b.SetDefBackgroundSortSet(ctx, addBackgroundIDList)
	}

	return nil
}

//GetTempImageUrlByID 通过cosID换取tmpURL
func (b *BackgroundImp) GetTempImageUrlByID(ctx context.Context, backgroundList []*pb.BackgroundInfo) error {
	picIds := make([]string, 0, len(backgroundList))
	for _, val := range backgroundList {
		if val.GetStrPicId() == "" {
			continue
		}
		picIds = append(picIds, val.GetStrPicId())
	}
	if len(picIds) == 0 {
		return nil
	}
	picUrlMap, err := rpc.BatchTempUrl(ctx, picIds)
	if err != nil {
		metrics.IncrCounter("rpc.BatchTempUrl fail", 1)
		log.ErrorContextf(ctx, "rpc.BatchTempUrl fail, err:%+v", err)
		return err
	}

	for _, data := range backgroundList {
		if data.GetStrPicId() == "" {
			continue
		}
		data.StrPicUrl = proto.String(picUrlMap[data.GetStrPicId()])
	}
	return nil
}

//GetMeetBackgroundSize ..
func (b *BackgroundImp) GetMeetBackgroundSize(ctx context.Context, meetingID uint64) (int64, error) {
	key := MakeMeetBackgroundKey(meetingID)

	rst, err := b.redisProxy.Do(ctx, "ZCARD", key)
	if err != nil {
		metrics.IncrCounter("GetMeetBackgroundSize error", 1)
		log.ErrorContextf(ctx, "GetMeetBackgroundSize error, cmd = ZCARD, key:%v, err:[%v]", key, err)
		return 0, err
	}

	tmpStr := fmt.Sprintf("%v", rst)
	nodeCount, _ := strconv.ParseInt(tmpStr, 10, 64)
	log.InfoContextf(ctx, "GetMeetBackgroundSize redis ok, cmd = ZCARD, key:%v  rst:[%v], nodeCount:[%v] ",
		key, rst, nodeCount)
	return nodeCount, nil
}

//MakeMeetBackgroundKey ..
func MakeMeetBackgroundKey(meetID uint64) string {
	return fmt.Sprintf("background_meet_%v", meetID)
}

//makeDefBackgroundKey ..
func makeDefBackgroundKey() string {
	return fmt.Sprintf("background_default")
}

//MakeBackgroundInfoKey ..
func MakeBackgroundInfoKey(index int64) string {
	return fmt.Sprintf("background_info_%v", index)
}

//MakeBackgroundIndexKey ..
func MakeBackgroundIndexKey() string {
	return fmt.Sprintf("background_index")
}

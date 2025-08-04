package namebadge

import (
	"context"
	"fmt"
	"time"

	"meeting_template/material_control/cache"

	"git.code.oa.com/trpc-go/trpc-database/redis"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
	"github.com/golang/protobuf/proto"
)

//NameBadgeImp ..
type NameBadgeImp struct {
	redisProxy redis.Client
}

//NewNameBadge 初始化
func NewNameBadge() cache.NameBadges {
	imp := &NameBadgeImp{
		redisProxy: redis.NewClientProxy("trpc.wemeet.wemeet_template.redis_storage"),
	}
	_, err := imp.redisProxy.Do(context.Background(), "ping")
	if err != nil {
		panic(err)
	}

	return imp
}

//SetDefNameBadgeList 设置默认名片列表
func (n *NameBadgeImp) SetDefNameBadgeList(ctx context.Context, nameBadgeList []*pb.NameBadgeInfo) error {
	if len(nameBadgeList) == 0 {
		log.InfoContextf(ctx, "SetDefNameBadgeList empty")
		return nil
	}

	InfoMap := make(map[string]string, 0)
	for _, val := range nameBadgeList {
		v, err := proto.Marshal(val)
		if err != nil {
			metrics.IncrCounter("NameBadgeInfo proto.Marshal fail", 1)
			log.ErrorContextf(ctx, "NameBadgeInfo proto.Marshal fail, val:%+v, err:%+v", val, err)
		}
		InfoMap[val.GetStrNamebadgeId()] = string(v)
	}

	_, err := n.redisProxy.Do(ctx, "MSET", redis.Args{}.AddFlat(InfoMap)...)
	if err != nil {
		metrics.IncrCounter("SetDefNameBadgeList cmd=MSET fail", 1)
		log.ErrorContextf(ctx, "SetDefNameBadgeList cmd=MSET fail, err:%+v, nameBadgeList:%+v",
			err, nameBadgeList)
		return err
	}

	metrics.IncrCounter("SetDefNameBadgeList Succ", 1)
	return nil
}

//QueryDefNameBadgeInfoList 查询名片样式列表
func (n *NameBadgeImp) QueryDefNameBadgeInfoList(ctx context.Context) ([]*pb.NameBadgeInfo, error) {

	ids, err := n.GetDefNameBadgeSortSet(ctx)
	if err != nil {
		metrics.IncrCounter("GetDefNameBadgeSortSet fail", 1)
		log.ErrorContextf(ctx, "GetDefNameBadgeSortSet fail, err:%+v", err)
		return nil, err
	}

	if len(ids) == 0 {
		log.InfoContextf(ctx, "GetDefNameBadgeSortSet len(ids)==0")
		return nil, nil
	}

	reply, err := redis.Strings(n.redisProxy.Do(ctx, "MGET", redis.Args{}.AddFlat(ids)...))
	if err != nil {
		metrics.IncrCounter("QueryDefNameBadgeInfoList fail", 1)
		log.ErrorContextf(ctx, "QueryDefNameBadgeInfoList fail, ids:%+v, err:%+v", ids, err)
		return nil, err
	}
	nameBadgeList := make([]*pb.NameBadgeInfo, 0, len(reply))
	for _, val := range reply {
		nameBadge := &pb.NameBadgeInfo{}
		err = proto.Unmarshal([]byte(val), nameBadge)
		if err != nil {
			metrics.IncrCounter("NameBadgeInfo proto.Unmarshal", 1)
			log.ErrorContextf(ctx, "NameBadgeInfo proto.Unmarshal, err:%+v", err)
			continue
		}
		nameBadgeList = append(nameBadgeList, nameBadge)
	}

	metrics.IncrCounter("QueryDefNameBadgeInfoList Succ", 1)
	return nameBadgeList, nil
}

//SetDefNameBadgeSortSet 设置默认名片sortset
func (n *NameBadgeImp) SetDefNameBadgeSortSet(ctx context.Context, iDs []string) error {

	if len(iDs) == 0 {
		log.InfoContextf(ctx, "SetDefNameBadgeSortSet len(IDs)==0")
		return nil
	}
	key := makeDefNameBadgeKey()
	IDMap := make(map[int64]string, 0)
	start := time.Now().Unix()
	for i := 0; i < len(iDs); i++ {
		IDMap[start+int64(i)] = iDs[i]
	}

	_, err := n.redisProxy.Do(ctx, "ZADD", redis.Args{}.Add(key).AddFlat(IDMap)...)
	if err != nil {
		metrics.IncrCounter("SetDefNameBadgeSortSet fail", 1)
		log.ErrorContextf(ctx, "SetDefNameBadgeSortSet fail,IDs:%+v, err:%+v", iDs, err)
		_, err = n.redisProxy.Do(ctx, "ZADD", redis.Args{}.Add(key).AddFlat(IDMap)...)
		if err != nil {
			metrics.IncrCounter("SetDefNameBadgeSortSet again fail", 1)
			log.ErrorContextf(ctx, "SetDefNameBadgeSortSet again fail,IDs:%+v, err:%+v", iDs, err)
			return err
		}
	}

	metrics.IncrCounter("SetDefNameBadgeSortSet Succ", 1)
	return nil
}

//GetDefNameBadgeSortSet ..
func (n *NameBadgeImp) GetDefNameBadgeSortSet(ctx context.Context) ([]string, error) {
	key := makeDefNameBadgeKey()
	reply, err := redis.Strings(n.redisProxy.Do(ctx, "ZRANGEBYSCORE", key, "-inf", "+inf"))
	if err != nil {
		metrics.IncrCounter("GetDefNameBadgeSortSet fail", 1)
		log.ErrorContextf(ctx, "GetDefNameBadgeSortSet fail err:%+v", err)
		return nil, err
	}

	ids := make([]string, 0, len(reply))
	for _, id := range reply {
		ids = append(ids, id)
	}

	return ids, nil
}

//makeDefNameBadgeKey ..
func makeDefNameBadgeKey() string {
	return fmt.Sprintf("namebadge_default")
}

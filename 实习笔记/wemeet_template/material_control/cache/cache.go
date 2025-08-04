package cache

import (
	"context"

	pb "git.code.oa.com/trpcprotocol/wemeet/meeting_template"
)

// Backgrounds 虚拟背景
type Backgrounds interface {
	GetBackgroundIndex(ctx context.Context) (index int64)
	SetBackgroundList(ctx context.Context, background []*pb.BackgroundInfo) error
	GetBackground(ctx context.Context, ids []int64) (background []*pb.BackgroundInfo, err error)

	SetMeetBackgroundSortSet(ctx context.Context, meetingID uint64, backgroundID []int64) error
	GetMeetBackgroundSortSet(ctx context.Context, meetingID uint64) (iDs []int64, err error)
	GetMeetBackgroundSize(ctx context.Context, meetingID uint64) (int64, error)

	SetDefBackgroundSortSet(ctx context.Context, backgroundID []int64) error
	GetDefaultBackgroundSortSet(ctx context.Context) (iDs []int64, err error)

	DeleteMeetBackground(ctx context.Context, meetingID uint64, backgroundID []int64) error

	//SetDefBackgroundListInfo 默认背景图设置
	SetDefBackgroundListInfo(ctx context.Context, backgroundList []*pb.BackgroundInfo) error

	DelDefBackgroundAllSortSet(ctx context.Context) error
	DelDefBackgroundIDListSortSet(ctx context.Context, ids []int64) error
	DelBackgroundInfo(ctx context.Context, ids []int64) error

	GetTempImageUrlByID(ctx context.Context, backgroundList []*pb.BackgroundInfo) error

	SetCache180DaysExpireTimeDuration(ctx context.Context, key string, startTime uint32) error
}

// NameBadges 名牌样式
type NameBadges interface {
	SetDefNameBadgeList(ctx context.Context, nameBadgeList []*pb.NameBadgeInfo) error
	QueryDefNameBadgeInfoList(ctx context.Context) ([]*pb.NameBadgeInfo, error)

	SetDefNameBadgeSortSet(ctx context.Context, iDs []string) error
	GetDefNameBadgeSortSet(ctx context.Context) ([]string, error)
}

// WelCome 欢迎语
type WelCome interface {
	GetMeetingWelComeInfo(ctx context.Context, meetingID uint64) (*pb.WelComeCache, error)
	SetMeetingWelComeInfo(ctx context.Context, meetingID uint64, welcomeInfo *pb.WelComeCache, expireTime uint32) error
	SetWelComeInfoExpireTime(ctx context.Context, meetingID uint64, expireTime uint32) error
}

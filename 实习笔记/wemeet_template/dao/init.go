package dao

import (
	"sync"

	"git.code.oa.com/trpc-go/trpc-database/gorm"
	rawGorm "gorm.io/gorm"
)

const (
	wemeetTemplateGormReadonlyMySQL = "trpc.wemeet.db_proxy.mysql_readonly"
)

//DbProxy struct
type DbProxy struct {
	TemplateReadonlyProxy *rawGorm.DB
}

//DbInstance 数据库单例
var DbInstance DbProxy
var DbOnce sync.Once

// InitDb 初始化数据库
func InitDb() {
	DbOnce.Do(func() {
		db, err := gorm.NewClientProxy(wemeetTemplateGormReadonlyMySQL)
		if err != nil {
			panic("连接数据库失败: " + err.Error())
		}
		DbInstance.TemplateReadonlyProxy = db
	})
}

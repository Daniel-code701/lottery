package datasource

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"lottery/conf"
	"sync"
)

var dbLock sync.Mutex
var masterInstance *xorm.Engine
var slaveInstance *xorm.Engine

//单例
func InstanceDbMaster() *xorm.Engine {
	if masterInstance != nil {
		return masterInstance
	}
	dbLock.Lock()
	defer dbLock.Unlock()

	//单例重点 这里需要再判断一次实例是否存在

	if masterInstance != nil {
		return masterInstance
	}
	return NewDbMaster()
}

func NewDbMaster() *xorm.Engine {
	sourcename := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		conf.DbMaster.User,
		conf.DbMaster.Pwd,
		conf.DbMaster.Host,
		conf.DbMaster.Port,
		conf.DbMaster.Database)
	instance, err := xorm.NewEngine(conf.DriverName, sourcename)

	if err != nil {
		//打印日志
		log.Fatal("dbhelper.InstanceDbMaster error", err)
	}
	//数据库调试模式 false 显示sql语句
	instance.ShowSQL(false)
	masterInstance = instance
	return masterInstance
}
package entities

import _ "github.com/go-sql-driver/mysql"
import "github.com/go-xorm/xorm"

var Engine *xorm.Engine

func init() {
	var err error
	Engine, err = xorm.NewEngine("mysql", "root:root@/activityplus?charset=utf8")
	if err != nil {
		panic(err)
	}
}

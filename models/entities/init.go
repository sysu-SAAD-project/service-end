package entities

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

// Engine is mysql engine
var Engine *xorm.Engine

func init() {
	DBADDRESS := os.Getenv("DATABASE_ADDRESS")
	DBPORT := os.Getenv("DATABASE_PORT")
	url := fmt.Sprintf("root:root@%s%s/activityplus?charset=utf8", DBADDRESS, DBPORT)
	var err error
	engine, err := xorm.NewEngine("mysql", url)
	if err != nil {
		panic(err)
	}
	Engine = engine
	if os.Getenv("DEVELOP") == "TRUE" {
		Engine.ShowSQL(true)
		Engine.Logger().SetLevel(core.LOG_DEBUG)
	}
}

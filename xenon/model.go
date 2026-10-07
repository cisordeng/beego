package xenon

import (
	"fmt"
	"github.com/cisordeng/beego"
	"github.com/cisordeng/beego/orm"
)

func init() {

}

func RegisterModels() {

	dbUsed, _ := beego.AppConfig.Bool("db::DB_USED")
	if !dbUsed {
		return
	}

	// set default database
	maxIdle := 30
	maxConn := 100

	dbDriver := beego.AppConfig.String("db::DB_DRIVER")
	dbURL := ""

	switch dbDriver {
	case "mysql":
		host := beego.AppConfig.String("db::DB_HOST")
		port := beego.AppConfig.String("db::DB_PORT")
		db := beego.AppConfig.String("db::DB_NAME")
		user := beego.AppConfig.String("db::DB_USER")
		password := beego.AppConfig.String("db::DB_PASSWORD")
		charset := beego.AppConfig.String("db::DB_CHARSET")
		dbURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=Asia%%2FShanghai", user, password, host, port, db, charset)
	case "sqlite3":
		path := beego.AppConfig.String("db::DB_PATH")
		dbURL = path
	default:
		beego.Warn(fmt.Sprintf("unknown database driver %s", dbDriver))
	}

	beego.Notice(fmt.Sprintf("connect %s: ", dbDriver), dbURL)
	orm.RegisterDataBase("default", dbDriver, dbURL, maxIdle, maxConn)
	orm.RunSyncdb("default", false, true)
}

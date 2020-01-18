package xenon

import (
	"context"

	"github.com/cisordeng/beego"
	"github.com/cisordeng/beego/orm"
)

func newContextWithOrm() context.Context {
	dbUsed, _ := beego.AppConfig.Bool("db::DB_USED")
	if !dbUsed {
		return nil
	}
	bContext := context.Background()
	bContext = context.WithValue(bContext, "orm", orm.NewOrm())
	return bContext
}

func GetOrmFromContext(ctx context.Context) orm.Ormer {
	o := ctx.Value("orm")
	if o == nil {
		return nil
	}
	return o.(orm.Ormer)
}
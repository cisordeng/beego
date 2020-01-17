package xenon

import (
	"context"
	"fmt"
	"runtime"

	"github.com/cisordeng/beego"
	"github.com/cisordeng/beego/logs"
	"github.com/cisordeng/beego/orm"
	"github.com/cisordeng/beego/toolbox"
)



type CronTask struct {
	Name string
	Spec string
	Func func(ctx context.Context)
	EnableTx bool
}

func newCronTask(name string, spec string, f func(ctx context.Context), enableTx ...bool) *CronTask {
	enable := true
	if len(enableTx) > 0 {
		enable = enableTx[0]
	}
	dbUsed, _ := beego.AppConfig.Bool("db::DB_USED")
	if !dbUsed {
		enable = false
	}
	return &CronTask{
		Name: name,
		Spec: spec,
		Func: f,
		EnableTx: enable,
	}
}

func newContextWithOrm() context.Context {
	dbUsed, _ := beego.AppConfig.Bool("db::DB_USED")
	if !dbUsed {
		return nil
	}
	bContext := context.Background()
	bContext = context.WithValue(bContext, "orm", orm.NewOrm())
	return bContext
}

func (this *CronTask) start() {
	bCtx := newContextWithOrm()

	f := func() (err error) {
		defer func() {
			if err := recover(); err != nil {
				beego.Error(fmt.Sprintf("cron: task [%s] is error", this.Name))
				beego.Error(err)
				if this.EnableTx {
					o := GetOrmFromContext(bCtx)
					e := o.Rollback()
					beego.Warn("[ORM] rollback transaction")
					if e != nil {
						beego.Error(e)
					}
				}

				msg := make([]string, 0)
				for i := 1; ; i++ {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					msg = append(msg, fmt.Sprintf("%s:%d", file, line))
				}
				msg = msg[2:len(msg) - 3]
				for _, m := range msg {
					logs.Critical(m)
				}
			}
		}()

		beego.Notice(fmt.Sprintf("cron: task [%s] is start", this.Name))
		if this.EnableTx {
			o := GetOrmFromContext(bCtx)
			e := o.Begin()
			beego.Debug("[ORM] begin transaction")
			if e != nil {
				beego.Error(e)
			}
		}
		this.Func(bCtx)
		if this.EnableTx {
			o := GetOrmFromContext(bCtx)
			e := o.Commit()
			beego.Debug("[ORM] commit transaction")
			if e != nil {
				beego.Error(e)
			}
		}
		beego.Notice(fmt.Sprintf("cron: task [%s] is end", this.Name))
		return err
	}
	task := toolbox.NewTask(this.Name, this.Spec, f)
	toolbox.AddTask(this.Name, task)
	toolbox.StartTask()
}

var CronTasks []*CronTask

func RegisterCronTask(taskName string, taskSpec string, taskFunc func(ctx context.Context), enableTx ...bool) {
	task := newCronTask(taskName, taskSpec, taskFunc, enableTx...)
	CronTasks = append(CronTasks, task)
}

func RegisterCronTasks() {
	for _, cronTask := range CronTasks {
		beego.Info("+cron: "+cronTask.Name, cronTask.Spec)
		cronTask.start()
	}
}
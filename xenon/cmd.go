package xenon

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"github.com/cisordeng/beego"
	"github.com/cisordeng/beego/logs"
)



type Cmd struct {
	Name string
	Func func(ctx context.Context)
	EnableTx bool
}

func newCmd(name string, f func(ctx context.Context), enableTx ...bool) *Cmd {
	enable := true
	if len(enableTx) > 0 {
		enable = enableTx[0]
	}
	dbUsed, _ := beego.AppConfig.Bool("db::DB_USED")
	if !dbUsed {
		enable = false
	}
	return &Cmd{
		Name: name,
		Func: f,
		EnableTx: enable,
	}
}

func (this *Cmd) run() {
	bCtx := newContextWithOrm()

	defer func() {
		if err := recover(); err != nil {
			beego.Error(fmt.Sprintf("cmd: cmd [%s] is error", this.Name))
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
			msg = msg[2:len(msg) - 6]
			for _, m := range msg {
				logs.Critical(m)
			}
		}
	}()

	beego.Notice(fmt.Sprintf("cmd: cmd [%s] is start", this.Name))
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
	beego.Notice(fmt.Sprintf("cmd: cmd [%s] is end", this.Name))
}

var FileName2Cmd = make(map[string]*Cmd, 0)

func RegisterCmd(cmdName string, cmdFunc func(ctx context.Context), enableTx ...bool) {
	_, file, _, _ := runtime.Caller(1)
	_, filename := path.Split(file)
	FileName2Cmd[filename] = newCmd(cmdName, cmdFunc, enableTx...)
}

func RunCmd(fileName string) {
	if len(fileName) < 3 || fileName[:3] != ".go" {
		fileName += ".go"
	}
	cmd := FileName2Cmd[fileName]
	if cmd == nil {
		beego.Error(fmt.Sprintf("cmd: file name is [%s] of cmd does not exist", fileName))
		return
	}
	cmd.run()
}
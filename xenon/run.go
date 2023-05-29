package xenon

import (
	"github.com/cisordeng/beego"
)

func Run(args []string) {
	RegisterModels()
	if len(args) > 1 {
		fileName := args[1]
		RunCmd(fileName)
		return
	}
	RegisterResources()
	RegisterCronTasks()

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = false
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.BConfig.RecoverFunc = RecoverPanic
	beego.Run()
}

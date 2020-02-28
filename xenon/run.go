package xenon

import (
	"github.com/cisordeng/beego"
	"github.com/cisordeng/beego/plugins/cors"
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

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = false
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.BConfig.RecoverFunc = RecoverPanic
	beego.Run()
}

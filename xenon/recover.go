package xenon

import (
	"github.com/cisordeng/beego/context"
)

func RecoverPanic(ctx *context.Context) {
	if err := recover(); err != nil {
		var resp Map
		if be, ok := err.(BusinessError); ok {
			resp = Map{
				"code":        531,
				"data":        "",
				"errCode":     be.ErrCode,
				"errMsg":      be.ErrMsg,
				"innerErrMsg": "",
			}
		} else {
			resp = Map{
				"code":        531,
				"data":        "",
				"errCode":     "",
				"errMsg":      err,
				"innerErrMsg": "",
			}
		}
		err = ctx.Output.JSON(resp, true, true)
	}
}
package xenon

import (
	"errors"

	"github.com/cisordeng/beego"
)

type Error struct {
	InnerErr error
	ErrCode string
	ErrMsg  string
}

func depictError(err error, errStr []string) Error {
	errCode := "internal error"
	errMsg := "内部错误"
	if len(errStr) > 0 {
		errCode = errStr[0]
	}
	if len(errStr) > 1 {
		errMsg = errStr[1]
	}
	return Error{
		ErrCode: errCode,
		ErrMsg:	errMsg,
		InnerErr: err,
	}
}

func RaiseException(errStr ...string) {
	beego.Error(errStr)
	panic(depictError(errors.New("raise exception"),  errStr))
}

func PanicNotNilError(err error, errStr ...string) {
	if err != nil {
		beego.Error(err)
		panic(depictError(err,  errStr))
	}
}

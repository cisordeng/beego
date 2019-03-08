package xenon

import (
	"fmt"
	"github.com/cisordeng/beego"
)

var Resources []RestResourceInterface

type Map map[string]interface{}

type RestResource struct {
	beego.Controller
	Error *Error
}

type RestResourceInterface interface {
	beego.ControllerInterface
	Resource() string
	Params() map[string][]string
}



type Error struct {
	Business *BusinessError
	Inner error
}

type BusinessError struct {
	ErrCode string
	ErrMsg string
}

type Response struct {
	Code 		int32 		`json:"code"`
	Data 		interface{} `json:"data"`
	ErrCode 	string 		`json:"errCode"`
	ErrMsg 		string 		`json:"errMsg"`
	InnerErrMsg string 		`json:"innerErrMsg"`
}

func (r *RestResource) Resource() string {
	return ""
}

func (r *RestResource) Params() map[string][]string {
	return nil
}

func (r *RestResource) CheckParams () {
	method := r.Ctx.Input.Method()
	app := r.AppController.(RestResourceInterface)
	method2params := app.Params()
	if method2params != nil {
		if params, ok := method2params[method]; ok {
			actualParams := r.Input()
			for _, param := range params {
				if _, ok := actualParams[param]; !ok {
					r.Error.Business = &BusinessError{
						"rest:missing_argument",
						fmt.Sprintf("missing or invalid argument: %s", param),
					}
					r.ReturnJSON(nil)
					return
				}
			}
		}
	}
}

func (r *RestResource) Prepare() {
	r.CheckParams()
}

func (r *RestResource) MakeResponse(data Map) *Response {
	response := &Response{
		200,
		data,
		"",
		"",
		"",
	}
	if r.Error != nil {
		response = &Response{ // 指针引用
			500,
			"",
			r.Error.Business.ErrCode,
			r.Error.Business.ErrMsg,
			r.Error.Inner.Error(),
		}
	}
	return response
}

func (r *RestResource) ReturnJSON(data Map) {
	response := r.MakeResponse(data)
	r.Data["json"] = response
	r.ServeJSON()
}
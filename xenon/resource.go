package xenon

import (
	"errors"
	"fmt"
	"github.com/cisordeng/beego"
	"net/http"
)

var Resources []RestResourceInterface

type Map map[string]interface{}

type RestResource struct {
	beego.Controller
	BCtx BCtx
}

type RestResourceInterface interface {
	beego.ControllerInterface
	Resource() string
	Params() map[string][]string
}

type BCtx struct {
	Req *http.Request
	Errors []Error
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

func RegisterResource(resourceInterface RestResourceInterface) {
	Resources = append(Resources, resourceInterface)
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
					r.BCtx.Errors = append(r.BCtx.Errors, Error{NewBusinessError("rest:missing_argument", fmt.Sprintf("missing or invalid argument: [%s]", param)), errors.New("missing param")})
					return
				}
			}
		}
	}
}

func (r *RestResource) Prepare() {
	r.BCtx.Req = r.Ctx.Input.Context.Request
	r.CheckParams()
}

func NewBusinessError(ErrCode string, ErrMsg string) *BusinessError {
	return &BusinessError{
		ErrCode: ErrCode,
		ErrMsg: ErrMsg,
	}
}

func (r *RestResource) MakeResponse(data Map) *Response {
	response := &Response{
		200,
		data,
		"",
		"",
		"",
	}
	for _, Error := range r.BCtx.Errors {
		if Error.Inner != nil {
			response = &Response{ // 指针引用
				500,
				"",
				Error.Business.ErrCode,
				Error.Business.ErrMsg,
				Error.Inner.Error(),
			}
			return response
		}
	}
	return response
}

func (r *RestResource) ReturnJSON(data Map) {
	response := r.MakeResponse(data)
	r.Data["json"] = response
	r.ServeJSON()
}
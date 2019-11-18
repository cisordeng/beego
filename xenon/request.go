package xenon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/cisordeng/beego"
)

func Get(service string, resource string, data Map, apiUrls ...string) Map {
	return request("GET", service, resource, data, apiUrls...)
}

func Put(service string, resource string, data Map, apiUrls ...string) Map {
	return request("PUT", service, resource, data, apiUrls...)
}

func Post(service string, resource string, data Map, apiUrls ...string) Map {
	return request("POST", service, resource, data, apiUrls...)
}

func Delete(service string, resource string, data Map, apiUrls ...string) Map {
	return request("DELETE", service, resource, data, apiUrls...)
}

func request(method string, service string, resource string, data Map, apiUrls ...string) Map {
	apiUrl := beego.AppConfig.String("api::apiUrl")
	if len(apiUrls) > 0 {
		apiUrl = apiUrls[0]
	}
	if apiUrl[len(apiUrl) - 1] != '/' {
		apiUrl = apiUrl + "/"
	}
	params := url.Values{"__source_service": {beego.AppConfig.String("appname")}}
	for k, v := range data {
		value := ""
		switch t := v.(type) {
		case int:
			value = fmt.Sprintf("%d", v)
		case bool:
			value = fmt.Sprintf("%t", v)
		case string:
			value = v.(string)
		case float64:
			value = fmt.Sprintf("%f", v)
		default:
			beego.Warn("unknown type: ", t)
		}
		params.Set(k, value)
	}

	requestUrl := fmt.Sprintf("%s%s/%s/?%s", apiUrl, service, strings.Replace(resource, ".", "/", -1), params.Encode())
	beego.Notice(fmt.Sprintf("request url: %s %s", requestUrl, method))

	request, err := http.NewRequest(method, requestUrl, nil)
	PanicNotNilError(err)
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	PanicNotNilError(err)
	bytes, err := ioutil.ReadAll(response.Body)
	PanicNotNilError(err)
	resMap := Map{}
	err = json.Unmarshal(bytes, &resMap)
	beego.Info(string(bytes))
	PanicNotNilError(err)
	return resMap
}
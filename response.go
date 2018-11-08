package main

import (
	"fmt"
	"net/http"
	"github.com/json-iterator/go"
)

var ffjson = jsoniter.ConfigCompatibleWithStandardLibrary

// NewResponse 创建一个新的response对象
func NewResponse(w http.ResponseWriter, r *http.Request) Response {
	r.ParseForm()
	return Response{
		w: w,
		r: r,
	}
}

// ReturnSuccess 返回正确的信息
func (r *Response) ReturnSuccess(data interface{}) {
	r.Return(data, 200, "")
}

// ReturnError 返回错误信息
func (r *Response) ReturnError(statuscode int, errMsg string) {
	r.Return(nil, statuscode, errMsg)
}

// Return 向客户返回回数据
func (r *Response) Return(data interface{}, code int, message string) {
	jsonp := r.IsJSONP()

	rs, err := ffjson.Marshal(data)
	if err != nil {
		code = 500
		rs = []byte(fmt.Sprintf(`{"code":500, "message":"%s"}`, err.Error()))
	}

	r.w.WriteHeader(code)
	if jsonp == "" {
		r.w.Header().Add("Content-Type", "application/json")
		r.w.Write([]byte(fmt.Sprintf(`{"code":%d, "message": "%s", "data":%s}`, code, message, rs)))
	} else {
		r.w.Header().Add("Content-Type", "application/javascript")
		r.w.Write([]byte(fmt.Sprintf(`%s({"code":%d, "message":"%s", "data":%s})`, jsonp, code, message, rs)))
	}
}

// IsJSONP 是否为jsonp 请求
func (r *Response) IsJSONP() string {
	if r.r.Form.Get("callback") != "" {
		return r.r.Form.Get("callback")
	}
	return ""
}

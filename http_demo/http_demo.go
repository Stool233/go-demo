package http_demo

import (
	"encoding/json"
	"fmt"
	"time"
)
import "github.com/valyala/fasthttp"

func Get(url string, timeout time.Duration) (code int, body []byte, e error)  {
	if timeout > 0 {
		return fasthttp.GetTimeout(nil, url, timeout)
	}
	return fasthttp.Get(nil, url)
}

func get(url string, resp interface{}) error {
	var (
		code int
		body []byte
		e    error
	)
	if url == "" {
		return fmt.Errorf("invalid url")
	}
	if code, body, e = Get(url, 30); e != nil {
		return e
	}
	if !HTTPStatusOk(code) {
		return fmt.Errorf("%d:%s", code, body)
	}
	e = json.Unmarshal(body, resp)
	return e
}

func HTTPStatusOk(code int) bool {
	return fasthttp.StatusOK == code
}

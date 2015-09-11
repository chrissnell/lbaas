package resterror

import (
	"encoding/json"

	"github.com/emicklei/go-restful"
)

type ErrorResponse struct {
	Err string `json:error`
}

func WriteErrorJSON(resp *restful.Response, respCode int, err error) {
	er := ErrorResponse{}

	er.Err = err.Error()

	content, _ := json.Marshal(er)

	resp.WriteHeader(respCode)
	resp.Write(content)
}

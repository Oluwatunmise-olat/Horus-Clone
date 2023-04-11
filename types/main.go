package types

import (
	"bytes"
	"net/http"
)

type RequestLogger struct {
	RequestBody     		string `json:"request_body"`
	RequestPath     		string `json:"request_path"`
	RequestHeaders  		string `json:"request_headers"`
	RequestMethod   		string `json:"request_method"`
	RequestHost     		string `json:"request_host"`
	ResponseBody       	string `json:"response_body"`
	ResponseStatus      string `json:"response_status"`
	ReqResTime       		float64 `json:"req_res_time"`
}

type WrappedResponseStatus struct {
	Status int
	http.ResponseWriter
	Body bytes.Buffer
}


package utility

import (
	"encoding/json"
	"net/http"
)

const (
	// ResSuccess is for Response regularization
	//{"Result" : "success"}
	ResSuccess = "success"

	// ResFailed is for Response regularization
	//{"Result" : "fail"}
	ResFailed = "fail"
)

// Response provide a regular http responce payload
type Response struct {
	Message string      `json:"msg"`
	Result  string      `json:"res"`
	Data    interface{} `json:"data"`
}

// ResponseWithJSON will set context-type as json and return payload
func ResponseWithJSON(response http.ResponseWriter, code int, payload interface{}) {
	result, _ := json.Marshal(payload)
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.WriteHeader(code)
	response.Write(result)
}

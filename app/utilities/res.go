package utilities

import "encoding/json"

const (
	Failure = "internal_failure"
	Success = "success"
	NoData  = "No data found"
)

func MsgJson(msg string) []byte {

	res := make(map[string]interface{})
	res["message"] = msg
	data, _ := json.Marshal(res)
	return data
}

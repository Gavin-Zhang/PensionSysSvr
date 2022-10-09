package unmarshal

import (
	"encoding/json"
)

// JSONMarshal 对象转码成json字符串
func JSONMarshal(data interface{}) (string, error) {
	str, err := json.Marshal(data)
	return string(str), err
}

// JSONUnmarshal json字符串转码成对象
func JSONUnmarshal(jsonStr string, out interface{}) error {
	return json.Unmarshal([]byte(jsonStr), &out)
}

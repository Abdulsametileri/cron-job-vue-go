package utils

import "encoding/json"

func PrettyPrint(i interface{}) (string, error) {
	s, err := json.MarshalIndent(i, "", "\t")
	return string(s), err
}

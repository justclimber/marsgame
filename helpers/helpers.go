package helpers

import (
	"encoding/json"
	"log"
)

func AbsInt64(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func PrettyPrint(msg string, obj interface{}) {
	str, err := json.MarshalIndent(obj, "", "   ")
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(msg, string(str))
}
